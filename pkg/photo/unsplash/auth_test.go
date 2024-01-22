package unsplash_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/photo/unsplash"
)

func TestLoadAuth(t *testing.T) {
	t.Parallel()

	a, err := unsplash.LoadAuth("auth-test.json")
	assert.NoError(t, err)
	assert.Equal(t, "test-token", a.AccessKey)
}
