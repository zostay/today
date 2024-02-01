package photo

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
)

// Info is the information about a photo. It combines the cache key, the
// loaded photo metadata, and the file handle to the JPEG.
type Info struct {
	// Key is a special value that is usually set.
	Key string

	// Meta is the photo metadata.
	*Meta

	// File, if not nil, holds a reference to a file handle open for reading the
	// image.
	*os.File
}

// HasPhoto returns true if the photo info has a downloaded file to work with.
func (pi *Info) HasDownload() bool {
	return pi.File != nil
}

// Close ensures the file handle is closed, if present. Should always be called
// when done with the photo info.
func (pi *Info) Close() error {
	if pi.File != nil {
		f := pi.File
		pi.File = nil
		return f.Close()
	}
	return nil
}

// Meta contains the metadata about a photo.
type Meta struct {
	Link  string `yaml:"link"`
	Type  string `yaml:"type"`
	Title string `yaml:"title,omitempty"`
	Color string `yaml:"color,omitempty"`
	Creator
}

// SetColor sets the color as a CSS string from a color.Color.
func (m *Meta) SetColor(c color.Color) {
	r, g, b, a := c.RGBA()
	m.Color = fmt.Sprintf("#%02x%02x%02x", r*256/a, g*256/a, b*256/a)
}

// GetColor returns a color.Color after decoding the CSS string stored. Returns
// nil and an error if the color cannot be decoded. Returns nil with no error if
// the color is not set.
func (m *Meta) GetColor() (color.Color, error) {
	if m.Color == "" {
		return nil, nil
	}

	if m.Color[0] != '#' || len(m.Color) != 7 {
		return nil, fmt.Errorf("invalid color format")
	}

	r, g, b := m.Color[1:3], m.Color[3:5], m.Color[5:7]
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
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}
