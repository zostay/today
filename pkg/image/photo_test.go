package image_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/image"
)

func TestPhotoInfo_HasDownload(t *testing.T) {
	t.Parallel()

	pi := image.PhotoInfo{}
	assert.False(t, pi.HasDownload())

	pi.File = &os.File{}
	assert.True(t, pi.HasDownload())
}

func TestPhotoInfo_Close(t *testing.T) {
	t.Parallel()

	pi := image.PhotoInfo{}
	assert.NoError(t, pi.Close())

	var err error
	pi.File, err = os.Open("unsplash/testdata/waa.jpg")
	require.NoError(t, err)
	assert.NoError(t, pi.Close())
	assert.Nil(t, pi.File)
}
