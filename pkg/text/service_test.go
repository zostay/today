package text_test

import (
	"context"
	"html/template"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/bible"
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
			"text": fjn41,
			"html": fjn41,
		},
		Link: "https://www.esv.org/" + url.PathEscape(ref.Ref()),
		Version: text.Version{
			Name: "ESV",
			Link: "https://www.esv.org/",
		},
	}, nil
}

func (t *testResolver) AsFormats() []string {
	return []string{"text", "html"}
}

type testVF struct {
	name        string
	description string
	ext         string
}

func (t testVF) Name() string {
	return t.name
}

func (t testVF) Description() string {
	return t.description
}

func (t testVF) Ext() string {
	return t.ext
}

func (t *testResolver) DescribeFormat(ofmt string) text.VerseFormat {
	switch ofmt {
	case "text":
		return testVF{
			name:        "text",
			description: "Plain text",
			ext:         "txt",
		}
	case "html":
		return testVF{
			name:        "html",
			description: "HTML",
			ext:         "html",
		}
	default:
		return nil
	}
}

func (t *testResolver) VerseAs(_ context.Context, ref *ref.Resolved, ofmt string) (string, error) {
	t.lastRef = ref
	return fjn41, nil
}

func (t *testResolver) VerseURI(_ context.Context, ref *ref.Resolved) (string, error) {
	t.lastRef = ref
	return "bible://" + url.PathEscape(ref.Book.Name) + "/" + url.PathEscape(ref.First.Ref()) + "-" + url.PathEscape(ref.Last.Ref()), nil
}

var _ text.Resolver = (*testResolver)(nil)

func TestService(t *testing.T) {
	t.Parallel()

	tr := &testResolver{}
	svc := text.NewService(tr)
	assert.NotNil(t, svc)

	b, err := bible.Protestant.Book("1 John")
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

func TestService_Sad(t *testing.T) {
	t.Parallel()

	tr := &testResolver{}
	svc := text.NewService(tr)
	assert.NotNil(t, svc)

	b, err := bible.Protestant.Book("1 John")
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
