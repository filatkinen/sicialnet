////go:build pgsql

package storage_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	pgsqlstorage "github.com/filatkinen/socialnet/internal/storage/pgsql"
	"github.com/stretchr/testify/require"
)

func TestPgsqlStorage(t *testing.T) {
	_ = os.Setenv("SOCIALNET_DB_USER", "socialnet")
	_ = os.Setenv("SOCIALNET_DB_PASS", "socialnet")
	_ = os.Setenv("SOCIALNET_DB", "snet")
	_ = os.Setenv("SOCIALNET_DB_ADDRESS", "localhost")
	_ = os.Setenv("SOCIALNET_PORT", "5432")
	_ = os.Setenv("SOCIALNET_DB_TYPE", "pgsql")
	conf, err := server.NewConfig("../../configs/server.yaml")
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}
	conf.DB.DBPort = "5432"
	pgsqlStorage, err := pgsqlstorage.New(conf.DB)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ctx := context.Background()
	defer ctx.Done()

	err = pgsqlStorage.Connect(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	defer pgsqlStorage.Close(ctx)

	var s storage.Storage = pgsqlStorage
	runTestStorage(t, ctx, s)
}
