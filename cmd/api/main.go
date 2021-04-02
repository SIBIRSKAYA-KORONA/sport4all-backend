package main

import (
	"flag"

	_ "sport4all/docs"
)

var args struct {
	configFilePath string
}

func main() {
	flag.StringVar(&args.configFilePath, "c", "", "path to configuration file")
	flag.Parse()

	srv := CreateServer(args.configFilePath)
	srv.Run()
}
