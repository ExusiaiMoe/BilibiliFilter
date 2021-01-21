package utils

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	Logger zerolog.Logger
)

func StartLoggerModule()  {
	Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	Logger.Info().Msg("Logger Started Successfully")
}