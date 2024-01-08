package main

import (
	"log"

	"github.com/softnetics/dotlocal/internal/daemon"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	err = daemon.Start(logger)
	if err != nil {
		logger.Fatal("Failed to start daemon", zap.Error(err))
	}
}
