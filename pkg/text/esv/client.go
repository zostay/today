package esv

import (
	"os"
	"path/filepath"

	"github.com/zostay/go-esv-api/pkg/esv"
)

type Resolver struct {
	*esv.Client
}

func New(auth *Auth) *Resolver {
	return &Resolver{
		Client: esv.New(auth.AccessKey),
	}
}

func NewFromAuthFile(path string) (*Resolver, error) {
	auth, err := LoadAuth(path)
	if err != nil {
		return nil, err
	}

	return New(auth), nil
}

func NewFromEnvironment() (*Resolver, error) {
	// try the environment first
	tok := os.Getenv("ESV_API_TOKEN")
	if tok != "" {
		return New(&Auth{AccessKey: tok}), nil
	}

	// try the auth file
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return NewFromAuthFile(filepath.Join(homePath, AuthFile))
}
