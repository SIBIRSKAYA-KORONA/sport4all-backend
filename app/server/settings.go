package server

import (
	"fmt"

	"github.com/spf13/viper"
)

type Settings struct {
	LogFile  string
	LogLevel string

	BaseURL     string
	SettingsURL string
	ProfileURL  string
	SessionsURL string

	ServerAddress string

	PsqlName string
	PsqlData string

	RedisAddress       string
	RedisProtocol      string
	RedisExpiresKeySec uint
}

func InitSettings(configFilePath string) Settings {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	dbHost := viper.GetString("psql.host")
	dbUser := viper.GetString("psql.user")
	dbPass := viper.GetString("psql.password")
	dbName := viper.GetString("psql.name")
	dbMode := viper.GetString("psql.sslmode")

	return Settings{
		LogFile:  viper.GetString("logger.logfile"),
		LogLevel: viper.GetString("logger.level"),

		BaseURL:     viper.GetString("api.baseURL"),
		SettingsURL: viper.GetString("api.settingsURL"),
		ProfileURL:  viper.GetString("api.profileURL"),
		SessionsURL: viper.GetString("api.sessionsUrl"),

		ServerAddress: viper.GetString("server.address"),

		PsqlName: viper.GetString("psql.dbms"),
		PsqlData: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode),

		RedisAddress:       viper.GetString("redis.address"),
		RedisProtocol:      viper.GetString("redis.protocol"),
		RedisExpiresKeySec: viper.GetUint("redis.expiresKeySec"),
	}
}
