package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "secret123"

	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash) // Harusnya hash tidak kosong
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	wrongPassword := "wrong123"

	hash, _ := HashPassword(password)

	// check password yang benar
	match := CheckPasswordHash(password, hash)
	assert.True(t, match, "password yang benar harusnya sesuai")

	// check password yang salah
	wrongMatch := CheckPasswordHash(wrongPassword, hash)
	assert.False(t, wrongMatch, "password yang salah harusnya tidak sesuai")
}
