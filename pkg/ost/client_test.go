package ost_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

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
	return template.HTML(today.Content), nil
}

var _ text.Resolver = (*testResolver)(nil)

func testServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			enc := yaml.NewEncoder(w)
			err := enc.Encode(today)
			if err != nil {
				panic(err)
			}
		},
	))
	return ts
}

func TestClient(t *testing.T) {
	t.Parallel()

	ts := testServer()
	defer ts.Close()

	c := &ost.Client{
		BaseURL: ts.URL,
		Client:  http.DefaultClient,
		Service: text.NewService(&testResolver{}),
	}

	v, err := c.TodayVerse()
	assert.NoError(t, err)
	assert.Equal(t, &today, v)

	txt, err := c.Today()
	assert.NoError(t, err)
	assert.Equal(t, string(today.Content), txt)

	htxt, err := c.TodayHTML()
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(today.Content), htxt)
}
