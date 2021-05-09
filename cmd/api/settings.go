package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Settings struct {
	LogFile  string
	LogLevel string

	BaseURL        string
	SettingsURL    string
	ProfileURL     string
	SessionsURL    string
	TeamsURL       string
	SportsURL      string
	TournamentsURL string
	MeetingsURL    string
	SkillsURL      string
	AttachURL      string
	MessageURL     string
	InviteURL      string
	SearchURL string

	ServerAddress string

	Origins []string

	OcrAddress string

	PsqlName string
	PsqlData string

	RedisAddress       string
	RedisProtocol      string
	RedisExpiresKeySec uint

	RabbitMQConnAddress  string
	RabbitMQEventQueueId string

	S3Bucket string
	S3Region string
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

		BaseURL:        viper.GetString("api.baseURL"),
		SettingsURL:    viper.GetString("api.settingsURL"),
		ProfileURL:     viper.GetString("api.profileURL"),
		SessionsURL:    viper.GetString("api.sessionsUrl"),
		TeamsURL:       viper.GetString("api.teamsURL"),
		SportsURL:      viper.GetString("api.sportsURL"),
		TournamentsURL: viper.GetString("api.tournamentsURL"),
		MeetingsURL:    viper.GetString("api.meetingsURL"),
		SkillsURL:      viper.GetString("api.skillsURL"),
		AttachURL:      viper.GetString("api.attachURL"),
		MessageURL:     viper.GetString("api.messageURL"),
		InviteURL:      viper.GetString("api.inviteURL"),
		SearchURL:      viper.GetString("api.searchURL"),

		ServerAddress: viper.GetString("server.address"),

		Origins: viper.GetStringSlice("cors.allowed_origins"),

		OcrAddress: viper.GetString("ocr.address"),

		PsqlName: viper.GetString("psql.dbms"),
		PsqlData: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode),

		RedisAddress:       viper.GetString("redis.address"),
		RedisProtocol:      viper.GetString("redis.protocol"),
		RedisExpiresKeySec: viper.GetUint("redis.expiresKeySec"),

		RabbitMQConnAddress:  fmt.Sprintf("amqp://%s:%s@%s/", rbmqUser, rbmqPass, rbmqAddress),
		RabbitMQEventQueueId: viper.GetString("rabbitmq.queueId"),

		S3Bucket: viper.GetString("s3.bucket"),
		S3Region: viper.GetString("s3.region"),
	}
}
