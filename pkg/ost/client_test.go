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
	today = ost.Verse{
		Reference: "Luke 10:25",
		Content:   "And behold, a lawyer stood up to put him to the test, saying, “Teacher, what shall I do to inherit eternal life?”",
		Version: ost.Version{
			Name: "ESV",
			Link: "https://www.esv.org/Luke+10:25",
		},
	}
	image = photo.Info{
		Key: "test/https-example-com",
		Meta: &photo.Meta{
			Link:  "https://example.com",
			Title: "",
			Creator: photo.Creator{
				Name: "Test Photographer",
				Link: "https://example.com/testuser",
			},
		},
	}
)

type testResolver struct {
	lastRef *ref.Resolved
}

func (t *testResolver) Verse(ref *ref.Resolved) (string, error) {
	t.lastRef = ref
	return string(today.Content), nil
}

func (t *testResolver) VerseHTML(ref *ref.Resolved) (template.HTML, error) {
	t.lastRef = ref
	return template.HTML(today.Content), nil //nolint:gosec // srsly?
}

var _ text.Resolver = (*testResolver)(nil)

type testSource struct{}

func (t *testSource) CacheKey(url string) (string, bool) {
	panic("implement me")
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

	v, err := c.TodayVerse()
	assert.NoError(t, err)
	assert.Equal(t, &today, v)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verse.yaml", ri.path)

	txt, err := c.Today()
	assert.NoError(t, err)
	assert.Equal(t, string(today.Content), txt)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verse.yaml", ri.path)

	htxt, err := c.TodayHTML()
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(today.Content), htxt) //nolint:gosec // srsly?
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verse.yaml", ri.path)

	on := time.Date(2023, 12, 30, 0, 0, 0, 0, time.Local)

	v, err = c.TodayVerse(ost.On(on))
	assert.NoError(t, err)
	assert.Equal(t, &today, v)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/verses/2023/12/30/verse.yaml", ri.path)

	pi, err := c.TodayPhoto()
	assert.NoError(t, err)
	assert.Equal(t, &image, pi)
	assert.NoError(t, ri.err)
	assert.Equal(t, "/photo.yaml", ri.path)
	assert.NoError(t, pi.Close())

	pi, err = c.TodayPhoto(ost.On(on))
	assert.NoError(t, err)
	assert.Equal(t, &image, pi)
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

	v, err := c.TodayVerse()
	assert.Error(t, err)
	assert.Nil(t, v)

	txt, err := c.Today()
	assert.Error(t, err)
	assert.Empty(t, txt)

	htxt, err := c.TodayHTML()
	assert.Error(t, err)
	assert.Empty(t, htxt)
}
