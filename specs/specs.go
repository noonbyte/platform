package specs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
	"os"
	"sort"
	"strings"
)

type Specs map[string]string

func NewSpecs() Specs {
	return Specs{}
}

func (s Specs) With(key, value string) Specs {
	if s == nil {
		s = make(Specs)
	}

	s[key] = value
	return s
}

func (s Specs) Get(key string) (string, bool) {
	value, ok := s[key]
	return value, ok
}

func (s Specs) GetOr(name, fallback string) string {
	if value, ok := s[name]; ok {
		return value
	}

	s[name] = fallback

	return fallback
}

func (s Specs) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Specs) Delete(key string) {
	delete(s, key)
}

func (s Specs) Clone() Specs {
	result := make(Specs, len(s))

	maps.Copy(result, s)

	return result
}

func (s Specs) Merge(other Specs) Specs {
	maps.Copy(s, other)
	return s
}

func (s Specs) IsEmpty() bool {
	return len(s) == 0
}

func (s Specs) Len() int {
	return len(s)
}

func (s Specs) String() *string {
	if len(s) == 0 {
		return nil
	}

	if os.Getenv("APP_ENVIRONMENT") == "development" {
		value := formatSpecs(s)
		return &value
	}

	data, err := json.Marshal(s)
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
