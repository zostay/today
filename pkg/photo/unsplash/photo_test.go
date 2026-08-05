package unsplash_test

import (
	"context"
	"encoding/json"
	"image"
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
						"html":              baseUrl + "/photos/a-test-photo-with-title-that-does-not-matter-abc123-_XYZ",
						"download":          baseUrl + "/photos/abc123-_XYZ/download",
						"download_location": baseUrl + "/photos/abc123-_XYZ/download",
					},
					"urls": map[string]any{
						"raw":     baseUrl + "/img/photo-abc123?ixid=raw",
						"full":    baseUrl + "/img/photo-abc123?ixid=full",
						"regular": baseUrl + "/img/photo-abc123?ixid=regular",
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
		Client: unsp.New(ts.Client()),
	}

	d, err := src.Photo(context.Background(), "https://unsplash.com/photos/a-test-photo-with-title-that-does-not-matter-abc123-_XYZ")
	assert.NoError(t, err)

	assert.EqualExportedValues(t,
		&photo.Descriptor{
			Link:  u.String() + "/photos/a-test-photo-with-title-that-does-not-matter-abc123-_XYZ",
			Type:  "unsplash",
			Title: "",
			Creator: photo.Creator{
				Name: "Test User",
				Link: u.String() + "/testuser",
			},
			ImageURL:         u.String() + "/img/photo-abc123?ixid=raw",
			DownloadLocation: u.String() + "/photos/abc123-_XYZ/download",
		}, d)

	item := d.GetImage(photo.Original)
	assert.NotNil(t, item)

	img, format, err := item.Image()
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", format)
	assert.Equal(t, image.Rect(0, 0, 4128, 2322), img.Bounds())

	assert.Equal(t, "YZ/download", item.Filename())
}

// hotlinkVariantServer serves a photo whose "urls" block contains only the
// variants given, so the preference order can be exercised. The variant values
// are opaque to the code under test, so they need not point anywhere real.
func hotlinkVariantServer(urls map[string]any) *httptest.Server {
	baseUrl := ""
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Source.Photo also resolves the download link, so that endpoint
			// has to answer even though this test only inspects the hotlink.
			if r.URL.Path == "/photos/abc123-_XYZ/download" {
				if err := json.NewEncoder(w).Encode(map[string]any{
					"url": baseUrl + "/photos/abc123-_XYZ/download/actual-file",
				}); err != nil {
					w.WriteHeader(500)
				}
				return
			}

			if r.URL.Path != "/photos/abc123-_XYZ" {
				w.WriteHeader(404)
				return
			}

			j := map[string]any{
				"id": "abc123-_XYZ",
				"links": map[string]any{
					"html":     baseUrl + "/photos/a-test-photo-abc123-_XYZ",
					"download": baseUrl + "/photos/abc123-_XYZ/download",
				},
				"user": map[string]any{
					"name":  "Test User",
					"links": map[string]any{"html": baseUrl + "/testuser"},
				},
			}
			if urls != nil {
				j["urls"] = urls
			}

			if err := json.NewEncoder(w).Encode(j); err != nil {
				w.WriteHeader(500)
			}
		},
	))
	baseUrl = ts.URL
	return ts
}

// Unsplash asks that photos be displayed from the URLs under "urls". Raw is
// preferred because it carries no size preset, but a response omitting it must
// still yield a usable hotlink rather than none.
func TestSourcePhotoHotlinkURL(t *testing.T) { //nolint:paralleltest // unsplash client has globals that have to be set
	const (
		raw     = "https://images.example/photo-abc123?ixid=raw"
		full    = "https://images.example/photo-abc123?ixid=full"
		regular = "https://images.example/photo-abc123?ixid=regular"
	)

	tests := []struct {
		name string
		urls map[string]any
		want string
	}{
		{"prefers raw", map[string]any{"raw": raw, "full": full, "regular": regular}, raw},
		{"falls back to full", map[string]any{"full": full, "regular": regular}, full},
		{"falls back to regular", map[string]any{"regular": regular}, regular},
		{"no urls block at all", nil, ""},
		{"empty urls block", map[string]any{}, ""},
	}

	//nolint:paralleltest // each case sets the client's global base URL
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := hotlinkVariantServer(tc.urls)
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)
			unsp.SetupBaseUrl(u.String() + "/")

			src := &unsplash.Source{Client: unsp.New(ts.Client())}
			d, err := src.Photo(context.Background(), "https://unsplash.com/photos/a-test-photo-abc123-_XYZ")
			require.NoError(t, err)

			assert.Equal(t, tc.want, d.ImageURL)
		})
	}
}
