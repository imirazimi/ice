package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL MySQLConfig
	Redis RedisConfig
	HTTP  HTTPConfig
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type HTTPConfig struct {
	Port string
}

func Load() *Config {
	v := viper.New()
	v.SetConfigFile(".env") // read .env if present
	_ = v.ReadInConfig()    // ignore error: file may not exist
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// MySQL defaults
	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", "3306")
	v.SetDefault("mysql.user", "root")
	v.SetDefault("mysql.password", "root")
	v.SetDefault("mysql.database", "todos")
	// Redis defaults
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	// HTTP default
	v.SetDefault("http.port", "8080")

	conf := &Config{
		MySQL: MySQLConfig{
			Host:     v.GetString("mysql.host"),
			Port:     v.GetString("mysql.port"),
			User:     v.GetString("mysql.user"),
			Password: v.GetString("mysql.password"),
			Database: v.GetString("mysql.database"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("redis.addr"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		HTTP: HTTPConfig{
			Port: v.GetString("http.port"),
		},
	}
	return conf
}
