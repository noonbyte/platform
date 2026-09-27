package logger

import (
	"io"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppName     string
	Version     string
	Environment string
	Level       string

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
		writers = append(
			writers,
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			},
		)
	} else {
		writers = append(
			writers,
			os.Stdout,
		)
	}

	if cfg.LogFile != "" {
		writers = append(
			writers,
			&lumberjack.Logger{
				Filename:   cfg.LogFile,
				MaxSize:    100, // MB
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   true,
			},
		)
	}

	multi := zerolog.MultiLevelWriter(writers...)

	logger := zerolog.New(multi).
		With().
		Timestamp().
		Logger()

	log.Logger = logger

	zerolog.DefaultContextLogger = &logger
}
