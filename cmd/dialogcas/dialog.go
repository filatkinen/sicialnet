package main

import (
	"context"
	"flag"
	"github.com/filatkinen/socialnet/internal/config/dialogcas"
	internalhttp "github.com/filatkinen/socialnet/internal/dialogcas/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/dialogcas.yaml", "Path to configuration file")
	flag.Parse()

	conf, err := dialogcas.NewConfig(configFile)
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)

	server, err := internalhttp.NewServer(conf, logger)
	if err != nil {
		log.Fatalf("error creating server http %v", err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			logger.Printf("errors  on closing server %s\n", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	signalFailedStart := make(chan struct{})
	go func() {
		if err := server.Start(); err != nil {
			logger.Println("failed to start dialog server: " + err.Error())
		}
		signalFailedStart <- struct{}{}
	}()

	select {
	case <-signalFailedStart:
	case <-signalCh:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			log.Println("failed to stop http dialog server: " + err.Error())
		}
	}
}
