package esv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zostay/today/pkg/text/esv"
)

func TestResolver_New(t *testing.T) {
	t.Parallel()

	cli := esv.New(&esv.Auth{
		AccessKey: "test-token",
	})

	assert.NotNil(t, cli)
	assert.Equal(t, "test-token", cli.Client.Token)
}

func TestResolver_NewFromAuthFile(t *testing.T) {
	t.Parallel()

	cli, err := esv.NewFromAuthFile("auth-test.json")
	assert.NoError(t, err)

	assert.NotNil(t, cli)
	assert.Equal(t, "test-token", cli.Client.Token)
}

func TestResolver_NewFromEnvironment(t *testing.T) {
	t.Parallel()

	err := os.Setenv("ESV_API_TOKEN", "test-token")
	assert.NoError(t, err)

	cli, err := esv.NewFromEnvironment()
	assert.NoError(t, err)

	assert.NotNil(t, cli)
	assert.Equal(t, "test-token", cli.Client.Token)
}
