package output

import (
	"errors"

	"github.com/zostay/go-std/maps"
)

// ErrUnknownFormat is the error returned when an unknown output format is
// provided.q
var ErrUnknownFormat = errors.New("unknown output format")

// Format describes each format.
type Format struct {
	Name        string
	Description string
	Ext         string
}

const (
	TextFormat = "text"
	HTMLFormat = "html"
	JPEGFormat = "jpeg"
	YAMLFormat = "yaml"
	JSONFormat = "json"
	MetaFormat = "meta"
)

var formatInfo = map[string]Format{
	TextFormat: {
		Name:        TextFormat,
		Description: "Plain text output",
		Ext:         "txt",
	},
	HTMLFormat: {
		Name:        HTMLFormat,
		Description: "HTML output",
		Ext:         "html",
	},
	JPEGFormat: {
		Name:        JPEGFormat,
		Description: "Output as a JPEG image (meme-style)",
		Ext:         "jpg",
	},
	YAMLFormat: {
		Name:        YAMLFormat,
		Description: "YAML output",
		Ext:         "yaml",
	},
	JSONFormat: {
		Name:        JSONFormat,
		Description: "JSON output",
		Ext:         "json",
	},
	MetaFormat: {
		Name:        MetaFormat,
		Description: "Output metadata about the verse",
		Ext:         "txt",
	},
}

// DefaultFormat returns the text format.
func DefaultFormat() Format {
	return formatInfo["text"]
}

// IsKnownFormat returns true if the given format is a known output format.
func LookupFormat(f string) (Format, error) {
	if fmt, known := formatInfo[f]; known {
		return fmt, nil
	}
	return Format{}, ErrUnknownFormat
}

// ListFormats returns a list of all known output formats.
func ListFormats() []string {
	return maps.Keys(formatInfo)
}
