package text_test

import (
	"context"
	"html/template"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

const fjn41 = `Beloved, do not believe every spirit, but test the spirits to see whether they are from God, for many false prophets have gone out into the world.`

type testResolver struct {
	lastRef *ref.Resolved
}

func (t *testResolver) VersionInformation(context.Context) (*text.Version, error) {
	return &text.Version{
		Name: "ESV",
		Link: "https://www.esv.org/",
	}, nil
}

func (t *testResolver) Verse(_ context.Context, ref *ref.Resolved) (*text.Verse, error) {
	t.lastRef = ref
	return &text.Verse{
		Reference: ref.Ref(),
		Content: text.Content{
			Text: fjn41,
			HTML: template.HTML(fjn41), //nolint:gosec // this is a test
		},
		Link: "https://www.esv.org/" + url.PathEscape(ref.Ref()),
		Version: text.Version{
			Name: "ESV",
			Link: "https://www.esv.org/",
		},
	}, nil
}

func (t *testResolver) VerseText(_ context.Context, ref *ref.Resolved) (string, error) {
	t.lastRef = ref
	return fjn41, nil
}

func (t *testResolver) VerseHTML(_ context.Context, ref *ref.Resolved) (template.HTML, error) {
	t.lastRef = ref
	return fjn41, nil
}

var _ text.Resolver = (*testResolver)(nil)

func TestService(t *testing.T) {
	t.Parallel()

	tr := &testResolver{}
	svc := text.NewService(tr)
	assert.NotNil(t, svc)

	b, err := ref.Canonical.Book("1 John")
	require.NoError(t, err)
	require.NotNil(t, b)

	ctx := context.Background()
	txt, err := svc.VerseText(ctx, "1 John 4:1")
	assert.NoError(t, err)
	assert.Equal(t, fjn41, txt)
	assert.Equal(t, &ref.Resolved{
		Book:  b,
		First: ref.CV{Chapter: 4, Verse: 1},
		Last:  ref.CV{Chapter: 4, Verse: 1},
	}, tr.lastRef)

	txt, err = svc.VerseText(ctx, "1jn 4:1")
	assert.NoError(t, err)
	assert.Equal(t, fjn41, txt)
	assert.Equal(t, &ref.Resolved{
		Book:  b,
		First: ref.CV{Chapter: 4, Verse: 1},
		Last:  ref.CV{Chapter: 4, Verse: 1},
	}, tr.lastRef)

	htxt, err := svc.VerseHTML(ctx, "1 John 4:1")
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(fjn41), htxt) //nolint:gosec // srsly?
	assert.Equal(t, &ref.Resolved{
		Book:  b,
		First: ref.CV{Chapter: 4, Verse: 1},
		Last:  ref.CV{Chapter: 4, Verse: 1},
	}, tr.lastRef)

	htxt, err = svc.VerseHTML(ctx, "1stjo 4:1")
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(fjn41), htxt) //nolint:gosec // srsly?
	assert.Equal(t, &ref.Resolved{
		Book:  b,
		First: ref.CV{Chapter: 4, Verse: 1},
		Last:  ref.CV{Chapter: 4, Verse: 1},
	}, tr.lastRef)

	r, txt, err := svc.RandomVerseText(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fjn41, txt)
	assert.NoError(t, r.Validate())

	r, htxt, err = svc.RandomVerseHTML(ctx)
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(fjn41), htxt) //nolint:gosec // srsly?
	assert.NoError(t, r.Validate())
}

// TestService_SingleChapterBook checks that looking up text for a book with no
// chapters works whether or not the caller names chapter 1. Both forms have to
// arrive at the same resolved reference.
func TestService_SingleChapterBook(t *testing.T) {
	t.Parallel()

	tr := &testResolver{}
	svc := text.NewService(tr)
	assert.NotNil(t, svc)

	b, err := ref.Canonical.Book("2 John")
	require.NoError(t, err)
	require.NotNil(t, b)

	want := &ref.Resolved{
		Book:  b,
		First: ref.N{Number: 1},
		Last:  ref.N{Number: 4},
	}

	ctx := context.Background()
	for _, vr := range []string{"2 John 1-4", "2 John 1:1-4", "2jn1.1-4"} {
		txt, err := svc.VerseText(ctx, vr)
		assert.NoError(t, err, "looking up %q", vr)
		assert.Equal(t, fjn41, txt)
		assert.Equal(t, want, tr.lastRef, "looking up %q", vr)
	}

	// a chapter the book does not have is still an error
	txt, err := svc.VerseText(ctx, "2 John 2:1")
	assert.Error(t, err)
	assert.Empty(t, txt)
}

func TestService_Sad(t *testing.T) {
	t.Parallel()

	tr := &testResolver{}
	svc := text.NewService(tr)
	assert.NotNil(t, svc)

	b, err := ref.Canonical.Book("1 John")
	require.NoError(t, err)
	require.NotNil(t, b)

	ctx := context.Background()
	txt, err := svc.Verse(ctx, "1 John 4:")
	assert.Error(t, err)
	assert.Empty(t, txt)

	htxt, err := svc.VerseHTML(ctx, "1 John 4:")
	assert.Error(t, err)
	assert.Empty(t, htxt)

	txt, err = svc.Verse(ctx, "1 John 400:1")
	assert.Error(t, err)
	assert.Empty(t, txt)

	htxt, err = svc.VerseHTML(ctx, "1 John 400:1")
	assert.Error(t, err)
	assert.Empty(t, htxt)

	txt, err = svc.Verse(ctx, "1 John 4:1; 5:1")
	assert.Error(t, err)
	assert.Empty(t, txt)

	htxt, err = svc.VerseHTML(ctx, "1 John 4:1; 5:1")
	assert.Error(t, err)
	assert.Empty(t, htxt)

	txt, err = svc.Verse(ctx, "1johnny 4:1")
	assert.Error(t, err)
	assert.Empty(t, txt)

	htxt, err = svc.VerseHTML(ctx, "1Jojo 4:1")
	assert.Error(t, err)
	assert.Empty(t, htxt)
}
