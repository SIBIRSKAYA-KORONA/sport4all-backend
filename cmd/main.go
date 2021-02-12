package main

import (
	"flag"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/server"
	_ "github.com/SIBIRSKAYA-KORONA/sport4all-backend/docs"
)

var args struct {
	configFilePath string
}

func main() {
	flag.StringVar(&args.configFilePath, "c", "", "path to configuration file")
	flag.Parse()

	srv := server.CreateServer(args.configFilePath)
	srv.Run()
}
