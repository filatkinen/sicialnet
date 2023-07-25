package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/filatkinen/socialnet/internal/config/server"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"github.com/filatkinen/socialnet/internal/storage"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func newStringPointer(s string) *string {
	return &s
}

func newTimePointer(s time.Time) *time.Time {
	return &s
}

func toJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}

func main() {
	var configFile string
	var dataFile string
	flag.StringVar(&configFile, "config", "./configs/server.mysql.yaml", "Path to configuration file")
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
	//var storageDB storage.Storage
	//switch config.StoreType {
	//case "mysql":
	//	stor, err := mysqlstorage.New(config)
	//	if err != nil {
	//		log.Fatalf("Error creating storageDB app %s", err)
	//	}
	//	err = stor.Connect(ctx)
	//	if err != nil {
	//		log.Fatalf("Error creating storageDB app %s", err)
	//	}
	//	log.Println("Using mysql DB")
	//	storageDB = stor
	//case "pgsql":
	//	stor, err := pgsqlstorage.New(config)
	//	if err != nil {
	//		log.Fatalf("Error creating storageDB app %s", err)
	//	}
	//	err = stor.Connect(ctx)
	//	if err != nil {
	//		log.Fatalf("Error creating storageDB app %s", err)
	//	}
	//	log.Println("Using pgsql DB")
	//	storageDB = stor
	//default:
	//	log.Fatal("Bad type store type in config file")
	//}
	//defer storageDB.Close(ctx)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

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
				dataUser := strings.Split(str, ",")
				names := strings.Split(dataUser[0], " ")
				var user storage.User
				user.FirstName = names[1]
				user.SecondName = newStringPointer(names[0])
				user.City = newStringPointer(dataUser[2])
				year, err := strconv.Atoi(dataUser[1])
				if err != nil {
					year = r1.Intn(65) + 2
				}
				dateBirthDay := time.Date(2023-year, time.Month(r1.Intn(11)+1), r1.Intn(27)+1, 0, 0, 0, 0, time.UTC).Truncate(time.Hour * 24).UTC()
				user.BirthDate = newTimePointer(dateBirthDay)
				uid, _ := storage.UUID()
				user.Id = uid
				_, err = app.UserAdd(ctx, &user, "password")
				if err != nil {
					fmt.Println(err)
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
		if count%100_000 == 0 {
			fmt.Println(count, time.Since(t1))
		}
	}
	close(chanWork)
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
