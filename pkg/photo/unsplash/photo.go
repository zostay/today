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
// hotlinkURL picks the URL to display the photo from. Unsplash asks that photos
// be shown from the URLs it returns under "urls" rather than from a copy the
// consumer hosts. Raw is preferred because it carries no size preset, leaving
// the caller free to append its own Imgix parameters; full and regular stand in
// when a response omits it.
func hotlinkURL(image *unsplash.Photo) string {
	if image.Urls == nil {
		return ""
	}

	for _, u := range []*unsplash.URL{image.Urls.Raw, image.Urls.Full, image.Urls.Regular} {
		if s := urlValueString(u); s != "" {
			return s
		}
	}

	return ""
}

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

// IDFromURL extracts the photo ID from a URL.
func IDFromURL(s string) (string, error) {
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
	photoId, err := IDFromURL(photoUrl)
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
		ImageURL:         hotlinkURL(image),
		DownloadLocation: urlValueString(image.Links.DownloadLocation),
	}

	filename, err := IDFromURL(urlValueString(image.Links.Download))
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
