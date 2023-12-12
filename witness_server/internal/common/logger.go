package common

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(logFilePath string, loggingLevel string, logFormat string) {
	// Create log directory
	logDir := filepath.Dir(logFilePath)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating log directory: '%s'", err)
	}

	// Open log file
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error opening log file: '%s'", err)
	}

	// Setup logger
	setLoggingLevel(loggingLevel)
	var consoleWriter io.Writer
	if logFormat == "json" {
		consoleWriter = os.Stdout
	} else {
		consoleWriter = zerolog.NewConsoleWriter()
	}
	multi := io.MultiWriter(consoleWriter, file)
	log.Logger = log.Output(multi).With().Caller().Logger()
}

func setLoggingLevel(level string) {
	if level == "trace" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if level == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if level == "info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if level == "warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if level == "error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if level == "fatal" {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if level == "panic" {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}
}
