//go:build pgsql

package storage_test

import (
	"context"
	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	pgsqlstorage "github.com/filatkinen/socialnet/internal/storage/pgsql"
	"log"
	"os"
	"testing"
)

func TestPgsqlStorage(t *testing.T) {
	_ = os.Setenv("SOCIALNET_DB_USER", "socialnet")
	_ = os.Setenv("SOCIALNET_DB_PASS", "socialnet")
	conf, err := server.NewConfig("../../configs/server.pgsql.yaml")
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}
	conf.DB.DBPort = "5432"
	pgsqlStorage, err := pgsqlstorage.New(conf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ctx := context.Background()
	defer ctx.Done()

	err = pgsqlStorage.Connect(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer pgsqlStorage.Close(ctx)

	var s storage.Storage = pgsqlStorage
	runTestStorage(t, ctx, s)
}
