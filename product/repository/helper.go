package repository

import (
	"encoding/base64"
	"encoding/binary"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(token string) (string, error) {
	byt, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	tokenString := string(byt)

	return tokenString, nil
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(token int64) (string, error) {
	byteToken := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteToken, uint64(token))
	return base64.StdEncoding.EncodeToString(byteToken), nil
}
