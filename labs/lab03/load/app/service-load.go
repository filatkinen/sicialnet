package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

const (
	port   = "5432"
	user   = "socialnet"
	dbname = "snet"
	maxID  = 1000000
)

func main() {
	var db *sql.DB

	password := os.Getenv("DBPASS")
	host := os.Getenv("DBHOST")
	// connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&Timezone=UTC",
		user, password, host, port, dbname)
	// open database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Printf("%s", "Error ping database")
		return
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	done := make(chan struct{})
	count := 0
	idLast := 0
	t1 := time.Now()
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			select {
			case <-signalCh:
				fmt.Println("Got signal to exit")
				fmt.Printf("Last insert ID=%d, RowsCount=%d\n", idLast, count)
				return
			default:
				idNew, err := insertRecord(db, TokenGenerator())
				if err != nil {
					fmt.Println(err)
					fmt.Printf("Last insert ID=%d, RowsCount=%d\n", idLast, count)
					return
				}
				idLast = idNew
				count++
				if count >= maxID {
					fmt.Printf("Without errors, get max rows. Last insert ID=%d, RowsCount=%d\n", idLast, count)
					return
				}
			}
		}
	}()

	<-done
	fmt.Println("Time to take:", time.Since(t1))
}

func TokenGenerator() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

func insertRecord(db *sql.DB, value string) (id int, err error) {
	query := `insert into test (uid)
              values ($1) returning id`
	err = db.QueryRow(query, value).Scan(&id)
	return id, err
}
