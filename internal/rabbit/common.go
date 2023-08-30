package rabbit

import (
	"fmt"
	"time"
)

const (
	DefaultChecKTimeSheduler time.Duration = time.Second * 5
)

type Config struct {
	Port          string
	Address       string
	User          string
	Password      string
	Queue         string
	CheckInterval time.Duration
	Tag           string
	ExchangeName  string
}

func (r Config) GetDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.User, r.Password, r.Address, r.Port)
}
