package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	dim   = "\033[2m"
	blue  = "\033[34m"
)

type PrettyWriter struct {
	Out io.Writer
}

func (w PrettyWriter) Write(p []byte) (int, error) {
	var event map[string]any

	if err := json.Unmarshal(p, &event); err != nil {
		return 0, err
	}

	var buf bytes.Buffer

	level := zerolog.InfoLevel
	if value, ok := event["level"]; ok {
		if parsed, err := zerolog.ParseLevel(fmt.Sprint(value)); err == nil {
			level = parsed
		}
	}

	levelName := strings.ToUpper(level.String())
	levelColor := zerolog.LevelColors[level]

	message := fmt.Sprint(event["message"])

	timestamp := ""
	if value, ok := event["time"]; ok {
		timestamp = fmt.Sprint(value)
	}

	if parsed, err := time.Parse(time.RFC3339Nano, timestamp); err == nil {
		timestamp = parsed.Format("15:04:05")
	}

	fmt.Fprintf(
		&buf,
		"%s%s%s %s%-5s%s %s%s%s\n",
		dim, timestamp, reset,
		fmt.Sprintf("\033[%dm", levelColor), levelName, reset,
		bold, message, reset,
	)

	keys := make([]string, 0, len(event))

	for key := range event {
		switch key {
		case "level", "time", "message":
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value, _ := json.Marshal(event[key])

		fmt.Fprintf(
			&buf,
			"             %s%-16s%s %s= %s%s\n",
			blue, key, reset,
			dim, string(value), reset,
		)
	}

	if w.Out == nil {
		return 0, fmt.Errorf("pretty writer: output is nil")
	}

	if _, err := w.Out.Write(buf.Bytes()); err != nil {
		return 0, err
	}

	return len(p), nil
}
