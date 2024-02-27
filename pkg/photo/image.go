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
	"sync"
	"sync/atomic"
)

const DefaultFormat = "jpeg"

type EncoderFunc func(w io.Writer, img image.Image, opts any) error
type imageEncoder struct {
	Name    string
	Ext     []string
	Encoder func(w io.Writer, img image.Image, opts any) error
}

var (
	formatLock         sync.Mutex
	registeredEncoders atomic.Value
)

func init() {
	RegisterEncoder(
		"jpeg",
		[]string{".jpg", ".jpeg"},
		func(w io.Writer, img image.Image, opts any) error {
			var jpegOpts *jpeg.Options
			if opts != nil && opts.(*jpeg.Options) != nil {
				jpegOpts = opts.(*jpeg.Options)
			}
			return jpeg.Encode(w, img, jpegOpts)
		},
	)
}

func RegisterEncoder(format string, ext []string, enc EncoderFunc) {
	if format == "" {
		panic("photo: cannot register an encoder with an empty format")
	}

	if enc == nil {
		panic("photo: cannot register an encoder with a nil encoder")
	}

	if len(ext) == 0 {
		panic("photo: cannot register an encoder with no extensions")
	}

	formatLock.Lock()
	defer formatLock.Unlock()
	formats, _ := registeredEncoders.Load().([]imageEncoder)
	registeredEncoders.Store(append(formats, imageEncoder{
		Name:    format,
		Ext:     ext,
		Encoder: enc,
	}))
}

func getEncoder(format string) (imageEncoder, bool) {
	formatLock.Lock()
	defer formatLock.Unlock()
	formats, _ := registeredEncoders.Load().([]imageEncoder)
	for i := len(formats) - 1; i >= 0; i-- {
		if formats[i].Name == format {
			return formats[i], true
		}
	}

	return imageEncoder{}, false
}

func PreferredExt(format string) string {
	enc, hasEnc := getEncoder(format)
	if !hasEnc {
		return ""
	}

	return enc.Ext[0]
}

// Image is the interface for a photo image. Every Image must also implemented
// ImageDecoded or ImageReader.
type Image interface {
	// Filename is the filename for the image. This may include a directory path
	// as well. It may be useful as a cache key, so a unique name is preferred.
	Filename() string
}

// FilenameWithoutFormat returns a filename with the format extension removed.
func FilenameWithoutFormat(format, path string) string {
	pathExt := strings.ToLower(filepath.Ext(path))
	if pathExt == "" {
		return path
	}

	if format == "" {
		format = DefaultFormat
	}

	enc, hasEnc := getEncoder(format)
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

// FilenameForFormat returns a filename for the given format. If the path already
// has an extension for the format, the path is returned as is. Otherwise, the
// path is returned with the extension for the format appended. If the format is
// empty, the default format is assumed. If the format is not registered, an
// error is returned.
func FilenameForFormat(format, path string) (string, error) {
	if format == "" {
		format = DefaultFormat
	}

	enc, hasEnc := getEncoder(format)
	if !hasEnc {
		return "", fmt.Errorf("photo: no encoder registered for format %q", format)
	}

	pathExt := strings.ToLower(filepath.Ext(path))
	for _, e := range enc.Ext {
		if e == pathExt {
			return path, nil
		}
	}

	return path + enc.Ext[0], nil
}

// Encode encodes the image to the given format and writes it to the writer. If
// the format is empty, the default format is assumed. If the format is not
// registered, an error is returned.
func Encode(format string, w io.Writer, img image.Image, opts any) error {
	if format == "" {
		format = DefaultFormat
	}

	enc, hasEnc := getEncoder(format)
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

func (i *Complete) Image() (image.Image, string, error) {
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

		enc, hasEnc := getEncoder(format)
		if !hasEnc {
			return nil, fmt.Errorf("photo: no encoder registered for format %q", format)
		}

		buf := &bytes.Buffer{}
		err = enc.Encoder(buf, img, nil)
		if err != nil {
			return nil, fmt.Errorf("photo: error encoding image to %q: %w", format, err)
		}

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
