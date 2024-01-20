package unsplash

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hbagdi/go-unsplash/unsplash"
	"golang.org/x/oauth2"

	"github.com/zostay/today/pkg/image"
)

// Source sources photos from unsplash.com.
type Source struct {
	Client *unsplash.Unsplash
}

var _ image.Source = (*Source)(nil)

// unsplashClient creates a new Source Client.
func unsplashClient(
	ctx context.Context,
	accessKey string,
) *unsplash.Unsplash {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "Client-ID " + accessKey},
	)
	client := oauth2.NewClient(ctx, ts)

	return unsplash.New(client)
}

func New(
	ctx context.Context,
	auth *Auth,
) *Source {
	client := unsplashClient(ctx, auth.AccessKey)

	return &Source{
		Client: client,
	}
}

func NewFromAuthFile(
	ctx context.Context,
	path string,
) (*Source, error) {
	auth, err := LoadAuth(path)
	if err != nil {
		return nil, err
	}

	return New(ctx, auth), nil
}

func NewFromEnvironment(ctx context.Context) (*Source, error) {
	// try the environment first
	tok := os.Getenv("UNSPLASH_API_TOKEN")
	if tok != "" {
		return New(ctx, &Auth{AccessKey: tok}), nil
	}

	// try the auth file
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return NewFromAuthFile(ctx, filepath.Join(homePath, AuthFile))
}
