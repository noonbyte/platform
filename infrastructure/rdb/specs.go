package rdb

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type Specs map[string]string

func NewSpecs() Specs {
	return Specs{}
}

func (s Specs) With(name, value string) Specs {
	s[name] = value
	return s
}

func GenerateKeyWithSpecs(specs Specs) *string {
	if len(specs) == 0 {
		return nil
	}

	if os.Getenv("APP_ENVIRONMENT") == "development" {
		value := formatSpecs(specs)
		return &value
	}

	data, err := json.Marshal(specs)
	if err != nil {
		return nil
	}

	hash := sha256.Sum256(data)
	value := hex.EncodeToString(hash[:])

	return &value
}

func formatSpecs(specs Specs) string {
	keys := make([]string, 0, len(specs))
	for key := range specs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, key+"="+specs[key])
	}

	return "[" + strings.Join(values, ",") + "]"
}
