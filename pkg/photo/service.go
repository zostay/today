package photo

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"net/http"

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

// Meta returns the photo info for a given photo URL. It will cache the
// photo info in the local filesystem and in S3.
func (s *Service) Photo(
	ctx context.Context,
	photoUrl string,
) (*Descriptor, error) {
	d, err := s.Source.Photo(ctx, photoUrl)
	if err != nil {
		return nil, err
	}

	return d, nil
}

type options struct {
	maxWidth  uint
	maxHeight uint
	image     string
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

func FromImage(img string) Option {
	return func(o *options) {
		o.image = img
	}
}

func processOptions(opts []Option) *options {
	o := &options{
		maxWidth:  DefaultMaxWidth,
		maxHeight: DefaultMaxHeight,
		image:     Original,
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
	d *Descriptor,
	opts ...Option,
) (string, error) {
	o := processOptions(opts)

	if !d.HasImage(o.image) {
		return "", fmt.Errorf("the %s image is not found", o.image)
	}

	item := d.GetImage(o.image)

	img, format, err := item.Image()
	if err != nil {
		return "", err
	}

	rect := img.Bounds()

	wi := rect.Max.X - rect.Min.X
	if wi <= 0 {
		return "", fmt.Errorf("invalid image width: %d", wi)
	}
	w := uint(wi)

	hi := rect.Max.Y - rect.Min.Y
	if hi <= 0 {
		return "", fmt.Errorf("invalid image height: %d", hi)
	}
	h := uint(hi)

	var (
		key               = fmt.Sprintf("resize:%dx%d", w, h)
		resizeWidth  uint = 0
		resizeHeight uint = 0
	)
	switch {
	case w > o.maxWidth:
		resizeWidth = o.maxWidth
	case h > o.maxHeight:
		resizeHeight = o.maxHeight
	default:
		d.AddImage(key,
			NewMemory(
				FilenameWithoutFormat(format, item.Filename()),
				"",
				img,
			),
		)
		return key, nil
	}

	rImg := resize.Resize(resizeWidth, resizeHeight, img, resize.Bicubic)

	d.AddImage(key,
		NewMemory(
			FilenameWithoutFormat(format, item.Filename()),
			"",
			rImg,
		),
	)
	return key, nil
}

// DominantImageColor returns the dominant color of the image. This is calculated
// by creating a histogram of all image pixels, then selecting the color that is
// used the most after excluding pure black and white (unless they are the only
// possible choices).
func DominantImageColor(
	ctx context.Context,
	img image.Image,
) (color.Color, error) {
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

	// disquality black and white unless they might be the only colors
	if len(hist) > 2 {
		delete(hist, color.RGBA{0, 0, 0, 255})
		delete(hist, color.RGBA{255, 255, 255, 255})
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
