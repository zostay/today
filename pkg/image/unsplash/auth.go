package unsplash

import (
	"os"

	"gopkg.in/yaml.v3"
)

const AuthFile = ".unsplash.yaml"

// Auth is the structure of the Source authentication file.
type Auth struct {
	AccessKey string `yaml:"access_key"`
}

// loadUnsplashAuth loads the Source authentication file.
func LoadAuth(path string) (*Auth, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var unsplashAuth Auth
	err = yaml.NewDecoder(r).Decode(&unsplashAuth)
	return &unsplashAuth, err
}
