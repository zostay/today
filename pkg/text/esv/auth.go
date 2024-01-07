package esv

import (
	"os"

	"gopkg.in/yaml.v3"
)

// AuthFile is the location where the local project should stash the ESV
// access key file.
const AuthFile = ".esv.yaml"

// Auth is the structure of the ESV access key file.
type Auth struct {
	// AccessKey is the ESV API access key.
	AccessKey string `yaml:"access_key"`
}

// LoadAuth loads the ESV access key file.
func LoadAuth(path string) (*Auth, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var esvAuth Auth
	err = yaml.NewDecoder(r).Decode(&esvAuth)
	return &esvAuth, err
}
