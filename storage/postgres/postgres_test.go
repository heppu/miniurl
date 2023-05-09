package postgres_test

import (
	"context"
	"testing"

	"github.com/heppu/miniurl/storage/postgres"
	"github.com/heppu/miniurl/storage/storagetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresStorage(t *testing.T) {
	setupTestEnv(t)

	const connStr = "postgres://root:root@localhost:5432/root?sslmode=disable"
	s, err := postgres.NewStorage(connStr)
	require.NoError(t, err)

	storagetest.RunSuite(t, s)
}

func setupTestEnv(t *testing.T) {
	stack, err := compose.NewDockerCompose("../../docker-compose.yaml")
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, stack.Down(context.Background(), compose.RemoveVolumes(true), compose.RemoveOrphans(true), compose.RemoveImagesLocal))
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	require.NoError(t, stack.WaitForService("db", wait.ForHealthCheck()).Up(ctx, compose.Wait(true)))
}
