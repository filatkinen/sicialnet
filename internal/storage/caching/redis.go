package caching

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq" // import pq
)

const maxPosts = 1000

type RedisCache struct {
	cache *redis.Client
	db    *sql.DB
	log   *log.Logger

	ready                atomic.Bool
	isUpdating           atomic.Bool
	exitSignalUpdateChan chan struct{}
}

func NewCache(config server.Config, log *log.Logger) (*RedisCache, error) {
	client, err := newRedis(config, log)
	if err != nil {
		return nil, err
	}
	db, err := newDB(config, log)
	if err != nil {
		client.Close()
		return nil, err
	}
	return &RedisCache{
		cache:                client,
		db:                   db,
		log:                  log,
		exitSignalUpdateChan: make(chan struct{}),
	}, nil

}

func newRedis(config server.Config, log *log.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.Redis.RedisAddress, config.Redis.RedisPort),
		Password: "",
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		log.Printf("%s", "Error opening redis cache")
		return nil, err
	}
	return client, nil
}

func newDB(config server.Config, log *log.Logger) (*sql.DB, error) {
	// connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&Timezone=UTC",
		config.DB.DBUser, config.DB.DBPass, config.DB.DBAddress, config.DB.DBPort, config.DB.DBName)
	// open database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("%s", "Error opening postgres database(redis cache)")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Printf("%s", "Error ping postgres database(redis cache)")
		db.Close()
		return nil, err
	}
	return db, nil
}

func (r *RedisCache) Ready() bool {
	return r.ready.Load()
}

func (r *RedisCache) Close() {
	r.StopUpdatePostAll()
	r.cache.Close()
	r.db.Close()
}

func (r *RedisCache) StopUpdatePostAll() {
	if r.isUpdating.Load() {
		go func() { r.exitSignalUpdateChan <- struct{}{} }()
		for {
			if !r.isUpdating.Load() {
				select {
				case <-r.exitSignalUpdateChan:
				default:
				}
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (r *RedisCache) UpdatePostAll() {
	if !r.isUpdating.Load() {
		r.isUpdating.Store(true)
		r.Clear()
		go func(r *RedisCache) {
			defer r.isUpdating.Store(false)
			r.ready.Store(false)
			t1 := time.Now()
			r.log.Printf("Starting process updating post cache %s\n", t1.UTC().Truncate(0))
			query := `select post_id,user_id,post_date,post from posts ORDER BY post_date DESC`
			ctx := context.Background()
			rows, err := r.db.QueryContext(ctx, query)
			if err != nil {
				r.log.Printf("Error  query postgres database(redis cache)%s\n", err)
				return
			}
			defer rows.Close()
			var post storage.Post
			for rows.Next() {
				select {
				case <-r.exitSignalUpdateChan:
					r.log.Printf("Got signal to exit from UpdatePosts%s\n", "")
					return
				default:
					err = rows.Scan(
						&post.PostId,
						&post.UserId,
						&post.PostDate,
						&post.PostText,
					)
					if err != nil {
						r.log.Printf("Error  query rows postgres database(redis cache)%s\n", err)
						return
					}
					err = r.AddPost(&post)
					if err != nil {
						r.Clear()
						r.log.Printf("Error  adding key redis cache%s\n", err)
						return
					}
				}
			}
			err = rows.Err()
			if err != nil {
				r.Clear()
				r.log.Printf("Error  row query postgres  (redis cache)%s\n", err)
				return
			}
			r.log.Printf("Finish process updating post cache %s\n", time.Since(t1).Truncate(0))
			r.ready.Store(true)
		}(r)
	}
}

func (r *RedisCache) Clear() {
	r.cache.FlushAll()
}

func (r *RedisCache) AddPost(postRecord *storage.Post) error {
	str := strings.Builder{}
	str.WriteString(postRecord.PostId)
	str.WriteString(postRecord.UserId)
	str.WriteString(postRecord.PostText)

	query := `select  friend_id from friends  where user_id=$1`
	rows, err := r.db.Query(query, postRecord.UserId)
	if err != nil {
		return err
	}
	defer rows.Close()
	var friend_id string
	for rows.Next() {
		err = rows.Scan(
			&friend_id,
		)
		if err != nil {
			return err
		}
		r.cache.RPush(friend_id, str.String())
		r.cache.LTrim(friend_id, 0, maxPosts-1)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) UserGetFriendsPosts(userID string, offset int, limit int) []*storage.Post {
	listPosts := r.cache.LRange(userID, int64(offset), int64(limit)-1).Val()
	var post storage.Post
	posts := make([]*storage.Post, 0, len(listPosts))
	for i := range listPosts {
		post.PostId = listPosts[i][0:36]
		post.UserId = listPosts[i][36:72]
		post.PostText = listPosts[i][72:]
		posts = append(posts, &post)
	}
	return posts
}
