package stores

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type User struct {
	Password   *string `db:"password"`
	Privileges *[]Privilege
	Name       *string `db:"name"`
	Email      *string `db:"email"`
	ID         int64   `db:"id"`
}

type UserFilter struct {
	Name  FilterString
	Email FilterString
	ID    FilterInt64
}

type UserList struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
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

func VerifyPassword(password string, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		panic("invalid encodedHash")
	}
	if parts[1] != "argon2id" {
		panic(fmt.Errorf("invalid encodedHash algorithm '%s'", parts[1]))
	}
	if parts[2] != fmt.Sprintf("v=%d", argon2.Version) {
		panic(fmt.Errorf("invalid encodedHash version '%s' != 'v=%v'", parts[2], argon2.Version))
	}
	iterations := uint32(1)
	memory := uint32(64 * 1024)
	parallelism := uint8(4)
	for _, pp := range strings.Split(parts[3], ",") {
		ps := strings.Split(pp, "=")
		if len(ps) == 2 {
			v, err := strconv.ParseInt(ps[1], 10, 64)
			if err != nil {
				panic(fmt.Errorf("invalid encodedHash '%s'", pp))
			}
			switch ps[0] {
			case "m":
				memory = uint32(v)
			case "t":
				iterations = uint32(v)
			case "p":
				parallelism = uint8(v)
			}
		}
	}
	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	hash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))
	return parts[5] == base64.StdEncoding.EncodeToString(comparisonHash)
}
