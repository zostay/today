package photo

import (
	"fmt"
	"image/color"
	"strconv"
)

const (
	Original = "original"
)

// Descriptor provides metadata data about the image, which is stored in a
// serialization format on openscripture.today, and also includes a reference to
// the image file.
type Descriptor struct {
	Link  string `yaml:"link" json:"link"`
	Type  string `yaml:"type" json:"type"`
	Title string `yaml:"title,omitempty" json:"title,omitempty"`
	Color string `yaml:"color,omitempty" json:"color,omitempty"`
	Creator

	images map[string]ImageComplete
}

// AddImage adds an image to the descriptor.
func (d *Descriptor) AddImage(key string, img Image) {
	if d.images == nil {
		d.images = map[string]ImageComplete{}
	}
	d.images[key] = CompleteImage(img)
}

// RemoveImage removes an image from the descriptor.
func (d *Descriptor) RemoveImage(key string) {
	delete(d.images, key)
}

// HasImage returns true if the descriptor has an image with the given key.
func (d *Descriptor) HasImage(key string) bool {
	_, ok := d.images[key]
	return ok
}

// GetImage returns an image from the descriptor.
func (d *Descriptor) GetImage(key string) ImageComplete {
	img := d.images[key]
	return img
}

// SetColor sets the color as a CSS string from a color.Color.
func (d *Descriptor) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	d.Color = fmt.Sprintf("#%02x%02x%02x", r*256/a, g*256/a, b*256/a)
}

// GetColor returns a color.Color after decoding the CSS string stored. Returns
// nil and an error if the color cannot be decoded. Returns nil with no error if
// the color is not set.
func (d *Descriptor) GetColor() (color.Color, error) {
	if d.Color == "" {
		return nil, nil
	}

	if d.Color[0] != '#' || len(d.Color) != 7 {
		return nil, fmt.Errorf("invalid color format")
	}

	r, g, b := d.Color[1:3], d.Color[3:5], d.Color[5:7]
	ri, err := strconv.ParseInt(r, 16, 8)
	if err != nil {
		return nil, err
	}

	gi, err := strconv.ParseInt(g, 16, 8)
	if err != nil {
		return nil, err
	}

	bi, err := strconv.ParseInt(b, 16, 8)
	if err != nil {
		return nil, err
	}

	rc, gc, bc := uint8(ri), uint8(gi), uint8(bi)
	return color.RGBA{rc, gc, bc, 255}, nil
}

// Creator contains the metadata about the creator of a photo.
type Creator struct {
	Name string `yaml:"name" json:"name"`
	Link string `yaml:"link" json:"link"`
}
