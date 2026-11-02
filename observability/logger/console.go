package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
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

	level := strings.ToUpper(fmt.Sprint(event["level"]))
	message := fmt.Sprint(event["message"])

	timestamp := ""
	if value, ok := event["time"]; ok {
		timestamp = fmt.Sprint(value)
	}

	if parsed, err := time.Parse(time.RFC3339Nano, timestamp); err == nil {
		timestamp = parsed.Format("15:04:05")
	}

	fmt.Fprintf(&buf, "%s %-5s %s\n", timestamp, level, message)

	keys := make([]string, 0, len(event))

	for key := range event {
		if key == "level" || key == "time" || key == "message" {
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Fprintf(&buf, "             %-12s = %v\n", key, event[key])
	}

	_, err := w.Out.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
