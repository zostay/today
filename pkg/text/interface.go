package text

import (
	"context"
	"html/template"

	"github.com/zostay/today/pkg/ref"
)

// Resolver is the interface used to retrieve text for a scripture passage.
type Resolver interface {
	// Verse turns a reference into a string of text.
	Verse(ctx context.Context, ref *ref.Resolved) (string, error)

	// VerseHTML turns a reference into a string of HTML.
	VerseHTML(ctx context.Context, ref *ref.Resolved) (template.HTML, error)

	// VersionInformation returns the metadata for the version of the Bible
	// used for the verse.
	VersionInformation(ctx context.Context) (*Version, error)
}

// Version is the metadata for the version of the Bible used for the verse.
type Version struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}
