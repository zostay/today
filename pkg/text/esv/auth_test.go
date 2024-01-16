package esv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/text/esv"
)

func TestLoadAuth(t *testing.T) {
	t.Parallel()

	a, err := esv.LoadAuth("auth-test.json")
	assert.NoError(t, err)
	assert.Equal(t, "test-token", a.AccessKey)
}
