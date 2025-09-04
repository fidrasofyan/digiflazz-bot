package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func GenerateSecretToken(length int) (string, error) {
	if length <= 0 || length > 1024 {
		return "", fmt.Errorf("length must be in [1, 1024]")
	}
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("read randomness: %w", err)
	}
	// URL-safe ([-_]) and no padding; shorter strings.
	return base64.RawURLEncoding.EncodeToString(b), nil
}
