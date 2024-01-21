package photo_test

import (
	"context"
	"fmt"
	"image/jpeg"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/photo"
)

type testSource struct{}

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func (t *testSource) CacheKey(url string) (string, bool) {
	url = nonAlnum.ReplaceAllString(url, "-")
	return "test/" + url, true
}

func (t *testSource) Photo(ctx context.Context, url string) (info *photo.Info, err error) {
	key, _ := t.CacheKey(url)
	return &photo.Info{
		Key: key,
		Meta: &photo.Meta{
			Link:  url,
			Title: "Test Meta",
			Creator: photo.Creator{
				Name: "Test Creator",
				Link: "https://example.com",
			},
		},
	}, nil
}

func (t *testSource) Download(ctx context.Context, info *photo.Info) (err error) {
	info.File, err = os.Open("unsplash/testdata/waa.jpg")
	return
}

var _ photo.Source = (*testSource)(nil)

func TestService(t *testing.T) {
	t.Parallel()

	s := photo.NewService(&testSource{})

	pi, err := s.Photo(context.Background(), "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, &photo.Info{
		Key: "test/https-example-com",
		Meta: &photo.Meta{
			Link:  "https://example.com",
			Title: "Test Meta",
			Creator: photo.Creator{
				Name: "Test Creator",
				Link: "https://example.com",
			},
		},
	}, pi)

	assert.False(t, pi.HasDownload())

	err = s.Download(context.Background(), pi)
	assert.NoError(t, err)
	assert.True(t, pi.HasDownload())
	assert.NotNil(t, pi.File)

	assert.NoError(t, pi.Close())

	resized, err := s.ResizedImage(
		context.Background(),
		pi,
		photo.MaxWidth(1000),
		photo.MaxHeight(1000),
	)
	assert.NoError(t, err)
	assert.NotNil(t, resized)

	img, err := jpeg.Decode(resized)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	assert.Equal(t, img.Bounds().Max.X-img.Bounds().Min.X, 1000)
	assert.Equal(t, img.Bounds().Max.Y-img.Bounds().Min.Y, 563)

	assert.NoError(t, resized.Close())

	assert.NoError(t, pi.Close())

	c, err := s.DominantImageColor(context.Background(), pi)
	assert.NoError(t, err)
	r, g, b, a := c.RGBA()
	bgc := fmt.Sprintf("#%02x%02x%02x", r*256/a, g*256/a, b*256/a)
	assert.Equal(t, "#aeb0a7", bgc)

	assert.NoError(t, pi.Close())
}
