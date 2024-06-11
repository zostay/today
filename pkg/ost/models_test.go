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
	firstIndex = &ost.Index{
		Metadata: ost.Metadata{
			Version:         0,
			OriginalVersion: 0,
		},
		Description: "Index of verses for 2023/12",
		Verses: map[string]ost.IndexEntry{
			"2023/12/20": {
				Reference: "Matthew 15:32-39",
			},
			"2023/12/21": {
				Reference: "Romans 8:31-39",
			},
			"2023/12/22": {
				Reference: "John 3:1-8",
			},
			"2023/12/23": {
				Reference: "Jude 17-23",
			},
			"2023/12/24": {
				Reference: "Ezekiel 36:22-38",
			},
			"2023/12/25": {
				Reference: "Genesis 3:14-24",
			},
			"2023/12/26": {
				Reference: "Luke 8:4-8",
			},
			"2023/12/27": {
				Reference: "Psalms 5",
			},
			"2023/12/28": {
				Reference: "Proverbs 16:4-9",
			},
			"2023/12/29": {
				Reference: "Revelation 4:1-4",
			},
			"2023/12/30": {
				Reference: "Daniel 2:20-21",
			},
			"2023/12/31": {
				Reference: "Genesis 39:19-23",
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

func TestLoadIndexYaml(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(`metadata:
    version: 0
    original_version: 0
description: Index of verses for 2023/12
verses:
    2023/12/20:
        reference: Matthew 15:32-39
    2023/12/21:
        reference: Romans 8:31-39
    2023/12/22:
        reference: John 3:1-8
    2023/12/23:
        reference: Jude 17-23
    2023/12/24:
        reference: Ezekiel 36:22-38
    2023/12/25:
        reference: Genesis 3:14-24
    2023/12/26:
        reference: Luke 8:4-8
    2023/12/27:
        reference: Psalms 5
    2023/12/28:
        reference: Proverbs 16:4-9
    2023/12/29:
        reference: Revelation 4:1-4
    2023/12/30:
        reference: Daniel 2:20-21
    2023/12/31:
        reference: Genesis 39:19-23

`)

	var idx ost.Index
	err := ost.LoadIndexYaml(r, &idx)
	assert.NoError(t, err)
	assert.Equal(t, firstIndex, &idx)
}

func TestLoadIndexJson(t *testing.T) {
	t.Parallel()

	r := strings.NewReader(`{"metadata":{"version":0,"original_version":0},"description":"Index of verses for 2023/12","verses":{"2023/12/20":{"reference":"Matthew 15:32-39"},"2023/12/21":{"reference":"Romans 8:31-39"},"2023/12/22":{"reference":"John 3:1-8"},"2023/12/23":{"reference":"Jude 17-23"},"2023/12/24":{"reference":"Ezekiel 36:22-38"},"2023/12/25":{"reference":"Genesis 3:14-24"},"2023/12/26":{"reference":"Luke 8:4-8"},"2023/12/27":{"reference":"Psalms 5"},"2023/12/28":{"reference":"Proverbs 16:4-9"},"2023/12/29":{"reference":"Revelation 4:1-4"},"2023/12/30":{"reference":"Daniel 2:20-21"},"2023/12/31":{"reference":"Genesis 39:19-23"}}}`)

	var idx ost.Index
	err := ost.LoadIndexJson(r, &idx)
	assert.NoError(t, err)
	assert.Equal(t, firstIndex, &idx)
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
