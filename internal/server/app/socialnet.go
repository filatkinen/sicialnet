package socialapp

import (
	"context"
	"errors"
	"github.com/filatkinen/socialnet/internal/common"
	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/storage"
	mysqlstorage "github.com/filatkinen/socialnet/internal/storage/mysql"
	pgsqlstorage "github.com/filatkinen/socialnet/internal/storage/pgsql"
	"log"
	"time"
)

type App struct {
	appLog  *log.Logger
	storage storage.Storage
}

const TokenTTL = time.Hour * 24

var (
	ErrorUserNotFound    = errors.New("user not found")
	ErrorUserPassInvalid = errors.New("user pass is invalid")

	ErrorTokenNotFound = errors.New("token not found")
	ErrorTokenExpire   = errors.New("token expire")
)

func New(servLog *log.Logger, config server.Config) (*App, error) {
	stor, err := newStorage(config)
	if err != nil {
		return nil, err
	}
	log.Println("application socialnet started")
	log.Println("application socialnet is using db:" + config.StoreType)

	return &App{
		appLog:  servLog,
		storage: stor,
	}, nil
}

func newStorage(config server.Config) (storage.Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	switch config.StoreType {
	//case "memory":
	//	return memorystorage.New(), nil
	case "mysql":
		stor, err := mysqlstorage.New(config)
		if err != nil {
			return nil, err
		}
		err = stor.Connect(ctx)
		if err != nil {
			return nil, err
		}
		return stor, err
	case "pgsql":
		stor, err := pgsqlstorage.New(config)
		if err != nil {
			return nil, err
		}
		err = stor.Connect(ctx)
		if err != nil {
			return nil, err
		}
		return stor, err
	default:
		return nil, errors.New("bad type store type in config file")
	}
}

func (a *App) Close(ctx context.Context) error {
	a.appLog.Println("application socialnet stopped")

	err := a.storage.Close(ctx)
	if err != nil {
		a.appLog.Println("DB was closed with error:" + err.Error())
		return err
	}
	a.appLog.Println("application socialnet DB connection was closed")

	return nil
}

// SetToken generate token for user and save to the DB as hash
func (a *App) SetToken(ctx context.Context, userID string) (string, error) {
	tokenString, err := common.TokenGenerator()
	if err != nil {
		return "", err
	}
	token := storage.Token{
		Hash:   common.Hasher(tokenString),
		Expire: time.Now().Add(TokenTTL).Truncate(time.Minute).UTC(),
		UserID: userID,
	}
	err = a.storage.TokenAdd(ctx, &token)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CheckToken check if token is valid and return userID
func (a *App) CheckToken(ctx context.Context, tokenString string) (string, error) {
	token, err := a.storage.TokenGet(ctx, common.Hasher(tokenString))
	if err != nil {
		if errors.Is(err, storage.ErrRecordNotFound) {
			return "", ErrorTokenNotFound
		}
		return "", err
	}
	if token.Expire.After(time.Now().UTC()) {
		err = a.storage.TokenDelete(ctx, common.Hasher(tokenString))
		return "", errors.Join(err, ErrorTokenExpire)
	}

	return token.UserID, nil
}

// UserLogin check user credentials and gives token
func (a *App) UserLogin(ctx context.Context, userID string, pass string) (string, error) {
	u, err := a.storage.UserCredentialGet(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrRecordNotFound) {
			return "", ErrorUserNotFound
		}
		return "", err
	}
	if !common.CheckPasswordHash(pass, u.PassBcrypt) {
		return "", ErrorUserPassInvalid
	}

	tokenString, err := a.SetToken(ctx, userID)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// UserAdd add user and returns user_id
func (a *App) UserAdd(ctx context.Context, user *storage.User, pass string) (string, error) {
	id, err := storage.UUID()
	if err != nil {
		return "", err
	}
	user.Id = id

	cryptedPass, err := common.HashPassword(pass)
	if err != nil {
		return "", err
	}

	err = a.storage.UserAdd(ctx, user)
	if err != nil {
		return "", err
	}

	err = a.storage.UserCredentialSet(ctx, &storage.UserCredential{
		UserID:     id,
		PassBcrypt: cryptedPass,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *App) UserGet(ctx context.Context, userID string) (*storage.User, error) {
	u, err := a.storage.UserGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (a *App) GetAge(ctx context.Context, birthDay time.Time) int {
	y1, _, _ := birthDay.Date()
	y2, _, _ := time.Now().Date()
	return y2 - y1
}