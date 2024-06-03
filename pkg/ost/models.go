package ost

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/photo"
	"github.com/zostay/today/pkg/text"
)

// Metadata describes general information about the documents published in the
// openscripture.net API.
type Metadata struct {
	Version         int  `yaml:"version" json:"version"`
	OriginalVersion int  `yaml:"original_version" json:"original_version"`
	Pruned          bool `yaml:"pruned,omitempty" json:"pruned,omitempty"`
}

// Photo describes the photos paired with scripture on openscripture.net.
type Photo struct {
	Metadata         Metadata `yaml:"metadata" json:"metadata"`
	photo.Descriptor `yaml:",inline" json:",inline"`
}

// Verse describes the verses paired with scripture on openscripture.net.
type Verse struct {
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	text.Verse `yaml:",inline" json:",inline"`
}

// IndexEntry describes an entry in the scripture index of openscripture.net.
type IndexEntry struct {
	Reference string `yaml:"reference" json:"reference"`
}

// Index describes a scripture index on openscripture.net.
type Index struct {
	Metadata    Metadata              `yaml:"metadata" json:"metadata"`
	Description string                `yaml:"description" json:"description"`
	Verses      map[string]IndexEntry `yaml:"verses" json:"verses"`
}

// LoadPhotoYaml loads a photo descriptor file in YAML format.
func LoadPhotoYaml(r io.Reader, p *Photo) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(p)
}

// LoadPhotoJson loads a photo descriptior file in JSON format.
func LoadPhotoJson(r io.Reader, p *Photo) error {
	dec := json.NewDecoder(r)
	return dec.Decode(p)
}

// LoadVerseYaml loads a verse reference file in YAML format.
func LoadVerseYaml(r io.Reader, v *Verse) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(v)
}

// LoadVerseJson loads a verse reference file in JSON format.
func LoadVerseJson(r io.Reader, v *Verse) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

// LoadIndexYaml loads a scripture index file in YAML format.
func LoadIndexYaml(r io.Reader, i *Index) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(i)
}

// LoadIndexJson loads a scripture index file in JSON format.
func LoadIndexJson(r io.Reader, i *Index) error {
	dec := json.NewDecoder(r)
	return dec.Decode(i)
}
