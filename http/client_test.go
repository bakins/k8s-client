package http_test

import (
	"fmt"
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

	opts := []http.OptionsFunc{
		http.SetServer(server),
	}

	if caFile := os.Getenv("K8S_CAFILE"); caFile != "" {
		fmt.Println("adding ca file ", caFile)
		opts = append(opts, http.SetCAFromFile(caFile))
	}

	c, err := http.New(opts...)
	require.Nil(t, err)

	return c
}
