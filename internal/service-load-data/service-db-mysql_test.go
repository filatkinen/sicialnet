//go:build manual

package main_test

import (
	"context"
	"log"
	"testing"
)

func BenchmarkAppMysqlAdd(b *testing.B) {
	app, users, err := getapp("../../configs/server.mysql.yaml", "../../p.csv")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	defer app.Close(ctx)

	for i := 0; i < b.N; i++ {
		for i := range users {
			_, err = app.UserAdd(ctx, users[i], "pass")
		}
	}
}
