package dialogapp

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/filatkinen/socialnet/internal/config/dialog"
	"github.com/filatkinen/socialnet/internal/storage"
	pgsqlstorage "github.com/filatkinen/socialnet/internal/storage/pgsql"
)

type DialogApp struct {
	appLog  *log.Logger
	Storage storage.Storage
}

var (
	ErrorUserNotFound    = errors.New("user not found")
	ErrorUserPassInvalid = errors.New("user pass is invalid")
)

func New(servLog *log.Logger, config dialog.Config) (*DialogApp, error) {
	stor, err := newStorage(config)
	if err != nil {
		return nil, err
	}
	log.Println("application dialog started")
	log.Println("application dialog is using db:" + config.StoreType)

	return &DialogApp{
		appLog:  servLog,
		Storage: stor,
	}, nil
}

func newStorage(config dialog.Config) (storage.Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	switch config.StoreType {
	case "pgsql":
		stor, err := pgsqlstorage.New(config.DB)
		if err != nil {
			return nil, err
		}
		err = stor.Connect(ctx)
		if err != nil {
			return nil, err
		}
		return stor, err
	default:
		return nil, errors.New("bad type store type in env")
	}
}

func (a *DialogApp) Close(ctx context.Context) error {
	a.appLog.Println("application dialog stopped")

	err := a.Storage.Close(ctx)
	if err != nil {
		a.appLog.Println("DB was closed with error:" + err.Error())
		return err
	}
	a.appLog.Println("application dialog DB connection was closed")

	return nil
}

func (a *DialogApp) UserDialogSendMessage(ctx context.Context, userID string, friendID string, message string) error {
	return a.Storage.UserDialogSendMessage(ctx, userID, friendID, message)
}

func (a *DialogApp) UserDialogListdMessages(ctx context.Context, userID string, friendID string) ([]*storage.DialogMessage, error) {
	return a.Storage.UserDialogListMessages(ctx, userID, friendID)
}
