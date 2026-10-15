package specs

import (
	"encoding/json"
	"strconv"
	"time"
)

func (s Specs) WithString(key string, value *string) Specs {
	if value != nil {
		return s.With(key, *value)
	}

	return s
}

func (s Specs) WithBool(key string, value *bool) Specs {
	if value != nil {
		return s.With(key, strconv.FormatBool(*value))
	}

	return s
}

func (s Specs) WithInt(key string, value *int) Specs {
	if value != nil {
		return s.With(key, strconv.Itoa(*value))
	}

	return s
}

func (s Specs) WithInt64(key string, value *int64) Specs {
	if value != nil {
		return s.With(key, strconv.FormatInt(*value, 10))
	}

	return s
}

func (s Specs) WithFloat64(key string, value *float64) Specs {
	if value != nil {
		return s.With(key, strconv.FormatFloat(*value, 'f', -1, 64))
	}

	return s
}

func (s Specs) WithTime(key string, value *time.Time) Specs {
	if value != nil {
		return s.With(key, value.Format(time.RFC3339Nano))
	}

	return s
}

func (s Specs) WithStrings(key string, values []string) Specs {
	if values == nil {
		return s
	}

	data, err := json.Marshal(values)
	if err != nil {
		return s
	}

	return s.With(key, string(data))
}

func (s Specs) GetString(key string) (string, bool) {
	return s.Get(key)
}

func (s Specs) GetBool(key string) (bool, bool, error) {
	value, ok := s.Get(key)
	if !ok {
		return false, false, nil
	}

	parsed, err := strconv.ParseBool(value)
	return parsed, true, err
}

func (s Specs) GetInt(key string) (int, bool, error) {
	value, ok := s.Get(key)
	if !ok {
		return 0, false, nil
	}

	parsed, err := strconv.Atoi(value)
	return parsed, true, err
}

func (s Specs) GetInt64(key string) (int64, bool, error) {
	value, ok := s.Get(key)
	if !ok {
		return 0, false, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	return parsed, true, err
}

func (s Specs) GetFloat64(key string) (float64, bool, error) {
	value, ok := s.Get(key)
	if !ok {
		return 0, false, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	return parsed, true, err
}

func (s Specs) GetTime(key string) (time.Time, bool, error) {
	value, ok := s.Get(key)
	if !ok {
		return time.Time{}, false, nil
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	return parsed, true, err
}

func (s Specs) GetStrings(key string) ([]string, bool, error) {
	value, ok := s.Get(key)
	if !ok {
		return nil, false, nil
	}

	var parsed []string
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return nil, true, err
	}

	return parsed, true, nil
}
