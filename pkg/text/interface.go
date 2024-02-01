package text

import (
	"context"
	"html/template"

	"github.com/zostay/today/pkg/ref"
)

// Resolver is the interface used to retrieve text for a scripture passage.
type Resolver interface {
	// Verse fetches a verse an associated metadata for the given reference..
	Verse(ctx context.Context, ref *ref.Resolved) (*Verse, error)

	// VerseText turns a reference into a string of text.
	VerseText(ctx context.Context, ref *ref.Resolved) (string, error)

	// VerseHTML turns a reference into a string of HTML.
	VerseHTML(ctx context.Context, ref *ref.Resolved) (template.HTML, error)

	// VersionInformation returns the metadata for the version of the Bible
	// used for the verse.
	VersionInformation(ctx context.Context) (*Version, error)
}

// Verse is the metadata and content for a verse of the day.
type Verse struct {
	Reference string  `yaml:"reference"`
	Content   Content `yaml:"content"`
	Link      string  `yaml:"link,omitempty"`
	Version
}

// Content holds the content of a scripture of the day.
type Content struct {
	Text string        `yaml:"text,omitempty"`
	HTML template.HTML `yaml:"html,omitempty"`
}

// Version is the metadata for the version of the Bible used for the verse.
type Version struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}
