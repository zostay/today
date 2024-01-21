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

// CacheKey returns the cache key for a given photo URL.
func (u *Source) CacheKey(photoUrl string) (string, bool) {
	id, err := idFromUrl(photoUrl)
	if err != nil {
		return "", false
	}
	return "unsplash/" + id, true
}

// Photo returns the photo info for a given photo URL.
func (u *Source) Photo(
	ctx context.Context,
	photoUrl string,
) (*photo.Info, error) {
	photoId, err := idFromUrl(photoUrl)
	if err != nil {
		return nil, err
	}

	image, _, err := u.Client.Photos.Photo(photoId, nil)
	if err != nil {
		return nil, err
	}

	photoKey, _ := u.CacheKey(photoUrl)
	return &photo.Info{
		Key: photoKey,
		Meta: &photo.Meta{
			Link: urlValueString(image.Links.HTML),
			Type: "unsplash",
			Creator: photo.Creator{
				Name: stringValue(image.Photographer.Name),
				Link: urlValueString(image.Photographer.Links.HTML),
			},
		},
	}, nil
}

// Download fetches the photo for the photo info.
func (u *Source) Download(
	ctx context.Context,
	info *photo.Info,
) error {
	photoId, err := idFromUrl(info.Meta.Link)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "bg.*.jpg")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	dl, _, err := u.Client.Photos.DownloadLink(photoId)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Get(dl.String())
	if err != nil {
		return err
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	info.File = f

	return nil
}
