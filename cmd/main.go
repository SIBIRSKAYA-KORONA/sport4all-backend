package main

import (
	"flag"
	"log"

	"github.com/spf13/viper"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/server"
	_ "github.com/SIBIRSKAYA-KORONA/sport4all-backend/docs"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

var args struct {
	configPath string
}

func main() {
	flag.StringVar(&args.configPath, "c", "", "path to configuration file")
	flag.StringVar(&args.configPath, "config", "", "path to configuration file")
	flag.Parse()

	log.Println("set up provided config file - " + args.configPath)
	viper.SetConfigFile(args.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("viper initialized ✓")

	logger.InitLogger()
	logger.Info("logger initialized ✓")

	srv := server.CreateServer()
	srv.Run()
}
