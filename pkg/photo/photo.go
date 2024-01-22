package photo

import "os"

// Info is the information about a photo. It combines the cache key, the
// loaded photo metadata, and the file handle to the JPEG.
type Info struct {
	// Key is a special value that is usually set.
	Key string

	// Meta is the photo metadata.
	*Meta

	// File, if not nil, holds a reference to a file handle open for reading the
	// image.
	*os.File
}

// HasPhoto returns true if the photo info has a downloaded file to work with.
func (pi *Info) HasDownload() bool {
	return pi.File != nil
}

// Close ensures the file handle is closed, if present. Should always be called
// when done with the photo info.
func (pi *Info) Close() error {
	if pi.File != nil {
		f := pi.File
		pi.File = nil
		return f.Close()
	}
	return nil
}

// Meta contains the metadata about a photo.
type Meta struct {
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
