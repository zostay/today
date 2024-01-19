package ost_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/ost"
	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

var today = ost.Verse{
	Reference: "Luke 10:25",
	Content:   "And behold, a lawyer stood up to put him to the test, saying, “Teacher, what shall I do to inherit eternal life?”",
	Version: ost.Version{
		Name: "ESV",
		Link: "https://www.esv.org/Luke+10:25",
	},
}

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

type requestInfo struct {
	path string
	err  error
}

func testServer() (*httptest.Server, *requestInfo) {
	ri := requestInfo{}
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ri.path = r.URL.Path
			enc := yaml.NewEncoder(w)
			err := enc.Encode(today)
			ri.err = err
		},
	))
	return ts, &ri
}

func TestClient(t *testing.T) {
	t.Parallel()

	ts, ri := testServer()
	defer ts.Close()

	c := &ost.Client{
		BaseURL: ts.URL,
		Client:  http.DefaultClient,
		Service: text.NewService(&testResolver{}),
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
}

func TestClient_Sad(t *testing.T) {
	t.Parallel()

	ts, _ := testServer()
	defer ts.Close()

	c := &ost.Client{
		BaseURL: "%^&*",
		Client:  http.DefaultClient,
		Service: text.NewService(&testResolver{}),
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
