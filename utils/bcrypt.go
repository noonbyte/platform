package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		3,       // iterations
		64*1024, // memory (64 MB)
		4,       // threads
		32,      // key length
	)

	return base64.RawStdEncoding.EncodeToString(
		append(salt, hash...),
	), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	data, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return false, err
	}

	if len(data) < 16+32 {
		return false, errors.New("invalid hash length")
	}

	salt := data[:16]
	expectedHash := data[16:]

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		3,       // iterations
		64*1024, // memory (64 MB)
		4,       // threads
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}
