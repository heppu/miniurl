//go:build integration

package miniurl_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/heppu/miniurl"
	"github.com/heppu/miniurl/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration(t *testing.T) {
	setupTestEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	}))
	t.Cleanup(srv.Close)

	errCh := make(chan error, 1)
	go func() {
		err := miniurl.Run()
		assert.NoError(t, err)
		errCh <- err
	}()
	time.Sleep(1 * time.Second)

	testUrl := fmt.Sprintf(`{"url": "%s"}`, srv.URL+"/testpath")
	resp, err := http.Post("http://127.0.0.1:8080/api/v1/url", "application/json", strings.NewReader(testUrl))
	require.NoError(t, err)

	var v api.AddUrlResp
	err = json.NewDecoder(resp.Body).Decode(&v)
	require.NoError(t, err)

	resp, err = http.Get("http://127.0.0.1:8080/" + v.Hash)
	require.NoError(t, err)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "/testpath", string(data))

	syscall.Kill(os.Getpid(), syscall.SIGINT)
	require.NoError(t, <-errCh)
}

func setupTestEnv(t *testing.T) {
	stack, err := compose.NewDockerCompose("docker-compose.yaml")
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, stack.Down(context.Background(), compose.RemoveVolumes(true), compose.RemoveOrphans(true), compose.RemoveImagesLocal))
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	require.NoError(t, stack.WaitForService("db", wait.ForHealthCheck()).Up(ctx, compose.Wait(true)))
}
