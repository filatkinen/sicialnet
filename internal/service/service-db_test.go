package main_test

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/filatkinen/socialnet/internal/config/server"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"github.com/filatkinen/socialnet/internal/storage"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

func getapp(configFile string, dataFile string) (*socialapp.App, []*storage.User, error) {
	_ = os.Setenv("SOCIALNET_DB_USER", "socialnet")
	_ = os.Setenv("SOCIALNET_DB_PASS", "socialnet")
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	config, err := server.NewConfig(configFile)
	if err != nil {
		log.Fatalf("error reading config file %v", err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
	app, err := socialapp.New(logger, config)
	if err != nil {
		return nil, nil, err
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var users []*storage.User
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dataUser := strings.Split(scanner.Text(), ",")
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
		users = append(users, &user)
	}
	if err := scanner.Err(); err != nil {
		app.Close(context.Background())
		return nil, nil, err
	}

	return app, users, nil
}
