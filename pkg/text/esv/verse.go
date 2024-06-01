package esv

import (
	"context"
	"fmt"
	"net/url"

	"github.com/zostay/go-esv-api/pkg/esv"
	"github.com/zostay/go-std/maps"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

type formatInfo struct {
	name string
	desc string
	ext  string
}

func (fi formatInfo) Name() string {
	return fi.name
}

func (fi formatInfo) Description() string {
	return fi.desc
}

func (fi formatInfo) Ext() string {
	return fi.ext
}

var formats = map[string]formatInfo{
	"text": {
		name: "text",
		desc: "Plain text",
		ext:  "txt",
	},
	"html": {
		name: "html",
		desc: "HTML",
		ext:  "html",
	},
}

// AsFormats returns "text" and "html".
func (r *Resolver) AsFormats() []string {
	return maps.Keys(formats)
}

// DescribeFormat returns the format information for the given format.
func (r *Resolver) DescribeFormat(fmt string) text.VerseFormat {
	fi, ok := formats[fmt]
	if !ok {
		return nil
	}

	return fi
}

// VersionInformation returns the metadata for the ESV from esv.org.
func (r *Resolver) VersionInformation(context.Context) (*text.Version, error) {
	return &text.Version{
		Name: "ESV",
		Link: "https://www.esv.org/",
	}, nil
}

// VerseAs returns the text of the given verse reference in the requested format.
func (r *Resolver) VerseAs(
	ctx context.Context,
	rs *ref.Resolved,
	ofmt string,
) (string, error) {
	switch ofmt {
	case "text":
		return r.verseText(ctx, rs)
	case "html":
		return r.verseHTML(ctx, rs)
	}
	return "", text.ErrUnsupportedFormat
}

// VerseURI returns the permalink to the esv.org website for the verse reference.
func (r *Resolver) VerseURI(ctx context.Context, rs *ref.Resolved) (string, error) {
	return url.JoinPath(
		"https://www.esv.org/",
		url.PathEscape(rs.Ref()),
	)
}

//// Verse fetches a verse an associated metadata for the given reference.
//func (r *Resolver) Verse(ctx context.Context, ref *ref.Resolved) (*text.Verse, error) {
//	txt, err := r.verseText(ctx, ref)
//	if err != nil {
//		return nil, err
//	}
//
//	html, err := r.verseHTML(ctx, ref)
//	if err != nil {
//		return nil, err
//	}
//
//	path, err := url.JoinPath(
//		"https://www.esv.org/",
//		url.PathEscape(ref.Ref()),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	vi, err := r.VersionInformation(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	vr, err := ref.CompactRef()
//	if err != nil {
//		return nil, err
//	}
//
//	return &text.Verse{
//		Reference: vr,
//		Content: text.Content{
//			Text: txt,
//			HTML: html,
//		},
//		Link:    path,
//		Version: *vi,
//	}, nil
//}

// verseText returns the text of the given verse reference using the ESV API.
func (r *Resolver) verseText(ctx context.Context, ref *ref.Resolved) (string, error) {
	tr, err := r.Client.PassageTextContext(ctx, ref.Ref(),
		esv.WithIncludeVerseNumbers(false),
		esv.WithIncludeHeadings(false),
		esv.WithIncludeFootnotes(false),
		esv.WithIncludePassageReferences(false),
	)
	if err != nil {
		return "", err
	}

	if len(tr.Passages) != 1 {
		return "", fmt.Errorf("expected a single passage returned but ESV API returned %d: %v", len(tr.Passages), tr)
	}

	return tr.Passages[0], nil
}

// verseHTML returns the HTML of the given verse reference using the ESV API.
func (r *Resolver) verseHTML(ctx context.Context, ref *ref.Resolved) (string, error) {
	tr, err := r.Client.PassageHtmlContext(ctx, ref.Ref(),
		esv.WithIncludeVerseNumbers(false),
		esv.WithIncludeHeadings(false),
		esv.WithIncludeFootnotes(false),
		esv.WithIncludeChapterNumbers(false),
		esv.WithIncludeAudioLink(false),
		esv.WithIncludeBookTitles(false),
		esv.WithIncludePassageReferences(false),
		esv.WithIncludeFirstVerseNumbers(false),
	)
	if err != nil {
		return "", err
	}

	if len(tr.Passages) != 1 {
		return "", fmt.Errorf("expected a single passage returned but ESV API returned %d: %v", len(tr.Passages), tr)
	}

	return tr.Passages[0], nil //nolint:gosec // we trust the ESV API
}

var _ text.Resolver = (*Resolver)(nil)
