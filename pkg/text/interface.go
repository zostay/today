package text

import (
	"html/template"
)

type Resolver interface {
	Verse(ref ref.ProperRef) (string, error)
	VerseHTML(ref ref.ProperRef) (template.HTML, error)
}
