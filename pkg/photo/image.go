package photo

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const DefaultFormat = "jpeg"

type ImageEncoder struct {
	// Ext must be lowercase. The first Ext is the preferred extension for the
	// format.
	Ext []string

	// Encoder is a formatter for the image format.
	Encoder func(w io.Writer, img image.Image, opts any) error
}

func (enc ImageEncoder) PreferredExt() string {
	return enc.Ext[0]
}

var registeredEncoders = map[string]ImageEncoder{
	"jpeg": ImageEncoder{
		Ext: []string{"jpg", "jpeg"},
		Encoder: func(w io.Writer, img image.Image, opts any) error {
			jpegOpts := opts.(*jpeg.Options)
			return jpeg.Encode(w, img, jpegOpts)
		},
	},
}

func RegisterEncoder(format string, ext []string, enc ImageEncoder) {
	if format == "" {
		panic("photo: cannot register an encoder with an empty format")
	}

	if enc.Encoder == nil {
		panic("photo: cannot register an encoder with a nil encoder")
	}

	if len(ext) == 0 {
		panic("photo: cannot register an encoder with no extensions")
	}

	registeredEncoders[format] = enc
}

// Image is the interface for a photo image. Every Image must also implemented
// ImageDecoded or ImageReader.
type Image interface {
	// Filename is the filename for the image. This may include a directory path
	// as well. It may be useful as a cache key, so a unique name is preferred.
	Filename() string
}

// FilenameWithoutFOrmat returns a filename with the format extension removed.
func FilenameWithoutFormat(format, path string) string {
	pathExt := strings.ToLower(filepath.Ext(path))
	if pathExt == "" {
		return path
	}

	if format == "" {
		format = DefaultFormat
	}

	enc, hasEnc := registeredEncoders[format]
	if !hasEnc {
		return path
	}

	for _, e := range enc.Ext {
		if e == pathExt {
			return path[:len(path)-len(pathExt)]
		}
	}

	return path
}

// FilanemeForFormat returns a filename for the given format. If the path already
// has an extension for the format, the path is returned as is. Otherwise, the
// path is returned with the extension for the format appended. If the format is
// empty, the default format is assumed. If the format is not registered, an
// error is returned.
func FilenameForFormat(format, path string) (string, error) {
	if format == "" {
		format = DefaultFormat
	}

	enc, hasEnc := registeredEncoders[format]
	if !hasEnc {
		return "", fmt.Errorf("photo: no encoder registered for format %q", format)
	}

	pathExt := strings.ToLower(filepath.Ext(path))
	for _, e := range enc.Ext {
		if e == pathExt {
			return path, nil
		}
	}

	return path + "." + enc.PreferredExt(), nil
}

// Encode encodes the image to the given format and writes it to the writer. If
// the format is empty, the default format is assumed. If the format is not
// registered, an error is returned.
func Encode(format string, img image.Image, w io.Writer, opts any) error {
	if format == "" {
		format = DefaultFormat
	}

	enc, hasEnc := registeredEncoders[format]
	if !hasEnc {
		return fmt.Errorf("photo: no encoder registered for format %q", format)
	}

	return enc.Encoder(w, img, opts)
}

// ImageDecoded is the interface that provides the image data in decoded form.
type ImageDecoded interface {
	// Image returns the image and the format of the image. The format is the
	// name of the format, like "jpeg" or "png". It should match the registered
	// format used to decode the image. If the image was not decoded, the format
	// should be empty. An error should be returned if there's a problem getting
	// the iamge data.
	Image() (img image.Image, format string, err error)
}

// ImageReader is the interface that provides the image data via io.Reader. After
// reading the image data, the caller must close the reader.
type ImageReader interface {
	Reader() (io.ReadCloser, error)
}

// ImageComplete is the interface for images that implement Image, ImageReader,
// and ImageDecoded.
type ImageComplete interface {
	Image
	ImageReader
	ImageDecoded
}

// Complete provides a complete implementation of the Image interface,
// providing ImageDecoded and ImageReader.
type Complete struct {
	image Image
}

func CompleteImage(img Image) ImageComplete {
	if img, isComplete := img.(ImageComplete); isComplete {
		return img
	}

	return &Complete{
		image: img,
	}
}

func (i *Complete) Filename() string {
	return i.image.Filename()
}

func (i *Complete) Image() (img image.Image, format string, err error) {
	switch it := i.image.(type) {
	case ImageDecoded:
		return it.Image()
	case ImageReader:
		r, err := it.Reader()
		if err != nil {
			return nil, "", err
		}

		return image.Decode(r)
	}

	return nil, "", fmt.Errorf("photo: image type %T does not implement ImageDecoded or ImageReader", i.image)
}

func (i *Complete) Reader() (io.ReadCloser, error) {
	switch it := i.image.(type) {
	case ImageReader:
		return it.Reader()
	case ImageDecoded:
		img, format, err := it.Image()
		if err != nil {
			return nil, err
		}

		if format == "" {
			format = DefaultFormat
		}

		enc, hasEnc := registeredEncoders[format]
		if !hasEnc {
			return nil, fmt.Errorf("photo: no encoder registered for format %q", format)
		}

		buf := &bytes.Buffer{}
		enc.Encoder(buf, img, nil)

		return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
	}

	return nil, fmt.Errorf("photo: image type %T does not implement ImageDecoded or ImageReader", i.image)
}

var (
	_ Image        = (*Complete)(nil)
	_ ImageReader  = (*Complete)(nil)
	_ ImageDecoded = (*Complete)(nil)
)

type Memory struct {
	Name   string
	Format string
	Img    image.Image
}

func NewMemory(name string, format string, img image.Image) *Memory {
	return &Memory{
		Name:   name,
		Format: format,
		Img:    img,
	}
}

func (m *Memory) Filename() string {
	return m.Name
}

func (m *Memory) Image() (image.Image, string, error) {
	return m.Img, m.Format, nil
}

var (
	_ Image        = (*Memory)(nil)
	_ ImageDecoded = (*Memory)(nil)
)

type File struct {
	Path string
}

func NewFile(path string) *File {
	return &File{
		Path: path,
	}
}

func (f *File) Filename() string {
	return f.Path
}

func (f *File) Reader() (io.ReadCloser, error) {
	return os.Open(f.Path)
}

var (
	_ Image       = (*File)(nil)
	_ ImageReader = (*File)(nil)
)
