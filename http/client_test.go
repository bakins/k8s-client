package http_test

import (
	"os"
	"testing"

	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/require"
)

// create a test client based on env variables.
func testClient(t *testing.T) *http.Client {
	server := os.Getenv("K8S_SERVER")

	if server == "" {
		server = "http://127.0.0.1:8001"
	}

	c, err := http.New(http.SetServer(server))
	require.Nil(t, err)

	return c
}
