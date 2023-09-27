package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/filatkinen/socialnet/internal/config/server"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"log"
	"os"
	"sync"
	"time"
)

func main() { //nolint
	var configFile string
	var dataFile string
	flag.StringVar(&configFile, "config", "../../configs/server.yaml", "Path to configuration file")
	flag.StringVar(&dataFile, "data", "", "Path to file with data file")
	flag.Parse()

	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	config, err := server.NewConfig(configFile)
	if err != nil {
		log.Printf("error reading config file %v", err)
		return
	}

	ctx := context.Background()
	defer ctx.Done()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
	app, err := socialapp.New(logger, config)
	if err != nil {
		log.Printf("Error creating app %s", err)
		return
	}
	defer app.Close(ctx)

	chanWork := make(chan string, 5)
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for {
				str, ok := <-chanWork
				if !ok {
					return
				}
				user, err1 := app.UserGetRandom(ctx)
				if err != nil {
					log.Println(err)
				}
				friend, err2 := app.UserGetRandom(ctx)
				if err != nil {
					log.Println(err)
				}

				if err1 == nil && err2 == nil {
					_, err = app.UserAddPost(ctx, user.Id, str)
					if err != nil {
						log.Println(err)
					}
					app.UserAddFriend(ctx, user.Id, friend.Id)
				}
			}
		}()
	}

	t1 := time.Now()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		chanWork <- scanner.Text()
		count++
		if count%1_000 == 0 {
			fmt.Println(count, time.Since(t1))
		}
	}
	close(chanWork)
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}
	fmt.Println(count, time.Since(t1))

}
