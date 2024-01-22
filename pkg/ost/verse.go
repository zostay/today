package ost

import (
	"html/template"

	"github.com/zostay/today/pkg/text"
)

// Verse is the metadata and content for a verse of the day.
type Verse struct {
	Reference string        `yaml:"reference"`
	Content   template.HTML `yaml:"content"`
	text.Version
}
