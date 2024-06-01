package esv_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	esvc "github.com/zostay/go-esv-api/pkg/esv"

	"github.com/zostay/today/pkg/bible"
	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text/esv"
)

const jn11 = "In the beginning was the Word, and the Word was with God, and the Word was God."

func testServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			j := map[string]any{
				"passages": []any{jn11},
			}

			data, err := json.Marshal(j)
			if err != nil {
				panic(err)
			}

			_, _ = w.Write(data)
		},
	))

	return ts
}

func TestResolver_VerseText(t *testing.T) {
	t.Parallel()

	ts := testServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	res := &esv.Resolver{
		Client: &esvc.Client{
			BaseURL: u,
			Client:  http.DefaultClient,
			Token:   "abc123",
		},
	}

	p, err := ref.ParseProper("John 1:1")
	require.NoError(t, err)

	ref, err := bible.Protestant.Resolve(p)
	require.NoError(t, err)

	txt, err := res.VerseAs(context.Background(), &ref[0], "text")
	assert.NoError(t, err)
	assert.Equal(t, jn11, txt)
}

func TestResolver_VerseHTML(t *testing.T) {
	t.Parallel()

	ts := testServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	res := &esv.Resolver{
		Client: &esvc.Client{
			BaseURL: u,
			Client:  http.DefaultClient,
			Token:   "abc123",
		},
	}

	p, err := ref.ParseProper("John 1:1")
	require.NoError(t, err)

	ref, err := bible.Protestant.Resolve(p)
	require.NoError(t, err)

	txt, err := res.VerseAs(context.Background(), &ref[0], "html")
	assert.NoError(t, err)
	assert.Equal(t, jn11, txt)
}

func TestResolver_VersionInformation(t *testing.T) {
	t.Parallel()

	res := &esv.Resolver{}

	v, err := res.VersionInformation(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "ESV", v.Name)
	assert.Equal(t, "https://www.esv.org/", v.Link)
}
