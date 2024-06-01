package text

import (
	"html/template"
)

// Verse is the metadata and content for a verse of the day.
type Verse struct {
	Reference string  `yaml:"reference" json:"reference"`
	Content   Content `yaml:"content" json:"content"`
	Link      string  `yaml:"link,omitempty" json:"link,omitempty"`
	Version
}

// Content holds the content of a scripture of the day.
type Content map[string]string

// Text returns the plain text rendering of the scripture.
func (c Content) Text() string {
	return c["text"]
}

// HTML returns the HTML rendering of the scripture.
func (c Content) HTML() template.HTML {
	return template.HTML(c["html"])
}

// Version is the metadata for the version of the Bible used for the verse.
type Version struct {
	Name string `yaml:"name" json:"name"`
	Link string `yaml:"link" json:"link"`
}
