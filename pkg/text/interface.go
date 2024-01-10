package text

import (
	"html/template"

	"github.com/zostay/today/pkg/ref"
)

type Resolver interface {
	Verse(ref ref.Resolved) (string, error)
	VerseHTML(ref ref.Resolved) (template.HTML, error)
}
