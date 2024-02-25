package unsplash

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/hbagdi/go-unsplash/unsplash"

	"github.com/zostay/today/pkg/photo"
)

// stringValue is a helper for use with the Source Client to pull out strings
// from responses.
func stringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// urlValueString is a helper for use with the Source Client to pull out URL
// strings from responses.
func urlValueString(u *unsplash.URL) string {
	if u == nil {
		return ""
	}
	return u.String()
}

// idFromUrl extracts the photo ID from a URL.
func idFromUrl(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	return u.Path[len(u.Path)-11:], nil
}

// Photo returns the photo info for a given photo URL.
func (u *Source) Photo(
	ctx context.Context,
	photoUrl string,
) (*photo.Descriptor, error) {
	photoId, err := idFromUrl(photoUrl)
	if err != nil {
		return nil, err
	}

	image, _, err := u.Client.Photos.Photo(photoId, nil)
	if err != nil {
		return nil, err
	}

	d := &photo.Descriptor{
		Link: urlValueString(image.Links.HTML),
		Type: "unsplash",
		Creator: photo.Creator{
			Name: stringValue(image.Photographer.Name),
			Link: urlValueString(image.Photographer.Links.HTML),
		},
	}

	filename, err := idFromUrl(urlValueString(image.Links.Download))
	if err != nil {
		return nil, err
	}

	dl, _, err := u.Client.Photos.DownloadLink(photoId)
	if err != nil {
		return nil, err
	}

	d.AddImage(photo.Original, &unsplashImage{
		filename: filename,
		link:     dl.String(),
	})

	return d, nil
}

type unsplashImage struct {
	filename string
	link     string
}

func (u *unsplashImage) Filename() string {
	return u.filename
}

func (u *unsplashImage) Reader() (io.ReadCloser, error) {
	f, err := os.CreateTemp("", "bg.*.jpg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	res, err := http.DefaultClient.Get(u.link)
	if err != nil {
		return nil, err
	}

	return res.Body, err
}

var (
	_ photo.Image       = (*unsplashImage)(nil)
	_ photo.ImageReader = (*unsplashImage)(nil)
)
