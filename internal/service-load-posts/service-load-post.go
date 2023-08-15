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

//func newStringPointer(s string) *string {
//	return &s
//}
//
//func newTimePointer(s time.Time) *time.Time {
//	return &s
//}
//
//func toJSON(value any) string {
//	data, err := json.Marshal(value)
//	if err != nil {
//		return ""
//	}
//	return string(data)
//}

func main() {
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
		log.Fatalf("error reading config file %v", err)
	}

	ctx := context.Background()
	defer ctx.Done()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
	app, err := socialapp.New(logger, config)
	if err != nil {
		log.Fatalf("Error creating app %s", err)
	}
	defer app.Close(ctx)

	chanWork := make(chan string, 5)
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(chwork chan string) {
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
				//fmt.Println(toJSON(user))
			}
		}(chanWork)
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
		log.Fatal(err)
	}
	fmt.Println(count, time.Since(t1))

}
