//go:build mysql

package internalhttp_test

import (
	"context"
	"github.com/filatkinen/socialnet/internal/config/server"
	internalhttp "github.com/filatkinen/socialnet/internal/server/http"
	"log"
	"os"
	"testing"
	"time"
)

func TestHttpMysql(t *testing.T) {
	_ = os.Setenv("SOCIALNET_DB_USER", "socialnet")
	_ = os.Setenv("SOCIALNET_DB_PASS", "socialnet")
	_ = os.Setenv("SOCIALNET_DB", "snet")
	_ = os.Setenv("SOCIALNET_DB_ADDRESS", "localhost")
	_ = os.Setenv("SOCIALNET_PORT", "3306")
	_ = os.Setenv("SOCIALNET_DB_TYPE", "mysql")

	config, err := server.NewConfig("../../../configs/server.yaml")
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)

	srv, err := internalhttp.NewServer(config, l)
	if err != nil {
		log.Fatalf("error crating http server %s", err)
	}
	defer srv.Close()

	chanEnd := make(chan struct{})
	go func() {
		errstop := srv.Start()
		if errstop != nil {
			log.Fatalf("error starting server %s", errstop)
		}
		chanEnd <- struct{}{}
	}()
	// test status
	time.Sleep(time.Second)
	t.Run("test HTTP status", func(t *testing.T) {
		testHTTPStatus(t, config.ServerPort)
	})

	// test API
	t.Run("test HTTP API", func(t *testing.T) {
		testHTTPAPI(t, config.ServerPort)
	})

	err = srv.Stop(context.Background())
	if err != nil {
		log.Printf("error stopping server %s", err)
	}
	<-chanEnd
}
