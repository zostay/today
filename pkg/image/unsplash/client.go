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
	client *unsplash.Unsplash
}

var _ image.Source = (*Source)(nil)

// unsplashClient creates a new Source client.
func unsplashClient(
	ctx context.Context,
	accessKey string,
) (*unsplash.Unsplash, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "Client-ID " + accessKey},
	)
	client := oauth2.NewClient(ctx, ts)

	return unsplash.New(client), nil
}

func New(
	ctx context.Context,
	auth *Auth,
) (*Source, error) {
	client, err := unsplashClient(ctx, auth.AccessKey)
	if err != nil {
		return nil, err
	}

	return &Source{
		client: client,
	}, nil
}

func NewFromAuthFile(
	ctx context.Context,
	path string,
) (*Source, error) {
	auth, err := LoadAuth(path)
	if err != nil {
		return nil, err
	}

	return New(ctx, auth)
}

func NewFromEnvironment(ctx context.Context) (*Source, error) {
	// try the environment first
	tok := os.Getenv("ESV_API_TOKEN")
	if tok != "" {
		return New(ctx, &Auth{AccessKey: tok})
	}

	// try the auth file
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return NewFromAuthFile(ctx, filepath.Join(homePath, AuthFile))
}
