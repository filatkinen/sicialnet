package server

import (
	"log"
	"os"
	"time"

	common "github.com/filatkinen/socialnet/internal/config"
	"github.com/filatkinen/socialnet/internal/rabbit"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	RedisPort    string
	RedisAddress string
}

type Config struct {
	StoreType         string
	ServerPort        string
	ServerGRPCPort    string
	ServerAddress     string
	ServerHTTPLogfile string
	DB                common.DBConfig
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
		ServerGRPCPort:    viper.GetString("bindings.grpcport"),
		ServerHTTPLogfile: viper.GetString("httplog.logfile"),
		DB: common.DBConfig{
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
		},
	}

	return config, nil
}
