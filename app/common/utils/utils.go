package utils

import (
	"encoding/base64"
	"math/rand"

	"github.com/google/uuid"
)

// RandString generates random string with given length.
func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))] // nolint:gosec
	}
	return string(b)
}

// UUIDtoB64URL convert UUID to base64 string
func UUIDtoB64URL(id uuid.UUID) string {
	return base64.URLEncoding.EncodeToString(id[:])
}

// B64URLtoUUID converts base64 string to UUID
func B64URLtoUUID(s string) (uuid.UUID, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.ParseBytes(b)
}
