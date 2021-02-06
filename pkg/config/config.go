package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Observer interface {
	GetServerIP() string
	GetServerPort() uint
	GetDBMS() string
	GetDBConnection() string
	GetBaseURL() string
	GetSettingsURL() string
}

type observerImpl struct {
	serverIp     string
	serverPort   uint
	dbms         string
	dbConnection string
	baseURL      string
	settingsURL  string
}

func CreateConfigObserver() Observer {
	dbHost := viper.GetString("database.host")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbMode := viper.GetString("database.sslmode")

	return &observerImpl{
		serverIp:     viper.GetString("server.ip"),
		serverPort:   viper.GetUint("server.port"),
		dbms:         viper.GetString("database.dbms"),
		dbConnection: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode),
		baseURL:      viper.GetString("api.baseURL"),
		settingsURL:  viper.GetString("api.settingsURL"),
	}
}

func (observer *observerImpl) GetServerIP() string {
	return observer.serverIp
}

func (observer *observerImpl) GetServerPort() uint {
	return observer.serverPort
}

func (observer *observerImpl) GetDBMS() string {
	return observer.dbms
}

func (observer *observerImpl) GetDBConnection() string {
	return observer.dbConnection
}

func (observer *observerImpl) GetBaseURL() string {
	return observer.baseURL
}

func (observer *observerImpl) GetSettingsURL() string {
	return observer.settingsURL
}
