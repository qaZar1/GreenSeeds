package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(lvl zerolog.Level) zerolog.Logger {
	zerolog.CallerSkipFrameCount = 2
	zerolog.LevelFieldName = "lvl"
	zerolog.TimestampFieldName = "dt"
	zerolog.CallerFieldName = "call"
	zerolog.MessageFieldName = "msg"

	return zerolog.New(os.Stdout).
		Level(lvl).
		With().Timestamp().Caller().
		Logger()
}
