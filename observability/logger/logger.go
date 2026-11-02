package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level   string
	Console bool
	LogFile string
}

func Init(cfg Config) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var writers []io.Writer

	if cfg.Console {
		writers = append(writers, PrettyWriter{
			Out: zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: "15:04:05",
				NoColor:    false,
			},
		})
	} else {
		writers = append(writers, os.Stdout)
	}

	if cfg.LogFile != "" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   cfg.LogFile,
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		})
	}

	logger := zerolog.New(zerolog.MultiLevelWriter(writers...)).
		With().
		Timestamp().
		Logger()

	log.Logger = logger
	zerolog.DefaultContextLogger = &logger
}
