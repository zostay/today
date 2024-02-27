package photo_test

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/photo"
)

func TestPreferredExt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ".jpg", photo.PreferredExt("jpeg"))
	assert.Equal(t, "", photo.PreferredExt("foo"))
}

func TestEncode(t *testing.T) {
	t.Parallel()

	img := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	buf := &bytes.Buffer{}

	err := photo.Encode("jpeg", buf, img, nil)
	assert.NoError(t, err)
	assert.Equal(t, "\xFF\xD8", string(buf.Bytes()[:2]))

	buf.Reset()

	err = photo.Encode("", buf, img, nil)
	assert.NoError(t, err)
	assert.Equal(t, "\xFF\xD8", string(buf.Bytes()[:2]))

	buf.Reset()

	err = photo.Encode("jpeg", buf, img, &jpeg.Options{
		Quality: 50,
	})
	assert.NoError(t, err)
	assert.Equal(t, "\xFF\xD8", string(buf.Bytes()[:2]))

	buf.Reset()

	err = photo.Encode("foo", buf, img, nil)
	assert.Error(t, err)
}

func TestRegisterEncoder(t *testing.T) {
	t.Parallel()

	photo.RegisterEncoder(
		"png",
		[]string{"plingplong", "png"},
		func(w io.Writer, img image.Image, opts any) error {
			return nil
		},
	)

	assert.Equal(t, "plingplong", photo.PreferredExt("png"))

	assert.Panics(t, func() {
		photo.RegisterEncoder("", nil, nil)
	})

	assert.Panics(t, func() {
		photo.RegisterEncoder("png", nil, nil)
	})

	assert.Panics(t, func() {
		photo.RegisterEncoder("png", []string{},
			func(w io.Writer, img image.Image, opts any) error {
				return nil
			},
		)
	})

	assert.Panics(t, func() {
		photo.RegisterEncoder("png", []string{"png"}, nil)
	})
}

func TestFilenameWithoutFormat(t *testing.T) {
	t.Parallel()

	p := photo.FilenameWithoutFormat("jpeg", "foo.jpeg")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("jpeg", "foo.jpg")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("jpeg", "foo.JPG")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("jpeg", "foo.gif")
	assert.Equal(t, "foo.gif", p)

	p = photo.FilenameWithoutFormat("jpeg", "foo")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("", "foo.jpeg")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("", "foo.jpg")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("", "foo.JPG")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("", "foo.gif")
	assert.Equal(t, "foo.gif", p)

	p = photo.FilenameWithoutFormat("", "foo")
	assert.Equal(t, "foo", p)

	p = photo.FilenameWithoutFormat("foo", "foo.jpeg")
	assert.Equal(t, "foo.jpeg", p)

	p = photo.FilenameWithoutFormat("foo", "foo.jpg")
	assert.Equal(t, "foo.jpg", p)

	p = photo.FilenameWithoutFormat("foo", "foo.JPG")
	assert.Equal(t, "foo.JPG", p)

	p = photo.FilenameWithoutFormat("foo", "foo.gif")
	assert.Equal(t, "foo.gif", p)

	p = photo.FilenameWithoutFormat("foo", "foo")
	assert.Equal(t, "foo", p)
}

func TestFilenameForFormat(t *testing.T) {
	t.Parallel()

	p, err := photo.FilenameForFormat("jpeg", "foo.jpeg")
	assert.NoError(t, err)
	assert.Equal(t, "foo.jpeg", p)

	p, err = photo.FilenameForFormat("jpeg", "foo.jpg")
	assert.NoError(t, err)
	assert.Equal(t, "foo.jpg", p)

	p, err = photo.FilenameForFormat("jpeg", "foo.JPG")
	assert.NoError(t, err)
	assert.Equal(t, "foo.JPG", p)

	// It's debatable whether this is a bug or a feature, but this is the
	// expected behavior, so...
	p, err = photo.FilenameForFormat("jpeg", "foo.gif")
	assert.NoError(t, err)
	assert.Equal(t, "foo.gif.jpg", p)

	p, err = photo.FilenameForFormat("", "foo.jpeg")
	assert.NoError(t, err)
	assert.Equal(t, "foo.jpeg", p)

	p, err = photo.FilenameForFormat("", "foo.jpg")
	assert.NoError(t, err)
	assert.Equal(t, "foo.jpg", p)

	p, err = photo.FilenameForFormat("", "foo.JPG")
	assert.NoError(t, err)
	assert.Equal(t, "foo.JPG", p)

	// It's debatable whether this is a bug or a feature, but this is the
	// expected behavior, so...
	p, err = photo.FilenameForFormat("", "foo.gif")
	assert.NoError(t, err)
	assert.Equal(t, "foo.gif.jpg", p)

	_, err = photo.FilenameForFormat("foo", "foo.jpg")
	assert.Error(t, err)
}

func TestCompleteImage_ImageReader(t *testing.T) {
	t.Parallel()

	fImg := photo.NewFile("unsplash/testdata/waa.jpg")
	cImg := photo.CompleteImage(fImg)

	assert.Equal(t, "unsplash/testdata/waa.jpg", cImg.Filename())

	r, err := cImg.Reader()
	assert.NoError(t, err)

	bs, err := io.ReadAll(r)
	assert.NoError(t, err)

	assert.Equal(t, "\xFF\xD8", string(bs[:2]))

	err = r.Close()
	assert.NoError(t, err)

	i, ext, err := cImg.Image()
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", ext)
	assert.Equal(t, image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{4128, 2322},
	}, i.Bounds())
}

func TestCompleteImage_ImageDecoded(t *testing.T) {
	t.Parallel()

	f, err := os.Open("unsplash/testdata/waa.jpg")
	assert.NoError(t, err)

	decImg, err := jpeg.Decode(f)
	assert.NoError(t, err)

	err = f.Close()
	assert.NoError(t, err)

	fImg := photo.NewMemory("unsplash/testdata/waa.jpg", "jpeg", decImg)
	cImg := photo.CompleteImage(fImg)

	assert.Equal(t, "unsplash/testdata/waa.jpg", cImg.Filename())

	r, err := cImg.Reader()
	assert.NoError(t, err)

	bs, err := io.ReadAll(r)
	assert.NoError(t, err)

	assert.Equal(t, "\xFF\xD8", string(bs[:2]))

	err = r.Close()
	assert.NoError(t, err)

	i, ext, err := cImg.Image()
	assert.NoError(t, err)
	assert.Equal(t, "jpeg", ext)
	assert.Equal(t, image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{4128, 2322},
	}, i.Bounds())
}

type testImageComplete struct{}

func (t *testImageComplete) Filename() string {
	return ""
}

func (t *testImageComplete) Image() (img image.Image, format string, err error) {
	return nil, "", nil
}

func (t *testImageComplete) Reader() (io.ReadCloser, error) {
	return nil, nil
}

func TestCompleteImage_ImageComplete(t *testing.T) {
	t.Parallel()

	origImg := &testImageComplete{}
	cImg := photo.CompleteImage(origImg)

	assert.Equal(t, origImg, cImg)
}

type testImageReader struct {
	ReaderErr error
}

func (t *testImageReader) Filename() string {
	return ""
}

func (t *testImageReader) Reader() (io.ReadCloser, error) {
	return nil, t.ReaderErr
}

func TestCompleteImage_ImageReader_Sad(t *testing.T) {
	t.Parallel()

	origImg := &testImageReader{
		ReaderErr: assert.AnError,
	}
	cImg := photo.CompleteImage(origImg)

	_, _, err := cImg.Image()
	assert.ErrorIs(t, err, assert.AnError)
}

type testImageDecoded struct {
	ImageErr  error
	Encoding  string
	ImageData image.Image
}

func (t *testImageDecoded) Filename() string {
	return ""
}

func (t *testImageDecoded) Image() (img image.Image, format string, err error) {
	return t.ImageData, t.Encoding, t.ImageErr
}

func EncoderError(w io.Writer, img image.Image, opts any) error {
	return assert.AnError
}

func TestCompleteImage_ImageDecoded_Sad(t *testing.T) {
	t.Parallel()

	origImg := &testImageDecoded{
		ImageErr: assert.AnError,
	}

	cImg := photo.CompleteImage(origImg)

	_, err := cImg.Reader()
	assert.ErrorIs(t, err, assert.AnError)

	origImg.ImageErr = nil
	origImg.ImageData = image.NewNRGBA(image.Rect(0, 0, 100, 100))

	r, err := cImg.Reader()
	assert.NoError(t, err)
	assert.NotNil(t, r)

	origImg.Encoding = "foo"

	_, err = cImg.Reader()
	assert.Error(t, err)

	photo.RegisterEncoder("foo", []string{".foo"}, EncoderError)
	_, err = cImg.Reader()
	assert.ErrorIs(t, err, assert.AnError)
}
