package image

import "os"

// PhotoInfo is the information about a photo. It combines the cache key, the
// loaded photo metadata, and the file handle to the JPEG.
type PhotoInfo struct {
	Key string
	*Photo
	*os.File
}

// HasPhoto returns true if the photo info has a downloaded file to work with.
func (pi *PhotoInfo) HasDownload() bool {
	return pi.File != nil
}

// Close ensures the file handle is closed, if present. Should always be called
// when done with the photo info.
func (pi *PhotoInfo) Close() error {
	if pi.File != nil {
		f := pi.File
		pi.File = nil
		return f.Close()
	}
	return nil
}

// Photo contains the metadata about a photo.
type Photo struct {
	Link  string `yaml:"link"`
	Type  string `yaml:"type"`
	Title string `yaml:"title,omitempty"`
	Creator
}

// Creator contains the metadata about the creator of a photo.
type Creator struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}
