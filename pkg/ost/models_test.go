package ost_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zostay/today/pkg/ost"
	"github.com/zostay/today/pkg/photo"
	"github.com/zostay/today/pkg/text"
)

var (
	firstPhoto = &ost.Photo{
		Metadata: ost.Metadata{
			Version:         2,
			OriginalVersion: 0,
			Pruned:          true,
		},
		Descriptor: photo.Descriptor{
			Link:  "https://unsplash.com/photos/stones-beside-river-Bjq6Toa-K_I",
			Type:  "unsplash",
			Color: "#f8f8f8",
			Creator: photo.Creator{
				Name: "Robert Bye",
				Link: "https://unsplash.com/@robertbye",
			},
		},
	}
	firstVerse = &ost.Verse{
		Metadata: ost.Metadata{
			Version:         3,
			OriginalVersion: 0,
			Pruned:          true,
		},
		Verse: text.Verse{
			Reference: "Matthew 15:32-39",
			Link:      "https://www.esv.org/Matthew%2015:32-15:39",
			Version: text.Version{
				Name: "ESV",
				Link: "https://www.esv.org/",
			},
		},
	}
)

func TestLoadPhotoYaml(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(`metadata:
    version: 2
    original_version: 0
    pruned: true
link: https://unsplash.com/photos/stones-beside-river-Bjq6Toa-K_I
type: unsplash
color: '#f8f8f8'
creator:
    name: Robert Bye
    link: https://unsplash.com/@robertbye
`)

	var p ost.Photo
	err := ost.LoadPhotoYaml(r, &p)
	assert.NoError(t, err)
	assert.Equal(t, firstPhoto, &p)
}

func TestLoadPhotoJson(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(`{"metadata":{"version":2,"original_version":0,"pruned":true},"link":"https://unsplash.com/photos/stones-beside-river-Bjq6Toa-K_I","type":"unsplash","color":"#f8f8f8","creator":{"name":"Robert Bye","link":"https://unsplash.com/@robertbye"}}`)

	var p ost.Photo
	err := ost.LoadPhotoJson(r, &p)
	assert.NoError(t, err)
	assert.Equal(t, firstPhoto, &p)
}

func TestLoadVerseYaml(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(`metadata:
    version: 3
    original_version: 0
    pruned: true
reference: Matthew 15:32-39
content: {}
link: https://www.esv.org/Matthew%2015:32-15:39
version:
    name: ESV
    link: https://www.esv.org/
`)

	var v ost.Verse
	err := ost.LoadVerseYaml(r, &v)
	assert.NoError(t, err)
	assert.Equal(t, firstVerse, &v)
}

func TestLoadVerseJson(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(`{"metadata":{"version":3,"original_version":0,"pruned":true},"reference":"Matthew 15:32-39","content":{},"link":"https://www.esv.org/Matthew%2015:32-15:39","version":{"name":"ESV","link":"https://www.esv.org/"}}`)

	var v ost.Verse
	err := ost.LoadVerseJson(r, &v)
	assert.NoError(t, err)
	assert.Equal(t, firstVerse, &v)
}
