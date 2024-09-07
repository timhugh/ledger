package argon2id

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

const (
	time    = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func HashPassword(password, salt, pepper string) string {
	hash := argon2.IDKey([]byte(password), []byte(salt), time, memory, threads, keyLen)
	buffer := hmac.New(sha256.New, []byte(pepper))
	buffer.Write(hash)
	return base64.StdEncoding.EncodeToString(buffer.Sum(nil))
}
