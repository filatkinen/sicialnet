package pgsqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	_ "github.com/lib/pq" // import pq
	"math/rand"
	"time"
)

type Storage struct { // TODO
	db  *sql.DB
	dsn string
}

func New(config server.Config) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&Timezone=UTC",
		config.DB.DBUser, config.DB.DBPass, config.DB.DBAddress, config.DB.DBPort, config.DB.DBName)
	db, err := openDB(config.DB, dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db:  db,
		dsn: dsn,
	}, nil
}

func openDB(config server.DBConfig, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(config.MaxIdleTime)
	return db, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

func (s *Storage) UserAdd(ctx context.Context, user *storage.User) error {
	query := `INSERT INTO users  (user_id, first_name, second_name, sex, birthdate, biography, city) 
			  VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING user_id`

	_, err := s.db.ExecContext(ctx, query, user.Id, user.FirstName, user.SecondName, user.Sex, user.BirthDate, user.Biography, user.City)
	return err
}

func (s *Storage) UserDelete(ctx context.Context, userID string) error {
	query := `DELETE FROM users 
              WHERE user_id=$1`
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return storage.ErrRecordNotFound
	}
	return nil
}

func (s *Storage) UserUpdate(ctx context.Context, user *storage.User) error {
	query := `UPDATE  users set 
                first_name=$1, 
                second_name=$2, 
                sex=$3, 
                birthdate=$4, 
                biography=$5, 
                city=$6
              WHERE user_id=$7`
	args := []any{
		user.FirstName,
		user.SecondName,
		user.Sex,
		user.BirthDate,
		user.Biography,
		user.City,
		user.Id,
	}
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return storage.ErrRecordNotFound
	}
	return nil
}
func (s *Storage) UserGet(ctx context.Context, userID string) (*storage.User, error) {
	r := storage.User{Id: userID}
	query := `select  first_name, second_name, sex,biography, city, birthdate from users
			where user_id=$1`
	if err := s.db.QueryRowContext(ctx, query, userID).
		Scan(&r.FirstName, &r.SecondName, &r.Sex, &r.Biography, &r.City, &r.BirthDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	if r.BirthDate != nil {
		*r.BirthDate = r.BirthDate.UTC()
	}
	return &r, nil
}

func (s *Storage) UserGetRandom(ctx context.Context) (*storage.User, error) {
	query := `select count(user_id) from users`
	var count int
	if err := s.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	offset := r1.Intn(count)

	r := storage.User{}
	query = `select  user_id, first_name, second_name, sex,biography, city, birthdate from users
			LIMIT $1 OFFSET 1`
	if err := s.db.QueryRowContext(ctx, query, offset).
		Scan(&r.Id, &r.FirstName, &r.SecondName, &r.Sex, &r.Biography, &r.City, &r.BirthDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	if r.BirthDate != nil {
		*r.BirthDate = r.BirthDate.UTC()
	}
	return &r, nil
}

func (s *Storage) UserAddPost(ctx context.Context, post *storage.Post) (string, error) {
	query := `INSERT INTO posts  (post_id, user_id, post_date, post) 
			  VALUES ($1,$2,$3,$4) RETURNING post_id`
	var id string
	err := s.db.QueryRowContext(ctx, query, post.PostId, post.UserId, post.PostDate, post.PostText).Scan(&id)
	return id, err
}

func (s *Storage) UserAddFriend(ctx context.Context, userID string, friendID string) error {
	query := `INSERT INTO friends  (user_id, friend_id) 
			  VALUES ($1,$2)`

	_, err := s.db.ExecContext(ctx, query, userID, friendID)
	return err
}

func (s *Storage) UserGetFriends(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT friend_id
	FROM friends
	WHERE user_id =$1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var friend string
		err = rows.Scan(
			&friend,
		)
		result = append(result, friend)
	}
	if len(result) == 0 {
		return nil, storage.ErrRecordNotFound
	}
	return result, nil
}

func (s *Storage) UserGetFriendsPosts(ctx context.Context, userID string, offset int, limit int) ([]*storage.Post, error) {
	query := `select post_id,user_id,post_date, post from posts
    where posts.user_id = ANY(
		select  friends.friend_id from friends
    	where friends.user_id=$1
        )
		ORDER BY post_date
		offset $2 limit $3`

	rows, err := s.db.QueryContext(ctx, query, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*storage.Post
	for rows.Next() {
		var friend storage.Post
		err = rows.Scan(
			&friend.PostId,
			&friend.UserId,
			&friend.PostDate,
			&friend.PostText,
		)
		result = append(result, &friend)
	}
	if len(result) == 0 {
		return nil, storage.ErrRecordNotFound
	}
	return result, nil
}

func (s *Storage) TokenAdd(ctx context.Context, token *storage.Token) error {
	query := `INSERT INTO token (hash,user_id,expires) 
			  VALUES ($1,$2,$3)`
	_, err := s.db.ExecContext(ctx, query, token.Hash, token.UserID, token.Expire)
	return err
}
func (s *Storage) TokenDelete(ctx context.Context, hash string) error {
	query := `DELETE FROM token 
              WHERE hash=$1`
	result, err := s.db.ExecContext(ctx, query, hash)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return storage.ErrRecordNotFound
	}
	return nil
}

func (s *Storage) TokenGet(ctx context.Context, hash string) (*storage.Token, error) {

	r := storage.Token{Hash: hash}
	query := `SELECT user_id,expires from token WHERE hash=$1`
	if err := s.db.QueryRowContext(ctx, query, hash).Scan(&r.UserID, &r.Expire); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	r.Expire = r.Expire.UTC()
	return &r, nil
}

func (s *Storage) TokenDeleteAllUser(ctx context.Context, userID string) error {
	query := `DELETE FROM token 
              WHERE user_id=$1`
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return storage.ErrRecordNotFound
	}
	return nil
}

func (s *Storage) UserCredentialSet(ctx context.Context, cred *storage.UserCredential) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		}
	}()
	query := `SELECT COUNT(*) FROM user_credentials WHERE user_id=$1`
	var count int
	err = tx.QueryRowContext(ctx, query, cred.UserID).Scan(&count)
	if err != nil {
		return err
	}
	switch count {
	case 0:
		query = `INSERT INTO user_credentials (password, user_id)
				VALUES ($1, $2)`
	default:
		query = `update user_credentials set password=$1
				where user_id=$2`
	}
	_, err = tx.ExecContext(ctx, query, cred.PassBcrypt, cred.UserID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return err
}

func (s *Storage) UserCredentialDelete(ctx context.Context, userID string) error {
	query := `DELETE FROM user_credentials 
              WHERE user_id=$1`
	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return storage.ErrRecordNotFound
	}
	return nil
}
func (s *Storage) UserCredentialGet(ctx context.Context, userID string) (*storage.UserCredential, error) {

	r := storage.UserCredential{UserID: userID}
	query := `SELECT password from user_credentials WHERE user_id=$1`
	if err := s.db.QueryRowContext(ctx, query, userID).Scan(&r.PassBcrypt); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	return &r, nil
}

func (s *Storage) UserSearch(ctx context.Context, firstNameMask string, secondNameMask string) ([]*storage.User, error) {
	query := `SELECT user_id, first_name, second_name, biography, city, birthdate
	FROM users
	WHERE first_name LIKE $1
  	AND second_name LIKE $2`
	rows, err := s.db.QueryContext(ctx, query, firstNameMask+"%", secondNameMask+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*storage.User
	for rows.Next() {
		var user storage.User
		err = rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.SecondName,
			&user.Biography,
			&user.City,
			&user.BirthDate,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &user)
	}
	if len(result) == 0 {
		return nil, storage.ErrRecordNotFound
	}
	return result, nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}
