package esv

import (
	"context"
	"fmt"
	"html/template"
	"net/url"

	"github.com/zostay/go-esv-api/pkg/esv"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

// VersionInformation returns the metadata for the ESV from esv.org.
func (r *Resolver) VersionInformation(context.Context) (*text.Version, error) {
	return &text.Version{
		Name: "ESV",
		Link: "https://www.esv.org/",
	}, nil
}

// Verse fetches a verse an associated metadata for the given reference.
func (r *Resolver) Verse(ctx context.Context, ref *ref.Resolved) (*text.Verse, error) {
	txt, err := r.VerseText(ctx, ref)
	if err != nil {
		return nil, err
	}

	html, err := r.VerseHTML(ctx, ref)
	if err != nil {
		return nil, err
	}

	path, err := url.JoinPath(
		"https://www.esv.org/",
		url.PathEscape(ref.Ref()),
	)
	if err != nil {
		return nil, err
	}

	vi, err := r.VersionInformation(ctx)
	if err != nil {
		return nil, err
	}

	vr, err := ref.CompactRef()
	if err != nil {
		return nil, err
	}

	return &text.Verse{
		Reference: vr,
		Content: text.Content{
			Text: txt,
			HTML: html,
		},
		Link:    path,
		Version: *vi,
	}, nil
}

// VerseText returns the text of the given verse reference using the ESV API.
func (r *Resolver) VerseText(ctx context.Context, ref *ref.Resolved) (string, error) {
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

// VerseHTML returns the HTML of the given verse reference using the ESV API.
func (r *Resolver) VerseHTML(ctx context.Context, ref *ref.Resolved) (template.HTML, error) {
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

	return template.HTML(tr.Passages[0]), nil //nolint:gosec // we trust the ESV API
}

var _ text.Resolver = (*Resolver)(nil)
