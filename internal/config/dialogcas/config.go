package dialogcas

import (
	"log"
	"os"
	"time"

	common "github.com/filatkinen/socialnet/internal/config"
	"github.com/spf13/viper"
)

type Config struct {
	StoreType          string
	ServerPort         string
	ServerAddress      string
	ServiceGRPCAddress string
	ServiceGRPCPort    string
	ServerHTTPLogfile  string
	DB                 common.DBConfig
	CAS                common.DBConfigCAS
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
		StoreType:          os.Getenv(viper.GetString("env.type")),
		ServerPort:         viper.GetString("bindings.port"),
		ServerAddress:      viper.GetString("bindings.address"),
		ServerHTTPLogfile:  viper.GetString("httplog.logfile"),
		ServiceGRPCAddress: viper.GetString("servicegrpc.address"),
		ServiceGRPCPort:    viper.GetString("servicegrpc.port"),
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
		CAS: common.DBConfigCAS{
			DBUser:       os.Getenv(viper.GetString("cas.dbuser")),
			DBPass:       os.Getenv(viper.GetString("cas.dbpass")),
			DBConnString: os.Getenv(viper.GetString("cas.address")),
			DBKeySpace:   os.Getenv(viper.GetString("cas.db")),
		},
	}
	return config, nil
}
