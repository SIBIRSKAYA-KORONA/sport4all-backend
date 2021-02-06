package hasher

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

const (
	lenSalt = 8
)

func hash(salt []byte, password string) []byte {
	hashPass := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return append(salt, hashPass...)
}

func HashPassword(password string) []byte {
	salt := make([]byte, lenSalt)
	if _, err := rand.Read(salt); err != nil {
		logger.Error(err)
	}
	return hash(salt, password)
}

func IsEqualPassword(checkPass string, realHashPass []byte) bool {
	var salt []byte
	salt = append(salt, realHashPass[0:lenSalt]...)
	return bytes.Equal(hash(salt, checkPass), realHashPass)
}
