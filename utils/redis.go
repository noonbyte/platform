package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Specs is a map with string keys and any type of value
type Specs map[string]string

// FormatSpecValue formats the value for human-readable output
func FormatSpecValue(key string, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

// GenerateKeyWithSpecs generates an MD5 hash string from the JSON-encoded representation of Specs
// If APP_ENV=development, returns a human-readable string instead of a hash
func GenerateKeyWithSpecs(specs Specs) string {
	// Check if APP_ENV is set to "development"
	env := os.Getenv("APP_ENVIRONMENT")
	if env == "development" {
		// Convert each spec to a human-readable format like [key=1,x=d3f34f3,f=true]
		var humanReadableSpecs []string
		for key, value := range specs {
			humanReadableSpecs = append(humanReadableSpecs, FormatSpecValue(key, value))
		}

		// Join the human-readable spec parts into a single string
		humanReadableString := "[" + strings.Join(humanReadableSpecs, ",") + "]"
		return humanReadableString
	}

	// If not in development mode, generate an MD5 hash
	specsJson, err := json.Marshal(specs)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	base64String := base64.StdEncoding.EncodeToString(specsJson)
	return base64String
}
