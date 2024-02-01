package photo_test

import (
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/photo"
)

func TestPhotoInfo_HasDownload(t *testing.T) {
	t.Parallel()

	pi := photo.Info{}
	assert.False(t, pi.HasDownload())

	pi.File = &os.File{}
	assert.True(t, pi.HasDownload())
}

func TestPhotoInfo_Close(t *testing.T) {
	t.Parallel()

	pi := photo.Info{}
	assert.NoError(t, pi.Close())

	var err error
	pi.File, err = os.Open("unsplash/testdata/waa.jpg")
	require.NoError(t, err)
	assert.NoError(t, pi.Close())
	assert.Nil(t, pi.File)
}

func TestPhotoMeta_SetColor(t *testing.T) {
	t.Parallel()

	pm := photo.Meta{}
	pm.SetColor(color.RGBA{0xaa, 0xbb, 0xcc, 0xff})
	assert.Equal(t, "#aabbcc", pm.Color)
}

func TestPhotoMeta_GetColor(t *testing.T) {
	t.Parallel()

	pm := photo.Meta{Color: "#123456"}
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
