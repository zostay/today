package esv

import (
	"html/template"

	"github.com/zostay/go-esv-api/pkg/esv"

	"github.com/zostay/today/pkg/ref"
	"github.com/zostay/today/pkg/text"
)

// VersionInformation returns the metadata for the ESV from esv.org
func (r *Resolver) VersionInformation() (*text.Version, error) {
	return &text.Version{
		Name: "ESV",
		Link: "https://www.esv.org/",
	}, nil
}

// Verse returns the text of the given verse reference using the ESV API.
func (r *Resolver) Verse(ref *ref.Resolved) (string, error) {
	tr, err := r.Client.PassageText(ref.Ref(),
		esv.WithIncludeVerseNumbers(false),
		esv.WithIncludeHeadings(false),
		esv.WithIncludeFootnotes(false),
	)
	if err != nil {
		return "", err
	}

	return tr.Passages[0], nil
}

// VerseHTML returns the HTML of the given verse reference using the ESV API.
func (r *Resolver) VerseHTML(ref *ref.Resolved) (template.HTML, error) {
	tr, err := r.Client.PassageHtml(ref.Ref(),
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

	return template.HTML(tr.Passages[0]), nil //nolint:gosec // we trust the ESV API
}
