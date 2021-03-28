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
}

func InitSettings(configFilePath string) Settings {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	rbmqAddress := viper.GetString("rabbitmq.address")
	rbmqUser := viper.GetString("rabbitmq.user")
	rbmqPass := viper.GetString("rabbitmq.password")

	return Settings{
		LogFile:  viper.GetString("logger.logfile"),
		LogLevel: viper.GetString("logger.level"),

		RabbitMQConnAddress:  fmt.Sprintf("amqp://%s:%s@%s/", rbmqUser, rbmqPass, rbmqAddress),
		RabbitMQEventQueueId: viper.GetString("rabbitmq.queueId"),
	}
}