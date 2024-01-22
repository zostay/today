package unsplash_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zostay/today/pkg/photo/unsplash"
)

func TestSource_New(t *testing.T) {
	t.Parallel()

	cli := unsplash.New(context.Background(), &unsplash.Auth{
		AccessKey: "test-token",
	})

	require.NotNil(t, cli)
	assert.NotNil(t, cli.Client)
}

func TestSource_NewFromAuthFile(t *testing.T) {
	t.Parallel()

	cli, err := unsplash.NewFromAuthFile(context.Background(), "auth-test.json")
	assert.NoError(t, err)

	require.NotNil(t, cli)
	assert.NotNil(t, cli.Client)
}

func TestSource_NewFromEnvironment(t *testing.T) {
	t.Parallel()

	err := os.Setenv("UNSPLASH_API_TOKEN", "test-token")
	assert.NoError(t, err)

	cli, err := unsplash.NewFromEnvironment(context.Background())
	assert.NoError(t, err)

	require.NotNil(t, cli)
	assert.NotNil(t, cli.Client)
}
