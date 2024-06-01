package text

import (
	"context"
	"errors"

	"github.com/zostay/today/pkg/ref"
)

var ErrUnsupportedFormat = errors.New("unsupported verse format")

// VerseFormat describes a format from the resolver.
type VerseFormat interface {
	// Name of the format. Should be the same name used to retrieve it from
	// DescribeFormat and passed to VerseAs.
	Name() string

	// Description is a human-readable description of the format.
	Description() string

	// Ext is the file extension to use for the format.
	Ext() string
}

type Resolver interface {
	// AsFormats returns the list of formats that this resolver can provide.
	// Ideally, at least "text" and "html" should always be supported.
	AsFormats() []string

	// DescribeFormat returns a description of each possible format.
	DescribeFormat(ofmt string) VerseFormat

	// VerseAs turns a reference into a string of text in the requested format.
	VerseAs(ctx context.Context, rs *ref.Resolved, ofmt string) (string, error)

	// VerseURI turns a reference into a permalink for the resolver. Ideally,
	// this is a URL reachable via a browser, but could just be a URI describing
	// the verse for this resolver.
	VerseURI(ctx context.Context, rs *ref.Resolved) (string, error)

	// VersionInformation returns the metadata for the version of the Bible
	// used for the verse.
	VersionInformation(ctx context.Context) (*Version, error)
}
