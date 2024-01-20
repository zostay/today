package image

import "context"

// Source is the interface for a source of photo information.
type Source interface {
	// CacheKey returns the cache key for a given URL. If the source does not
	// have a value for the URL, it returns false.
	CacheKey(url string) (string, bool)

	// Photo returns the photo info for a given URL.
	Photo(ctx context.Context, url string) (info *PhotoInfo, err error)

	// Download downloads the photo and attaches it to the given PhotoInfo.
	Download(ctx context.Context, info *PhotoInfo) error
}
