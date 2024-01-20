package unsplash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAuth(t *testing.T) {
	t.Parallel()

	a, err := LoadAuth("auth-test.json")
	assert.NoError(t, err)
	assert.Equal(t, "test-token", a.AccessKey)
}
