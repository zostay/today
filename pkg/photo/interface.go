package photo

import "context"

// Source is the interface for a source of photo information.
type Source interface {
	// Meta returns the photo info for a given URL.
	Photo(ctx context.Context, url string) (desc *Descriptor, err error)
}
