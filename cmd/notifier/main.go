package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sport4all/pkg/logger"
	"syscall"
)

var args struct {
	configFilePath string
}

func main() {
	flag.StringVar(&args.configFilePath, "c", "", "path to configuration file")
	flag.Parse()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-signals
		logger.Info("received signal: ", sig)
		cancel()
	}()

	service := CreateService(args.configFilePath)
	service.Run(ctx)

	logger.Info("notifier finished")
}
