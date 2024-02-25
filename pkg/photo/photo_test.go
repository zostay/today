package photo_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/photo"
)

type testImage struct{}

func (i testImage) Filename() string {
	return "test"
}

func (i testImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (i testImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{1, 1},
	}
}

func (i testImage) At(x, y int) color.Color {
	return color.RGBA{0, 0, 0, 0}
}

func (i testImage) Image() (image.Image, string, error) {
	return i, "", nil
}

var _ image.Image = testImage{}
var _ photo.Image = testImage{}

func TestDescriptor_Image(t *testing.T) {
	t.Parallel()

	pm := photo.Descriptor{}

	pm.AddImage(photo.Original, testImage{})

	ok := pm.HasImage(photo.Original)
	require.True(t, ok)

	img := pm.GetImage(photo.Original)
	assert.Equal(t, photo.CompleteImage(testImage{}), img)

	pm.RemoveImage(photo.Original)

	ok = pm.HasImage(photo.Original)
	require.False(t, ok)

	img = pm.GetImage(photo.Original)
	assert.Nil(t, img)
}

func TestDescriptor_SetColor(t *testing.T) {
	t.Parallel()

	pm := photo.Descriptor{}
	pm.SetColor(color.RGBA{0xaa, 0xbb, 0xcc, 0xff})
	assert.Equal(t, "#aabbcc", pm.Color)
}

func TestDescriptor_GetColor(t *testing.T) {
	t.Parallel()

	pm := photo.Descriptor{Color: "#123456"}
	c, err := pm.GetColor()
	require.NoError(t, err)
	assert.Equal(t, color.RGBA{0x12, 0x34, 0x56, 0xff}, c)

	pm.Color = "123456"
	c, err = pm.GetColor()
	require.Error(t, err)
	assert.Nil(t, c)

	pm.Color = ""
	c, err = pm.GetColor()
	require.NoError(t, err)
	assert.Nil(t, c)

	pm.Color = "#12345"
	c, err = pm.GetColor()
	require.Error(t, err)
	assert.Nil(t, c)

	pm.Color = "#1x3456"
	c, err = pm.GetColor()
	require.Error(t, err)
	assert.Nil(t, c)

	pm.Color = "#12x456"
	c, err = pm.GetColor()
	require.Error(t, err)
	assert.Nil(t, c)

	pm.Color = "#12345x"
	c, err = pm.GetColor()
	require.Error(t, err)
	assert.Nil(t, c)
}
