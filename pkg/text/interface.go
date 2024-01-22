package text

import (
	"html/template"

	"github.com/zostay/today/pkg/ref"
)

// Resolver is the interface used to retrieve text for a scripture passage.
type Resolver interface {
	// Verse turns a reference into a string of text.
	Verse(ref *ref.Resolved) (string, error)

	// VerseHTML turns a reference into a string of HTML.
	VerseHTML(ref *ref.Resolved) (template.HTML, error)

	// VersionInformation returns the metadata for the version of the Bible
	// used for the verse.
	VersionInformation() (*Version, error)
}

// Version is the metadata for the version of the Bible used for the verse.
type Version struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}
