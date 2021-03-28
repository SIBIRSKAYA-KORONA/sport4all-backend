package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Settings struct {
	LogFile  string
	LogLevel string

	RabbitMQConnAddress  string
	RabbitMQEventQueueId string

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

	rbmqAddress := viper.GetString("rabbitmq.address")
	rbmqUser := viper.GetString("rabbitmq.user")
	rbmqPass := viper.GetString("rabbitmq.password")

	return Settings{
		LogFile:  viper.GetString("logger.logfile"),
		LogLevel: viper.GetString("logger.level"),

		RabbitMQConnAddress:  fmt.Sprintf("amqp://%s:%s@%s/", rbmqUser, rbmqPass, rbmqAddress),
		RabbitMQEventQueueId: viper.GetString("rabbitmq.queueId"),

		RedisAddress:       viper.GetString("redis.address"),
		RedisProtocol:      viper.GetString("redis.protocol"),
		RedisExpiresKeySec: viper.GetUint("redis.expiresKeySec"),

		PsqlName: viper.GetString("psql.dbms"),
		PsqlData: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode),
	}
}
