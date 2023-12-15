package main

import (
	// "fmt"
	"platform/config"
	"platform/logging"
	"platform/services"
)

func writeMessage(logger logging.Logger, cfg config.Configuration) {
	section, ok := cfg.GetSection("main")
	if ok {
		message, ok := section.GetString("message")
		if ok {
			logger.Info(message)
		} else {
			logger.Panic("Cannot find configuration setting")
		}
	} else {
		logger.Panic("Config section not found")
	}
}

func main() {

	services.RegisterDefaultServices()

	// var cfg config.Configuration
	// services.GetService(&cfg)

	// var logger logging.Logger
	// services.GetService(&logger)

	services.Call(writeMessage)
}
