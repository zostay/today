package esv

import (
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
