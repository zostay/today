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
	OriginalVersion int  `yaml:"original_version,omitempty" json:"original_version,omitempty"`
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
