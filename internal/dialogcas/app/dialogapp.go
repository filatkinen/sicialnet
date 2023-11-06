package dialogappcas

import (
	"context"
	"errors"
	"github.com/filatkinen/socialnet/internal/config/dialogcas"
	"github.com/gocql/gocql"
	"log"
	"strings"
	"time"

	"github.com/filatkinen/socialnet/internal/storage"
	pgsqlstorage "github.com/filatkinen/socialnet/internal/storage/pgsql"
)

type DialogApp struct {
	appLog  *log.Logger
	Storage storage.Storage
	Session *gocql.Session
}

var (
	ErrorUserNotFound    = errors.New("user not found")
	ErrorUserPassInvalid = errors.New("user pass is invalid")
)

func New(servLog *log.Logger, config dialogcas.Config) (*DialogApp, error) {
	stor, err := newStorage(config)
	if err != nil {
		return nil, err
	}
	log.Println("application dialog started")
	log.Println("application dialog is using db:" + config.StoreType)

	conServers := strings.Split(config.CAS.DBConnString, ",")
	cluster := gocql.NewCluster(conServers...)
	cluster.Keyspace = config.CAS.DBKeySpace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.CAS.DBUser,
		Password: config.CAS.DBPass,
	}

	Session, err := cluster.CreateSession()
	if err != nil {
		stor.Close(context.Background())
		return nil, err
	}
	log.Println("application dialog is using cassandra")

	return &DialogApp{
		appLog:  servLog,
		Storage: stor,
		Session: Session,
	}, nil
}

func newStorage(config dialogcas.Config) (storage.Storage, error) {
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
	a.Session.Close()
	a.appLog.Println("application dialog Cassandra connection was closed")
	return nil
}

func (a *DialogApp) UserDialogSendMessage(_ context.Context, userID string, friendID string, message string) error {
	err := a.Session.Query("INSERT INTO dialogs (user_id, friend_id, dialog_id, message) VALUES (?, ?, uuid(), ?)",
		userID, friendID, message).Exec()
	return err
}

func (a *DialogApp) UserDialogListMessages(_ context.Context, userID string, friendID string) ([]*storage.DialogMessage, error) {
	var mess storage.DialogMessage
	messages := make([]*storage.DialogMessage, 0)

	m := map[string]interface{}{}
	iter := a.Session.Query("select friend_id, message from dialogs WHERE user_id=? and friend_id=?",
		userID, friendID).Iter()
	for iter.MapScan(m) {
		mess.From = m["user_id"].(string)
		mess.To = m["friend_id"].(string)
		mess.Text = m["message"].(string)
		messages = append(messages, &mess)
	}
	return messages, nil

}
