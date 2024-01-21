package photo

import "context"

// Source is the interface for a source of photo information.
type Source interface {
	// CacheKey returns the cache key for a given URL. If the source does not
	// have a value for the URL, it returns false.
	CacheKey(url string) (string, bool)

	// Meta returns the photo info for a given URL.
	Photo(ctx context.Context, url string) (info *Info, err error)

	// Download downloads the photo and attaches it to the given Info.
	Download(ctx context.Context, info *Info) error
}
