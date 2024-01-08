package store

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID       int64   `db:"id"`
	Name     string  `db:"name"`
	Email    string  `db:"email"`
	Password *string `db:"password"`
}

type UserFilter struct {
	ID   FilterInt64
	Name FilterString
}

type Privilege struct {
	ID          int64
	Name        string
	Description string
}

func HashPassword(password string) *string {
	b := make([]byte, 16) // salt length
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	iterations := uint32(1)
	memory := uint32(64 * 1024)
	parallelism := uint8(4)
	hash := argon2.IDKey([]byte(password), b, iterations, memory, parallelism, 32)
	b64Salt := base64.StdEncoding.EncodeToString(b)
	b64Hash := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)
	return &encodedHash
}
