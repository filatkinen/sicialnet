package commondbconf

import "time"

type DBConfig struct {
	DBUser       string
	DBPass       string
	DBAddress    string
	DBPort       string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}
