package ost

import "html/template"

// Verse is the metadata and content for a verse of the day.
type Verse struct {
	Reference string        `yaml:"reference"`
	Content   template.HTML `yaml:"content"`
	Version
}

// Version is the metadata for the version of the Bible used for the verse.
type Version struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}
