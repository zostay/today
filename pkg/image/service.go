package image

import (
	"context"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

const (
	DefaultMaxWidth  uint = 3840
	DefaultMaxHeight uint = 3840
)

type Service struct {
	Source
	Client *http.Client
}

func NewService(s Source) *Service {
	return &Service{
		Source: s,
		Client: http.DefaultClient,
	}
}

// Photo returns the photo info for a given photo URL. It will cache the
// photo info in the local filesystem and in S3.
func (s *Service) Photo(
	ctx context.Context,
	photoUrl string,
) (*PhotoInfo, error) {
	pi, err := s.Source.Photo(ctx, photoUrl)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

// Download fetches file.
func (s *Service) Download(
	ctx context.Context,
	info *PhotoInfo,
) error {
	return s.Source.Download(ctx, info)
}

type options struct {
	maxWidth  uint
	maxHeight uint
}

type Option func(*options)

func MaxWidth(w uint) Option {
	return func(o *options) {
		o.maxWidth = w
	}
}

func MaxHeight(h uint) Option {
	return func(o *options) {
		o.maxHeight = h
	}
}

func processOptions(opts []Option) *options {
	o := &options{
		maxWidth:  DefaultMaxWidth,
		maxHeight: DefaultMaxHeight,
	}
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// ResizedImage returns a resized version of the photo if it is larger than
// the maximum width and height I have set.
func (s *Service) ResizedImage(
	ctx context.Context,
	info *PhotoInfo,
	opts ...Option,
) (*os.File, error) {
	o := processOptions(opts)

	if info.HasDownload() {
		err := s.Source.Download(ctx, info)
		if err != nil {
			return nil, err
		}
	}

	img, err := jpeg.Decode(info.File)
	if err != nil {
		return nil, err
	}

	_, err = info.File.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	rect := img.Bounds()
	w := uint(rect.Max.X - rect.Min.X)
	h := uint(rect.Max.Y - rect.Min.Y)

	var (
		resizeWidth, resizeHeight uint = 0, 0
	)
	if w > o.maxWidth {
		resizeWidth = o.maxWidth
	} else if h > o.maxHeight {
		resizeHeight = o.maxHeight
	} else {
		return info.File, nil
	}

	rImg := resize.Resize(resizeWidth, resizeHeight, img, resize.Bicubic)

	tmpJr, err := os.CreateTemp("", "bg.*.jpg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpJr.Name())

	err = jpeg.Encode(tmpJr, rImg, nil)
	if err != nil {
		return nil, err
	}

	_, err = tmpJr.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return tmpJr, nil
}

// DominantImageColor returns a resized version of the photo if it is larger than
// the maximum width and height I have set.
func (s *Service) DominantImageColor(
	ctx context.Context,
	photo *PhotoInfo,
) (color.Color, error) {
	if photo.HasDownload() {
		err := s.Source.Download(ctx, photo)
		if err != nil {
			return nil, err
		}
	}

	img, err := jpeg.Decode(photo.File)
	if err != nil {
		return nil, err
	}

	_, err = photo.File.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	hist := make(map[color.RGBA]uint32)

	rect := img.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			r := uint8(256 * float64(r1) / float64(a1))
			g := uint8(256 * float64(g1) / float64(a1))
			b := uint8(256 * float64(b1) / float64(a1))

			hist[color.RGBA{r, g, b, 255}]++
		}
	}

	var best color.RGBA
	var bestCount uint32
	for color, count := range hist {
		if count > bestCount {
			best = color
			bestCount = count
		}
	}

	return best, nil
}
