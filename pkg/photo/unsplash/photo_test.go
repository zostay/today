package unsplash_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	unsp "github.com/hbagdi/go-unsplash/unsplash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/photo"
	"github.com/zostay/today/pkg/photo/unsplash"
)

func testServer() *httptest.Server {
	baseUrl := ""
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/photos/abc123-_XYZ":
				j := map[string]any{
					"id": "abc123-_XYZ",
					"links": map[string]any{
						"html":     baseUrl + "/photos/a-test-photo-with-title-that-does-not-matter-abc123-_XYZ",
						"download": baseUrl + "/photos/abc123-_XYZ/download",
					},
					"user": map[string]any{
						"name": "Test User",
						"links": map[string]any{
							"html": baseUrl + "/testuser",
						},
					},
				}

				enc := json.NewEncoder(w)
				err := enc.Encode(j)
				if err != nil {
					w.WriteHeader(500)
				}
			case "/photos/abc123-_XYZ/download":
				j := map[string]any{
					"url": baseUrl + "/photos/abc123-_XYZ/download/actual-file",
				}

				enc := json.NewEncoder(w)
				err := enc.Encode(j)
				if err != nil {
					w.WriteHeader(500)
				}
			case "/photos/abc123-_XYZ/download/actual-file":
				r, err := os.Open("testdata/waa.jpg")
				if err != nil {
					w.WriteHeader(500)
				}
				defer r.Close()

				w.Header().Add("Content-Type", "image/jpeg")

				_, err = io.Copy(w, r)
				if err != nil {
					w.WriteHeader(500)
				}
			default:
				w.WriteHeader(404)
			}
		},
	))

	baseUrl = ts.URL

	return ts
}

func TestSource(t *testing.T) { //nolint:paralleltest // unsplash client has globals that have to be set
	ts := testServer()
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	// A global variable for this? Srsly?
	unsp.SetupBaseUrl(u.String() + "/")
	src := &unsplash.Source{
		Client: unsp.New(http.DefaultClient),
	}

	pi, err := src.Photo(context.Background(), "https://unsplash.com/photos/a-test-photo-with-title-that-does-not-matter-abc123-_XYZ")
	assert.NoError(t, err)
	assert.Equal(t, &photo.Info{
		Key: "unsplash/abc123-_XYZ",
		Meta: &photo.Meta{
			Link:  u.String() + "/photos/a-test-photo-with-title-that-does-not-matter-abc123-_XYZ",
			Type:  "unsplash",
			Title: "",
			Creator: photo.Creator{
				Name: "Test User",
				Link: u.String() + "/testuser",
			},
		},
	}, pi)

	assert.False(t, pi.HasDownload())

	err = src.Download(context.Background(), pi)
	assert.NoError(t, err)
	assert.True(t, pi.HasDownload())
	assert.NotNil(t, pi.File)

	assert.NoError(t, pi.Close())
}
