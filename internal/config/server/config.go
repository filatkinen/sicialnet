package server

import (
	"github.com/filatkinen/socialnet/internal/rabbit"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

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

type RedisConfig struct {
	RedisPort    string
	RedisAddress string
}

type Config struct {
	StoreType         string
	ServerPort        string
	ServerAddress     string
	ServerHTTPLogfile string
	DB                DBConfig
	Redis             RedisConfig
	Rabbit            rabbit.Config
}

func NewConfig(in string) (Config, error) {
	const DefaultMaxIdleTime time.Duration = time.Minute * 15
	viper.SetConfigType("yaml")
	viper.SetConfigFile(in)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	viper.SetDefault("db.maxopenconns", 10)
	viper.SetDefault("db.maxidleconns", 10)
	viper.SetDefault("db.maxidletime", DefaultMaxIdleTime.String())

	viper.SetDefault("bindings.grpcport", "50051")

	duration, err := time.ParseDuration(viper.GetString("db.maxidletime"))
	if err != nil {
		log.Printf("Error parsing db.maxidletime value: %s, using defaul tvalue:%s", err.Error(), DefaultMaxIdleTime)
		duration = DefaultMaxIdleTime
	}

	config := Config{
		StoreType:         os.Getenv(viper.GetString("env.type")),
		ServerPort:        viper.GetString("bindings.port"),
		ServerAddress:     viper.GetString("bindings.address"),
		ServerHTTPLogfile: viper.GetString("httplog.logfile"),
		DB: DBConfig{
			DBUser:    os.Getenv(viper.GetString("env.dbuser")),
			DBPass:    os.Getenv(viper.GetString("env.dbpass")),
			DBAddress: os.Getenv(viper.GetString("env.address")),
			DBPort:    os.Getenv(viper.GetString("env.port")),
			DBName:    os.Getenv(viper.GetString("env.db")),

			MaxOpenConns: viper.GetInt("db.maxopenconns"),
			MaxIdleConns: viper.GetInt("db.maxidleconns"),
			MaxIdleTime:  duration,
		},
		Redis: RedisConfig{
			RedisPort:    os.Getenv(viper.GetString("env.redisport")),
			RedisAddress: os.Getenv(viper.GetString("env.redisaddress")),
		},
		Rabbit: rabbit.Config{
			Port:         viper.GetString("rabbit.port"),
			Address:      viper.GetString("rabbit.address"),
			ExchangeName: viper.GetString("rabbit.exchange"),
			User:         os.Getenv(viper.GetString("env.rabbituser")),
			Password:     os.Getenv(viper.GetString("env.rabbitpass")),
			//Queue:        viper.GetString("rabbit.queue"),
			//Tag:          viper.GetString("rabbit.tag"),
		},
	}

	return config, nil
}
