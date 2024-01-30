package ost_test

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/ost"
	"github.com/zostay/today/pkg/photo"
	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

var (
	today = text.Verse{
		Reference: "Luke 10:25",
		Content: text.Content{
			Text: "And behold, a lawyer stood up to put him to the test, saying, “Teacher, what shall I do to inherit eternal life?”",
			HTML: "And behold, a lawyer stood up to put him to the test, saying, “Teacher, what shall I do to inherit eternal life?”",
		},
		Link: "https://www.esv.org/Luke+10:25",
		Version: text.Version{
			Name: "ESV",
			Link: "https://www.esv.org/",
		},
	}
	image = photo.Meta{
		Link:  "https://example.com",
		Title: "",
		Creator: photo.Creator{
			Name: "Test Photographer",
			Link: "https://example.com/testuser",
		},
	}
)

type testResolver struct {
	lastRef *ref.Resolved
}

func (t *testResolver) VersionInformation(_ context.Context) (*text.Version, error) {
	return &today.Version, nil
}

func (t *testResolver) Verse(_ context.Context, ref *ref.Resolved) (string, error) {
	t.lastRef = ref
	return today.Content.Text, nil
}

func (t *testResolver) VerseHTML(_ context.Context, ref *ref.Resolved) (template.HTML, error) {
	t.lastRef = ref
	return today.Content.HTML, nil //nolint:gosec // srsly?
}

var _ text.Resolver = (*testResolver)(nil)

type testSource struct{}

func (t *testSource) CacheKey(url string) (string, bool) {
	return "test/" + url, true
}

func (t *testSource) Photo(ctx context.Context, url string) (info *photo.Info, err error) {
	panic("implement me")
}

func (t *testSource) Download(ctx context.Context, info *photo.Info) error {
	panic("implement me")
}

var _ photo.Source = (*testSource)(nil)

type requestInfo struct {
	path string
	err  error
}

func testServer() (*httptest.Server, *requestInfo) {
	ri := requestInfo{}
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.HasSuffix(r.URL.Path, "/verse.yaml"):
				ri.path = r.URL.Path
				enc := yaml.NewEncoder(w)
				err := enc.Encode(today)
				ri.err = err
			case strings.HasSuffix(r.URL.Path, "/photo.yaml"):
				ri.path = r.URL.Path
				enc := yaml.NewEncoder(w)
				err := enc.Encode(image)
				ri.err = err
			default:
				w.WriteHeader(404)
			}
		},
	))
	return ts, &ri
}

func TestClient(t *testing.T) {
	t.Parallel()

	ts, ri := testServer()
	defer ts.Close()

	c := &ost.Client{
		BaseURL:      ts.URL,
		Client:       http.DefaultClient,
		TextService:  text.NewService(&testResolver{}),
		PhotoService: photo.NewService(&testSource{}),
	}

	v, err := c.TodayVerse(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, &today, v)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verse.yaml", ri.path)

	txt, err := c.Today(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, today.Content.Text, txt)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verse.yaml", ri.path)

	htxt, err := c.TodayHTML(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, today.Content.HTML, htxt) //nolint:gosec // srsly?
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verse.yaml", ri.path)

	on := time.Date(2023, 12, 30, 0, 0, 0, 0, time.Local)

	v, err = c.TodayVerse(context.Background(), ost.On(on))
	assert.NoError(t, err)
	assert.Equal(t, &today, v)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verses/2023/12/30/verse.yaml", ri.path)

	pi, err := c.TodayPhoto(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, &photo.Info{
		Key:  "test/https://example.com",
		Meta: &image,
	}, pi)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/photo.yaml", ri.path)
	assert.NoError(t, pi.Close())

	pi, err = c.TodayPhoto(context.Background(), ost.On(on))
	assert.NoError(t, err)
	assert.Equal(t, &photo.Info{
		Key:  "test/https://example.com",
		Meta: &image,
	}, pi)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verses/2023/12/30/photo.yaml", ri.path)
	assert.NoError(t, pi.Close())
}

func TestClient_Sad(t *testing.T) {
	t.Parallel()

	ts, _ := testServer()
	defer ts.Close()

	c := &ost.Client{
		BaseURL:      "%^&*",
		Client:       http.DefaultClient,
		TextService:  text.NewService(&testResolver{}),
		PhotoService: photo.NewService(&testSource{}),
	}

	v, err := c.TodayVerse(context.Background())
	assert.Error(t, err)
	assert.Nil(t, v)

	txt, err := c.Today(context.Background())
	assert.Error(t, err)
	assert.Empty(t, txt)

	htxt, err := c.TodayHTML(context.Background())
	assert.Error(t, err)
	assert.Empty(t, htxt)
}
