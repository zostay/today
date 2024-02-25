package photo_test

import (
	"context"
	"fmt"
	"image"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/photo"
)

type testSource struct{}

var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func (t *testSource) Photo(ctx context.Context, url string) (info *photo.Descriptor, err error) {
	d := &photo.Descriptor{
		Link:  url,
		Title: "Test Meta",
		Creator: photo.Creator{
			Name: "Test Creator",
			Link: "https://example.com",
		},
	}

	d.AddImage(photo.Original, &photo.File{
		Path: "unsplash/testdata/waa.jpg",
	})

	return d, nil
}

var _ photo.Source = (*testSource)(nil)

func TestService(t *testing.T) {
	t.Parallel()

	s := photo.NewService(&testSource{})

	d, err := s.Photo(context.Background(), "https://example.com")
	assert.NoError(t, err)
	assert.True(t,
		assert.ObjectsExportedFieldsAreEqual(
			&photo.Descriptor{
				Link:  "https://example.com",
				Title: "Test Meta",
				Creator: photo.Creator{
					Name: "Test Creator",
					Link: "https://example.com",
				},
			}, d),
	)

	assert.True(t, d.HasImage(photo.Original))

	img, format, err := d.GetImage(photo.Original).Image()
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", format)
	assert.Equal(t, image.Rect(0, 0, 4128, 2322), img.Bounds())

	resizeKey, err := s.ResizedImage(
		context.Background(),
		d,
		photo.MaxWidth(1000),
		photo.MaxHeight(1000),
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, resizeKey)

	assert.True(t, d.HasImage(resizeKey))

	resized, format, err := d.GetImage(resizeKey).Image()
	assert.NoError(t, err)
	assert.Equal(t, "", format)
	assert.Equal(t, image.Rect(0, 0, 1000, 563), resized.Bounds())

	c, err := photo.DominantImageColor(context.Background(), img)
	assert.NoError(t, err)
	r, g, b, a := c.RGBA()
	bgc := fmt.Sprintf("#%02x%02x%02x", r*256/a, g*256/a, b*256/a)
	assert.Equal(t, "#aeb0a7", bgc)
}
