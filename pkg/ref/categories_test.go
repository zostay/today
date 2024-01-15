package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategories(t *testing.T) {
	t.Parallel()

	cat := Categories()
	assert.Len(t, cat, 7)
}
