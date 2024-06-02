package ost

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/photo"
	"github.com/zostay/today/pkg/text"
)

type Metadata struct {
	Version         int  `yaml:"version" json:"version"`
	OriginalVersion int  `yaml:"original_version" json:"original_version"`
	Pruned          bool `yaml:"pruned,omitempty" json:"pruned,omitempty"`
}

type Photo struct {
	Metadata         Metadata `yaml:"metadata" json:"metadata"`
	photo.Descriptor `yaml:",inline" json:",inline"`
}

type Verse struct {
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	text.Verse `yaml:",inline" json:",inline"`
}

type IndexEntry struct {
	Reference string `yaml:"reference" json:"reference"`
}

type Index struct {
	Metadata    Metadata              `yaml:"metadata" json:"metadata"`
	Description string                `yaml:"description" json:"description"`
	Verses      map[string]IndexEntry `yaml:"verses" json:"verses"`
}

func LoadPhotoYaml(r io.Reader, p *Photo) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(p)
}

func LoadPhotoJson(r io.Reader, p *Photo) error {
	dec := json.NewDecoder(r)
	return dec.Decode(p)
}

func LoadVerseYaml(r io.Reader, v *Verse) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(v)
}

func LoadVerseJson(r io.Reader, v *Verse) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

func LoadIndexYaml(r io.Reader, i *Index) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(i)
}

func LoadIndexJson(r io.Reader, i *Index) error {
	dec := json.NewDecoder(r)
	return dec.Decode(i)
}
