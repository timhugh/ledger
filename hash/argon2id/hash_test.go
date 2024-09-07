package argon2id_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/timhugh/ledger/hash/argon2id"
	"testing"
)

const pepper = "pepper"

func Test_HashPassword(t *testing.T) {
	t.Run("creates same hash for same inputs", func(t *testing.T) {
		const password = "password"
		const salt = "salt"

		hash1 := argon2id.HashPassword(password, salt, pepper)
		hash2 := argon2id.HashPassword(password, salt, pepper)

		assert.Equal(t, hash1, hash2)
	})
	t.Run("creates different hash for different inputs", func(t *testing.T) {
		hash1 := argon2id.HashPassword("password", "salt", pepper)
		hash2 := argon2id.HashPassword("different password", "different salt", pepper)

		assert.NotEqual(t, hash1, hash2)
	})
}
