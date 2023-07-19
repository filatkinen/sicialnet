//go:build mysql

package storage_test

import (
	"context"
	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	mysqlstorage "github.com/filatkinen/socialnet/internal/storage/mysql"
	"log"
	"os"
	"testing"
)

func TestMysqlStorage(t *testing.T) {
	_ = os.Setenv("SOCIALNET_DB_USER", "socialnet")
	_ = os.Setenv("SOCIALNET_DB_PASS", "socialnet")
	conf, err := server.NewConfig("../../configs/server.mysql.yaml")
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}
	conf.DB.DBPort = "3306"
	mysqlStorage, err := mysqlstorage.New(conf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ctx := context.Background()
	defer ctx.Done()

	err = mysqlStorage.Connect(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer mysqlStorage.Close(ctx)

	var s storage.Storage = mysqlStorage
	runTestStorage(t, ctx, s)
}
