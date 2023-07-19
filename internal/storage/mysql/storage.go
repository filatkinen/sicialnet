package mysqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	_ "github.com/go-sql-driver/mysql" // import mysql
)

type Storage struct { // TODO
	db  *sql.DB
	dsn string
}

func New(config server.Config) (*Storage, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
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
	db, err := sql.Open("mysql", dsn)
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
			  VALUES (?,?,?,?,?,?,?)`

	_, err := s.db.ExecContext(ctx, query, user.Id, user.FirstName, user.SecondName, user.Sex, user.BirthDate, user.Biography, user.City)
	return err
}

func (s *Storage) UserDelete(ctx context.Context, userID string) error {
	query := `DELETE FROM users 
              WHERE user_id=?`
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
                first_name=?, 
                second_name=?, 
                sex=?, 
                birthdate=?, 
                biography=?, 
                city=?
              WHERE user_id=?`
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
			where user_id=?`
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

func (s *Storage) TokenAdd(ctx context.Context, token *storage.Token) error {
	query := `INSERT INTO token (hash,user_id,expires) 
			  VALUES (?,?,?)`
	_, err := s.db.ExecContext(ctx, query, token.Hash, token.UserID, token.Expire)
	return err
}
func (s *Storage) TokenDelete(ctx context.Context, hash string) error {
	query := `DELETE FROM token 
              WHERE hash=?`
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
	query := `SELECT user_id,expires from token WHERE hash=?`
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
              WHERE user_id=?`
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
	query := `SELECT COUNT(*) FROM user_credentials WHERE user_id=?`
	var count int
	err = tx.QueryRowContext(ctx, query, cred.UserID).Scan(&count)
	if err != nil {
		return err
	}
	switch count {
	case 0:
		query = `INSERT INTO user_credentials (password, user_id)
				VALUES (?, ?)`
	default:
		query = `update user_credentials set password=?
				where user_id=?`
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
              WHERE user_id=?`
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
	query := `SELECT password from user_credentials WHERE user_id=?`
	if err := s.db.QueryRowContext(ctx, query, userID).Scan(&r.PassBcrypt); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	return &r, nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}
