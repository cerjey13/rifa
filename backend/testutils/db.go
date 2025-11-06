package testutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"rifa/backend/pkg/config"
	"rifa/backend/pkg/db"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type pgState struct {
	Ctx       context.Context
	Cancel    context.CancelFunc
	Container tc.Container
	Database  db.DB
}

// ensureBackendDir change the working directory to be "backend" to avoid issues
// with the migrations folder when starting the test db.
func ensureBackendDir() error {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("could not determine current file path")
	}

	rootPath := filepath.Dir(currentFile)

	for !strings.HasSuffix(rootPath, "backend") && rootPath != "/" {
		rootPath = filepath.Dir(rootPath)
	}

	if rootPath == "/" {
		return errors.New("could not find backend directory in path")
	}

	return os.Chdir(rootPath)
}

// StartPostgres spins up a temporary Postgres container for integration tests.
func StartPostgres(t *testing.T) pgState {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)

	req := tc.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "rifa",
			"POSTGRES_PASSWORD": "rifa",
			"POSTGRES_DB":       "rifa",
		},
		Cmd: []string{"postgres"},
		WaitingFor: wait.ForSQL(
			"5432/tcp",
			"postgres",
			func(host string, port nat.Port) string {
				return fmt.Sprintf(
					"postgres://rifa:rifa@%s:%s/rifa?sslmode=disable",
					host,
					port.Port(),
				)
			},
		).WithStartupTimeout(60 * time.Second),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	err = ensureBackendDir()
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgres://rifa:rifa@%s:%s/rifa?sslmode=disable",
		host,
		port.Port(),
	)

	driver := db.NewPostgresDriver()

	var database db.DB
	const maxRetries = 6
	for i := 1; i <= maxRetries; i++ {
		if ctx.Err() != nil {
			break
		}
		database, err = db.Connect(
			ctx,
			driver,
			&config.DatabaseOpts{DatabaseUrl: dsn},
		)
		if err == nil {
			break
		}
		t.Logf("retrying DB connect (%d/%d): %v", i, maxRetries, err)
		time.Sleep(time.Duration(i) * time.Second)
	}

	if err != nil {
		logs, _ := container.Logs(ctx)
		buf := new(bytes.Buffer)
		_, _ = io.Copy(buf, logs)
		t.Logf("Postgres logs:\n%s", buf.String())

		t.Fatalf("failed to connect after retries: %v", err)
	}

	return pgState{
		Ctx:       ctx,
		Cancel:    cancel,
		Container: container,
		Database:  database,
	}
}

// StopPostgres cleans up the container and database connection.
func StopPostgres(t *testing.T, st pgState) {
	t.Helper()
	if st.Database != nil {
		st.Database.Close()
	}
	if st.Container != nil {
		if err := st.Container.Terminate(st.Ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}
	if st.Cancel != nil {
		st.Cancel()
	}
}
