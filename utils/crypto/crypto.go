package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

/*
// bcrypt hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
} */

type Argon2Config struct {
	saltLength  uint32
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

var argon2Config = &Argon2Config{
	saltLength:  16,
	memory:      64 * 1024,
	iterations:  4,
	parallelism: 1,
	keyLength:   32,
}

// argon2 hash
func HashPassword(password string) (string, error) {
	config := argon2Config

	salt, err := generateRandomBytes(config.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, config.iterations, config.memory, config.parallelism, config.keyLength)
	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, config.memory, config.iterations, config.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func VerifyPassword(password, encodedHash string) (bool, error) {
	config, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, config.iterations, config.memory, config.parallelism, config.keyLength)
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

func decodeHash(encodedHash string) (config *Argon2Config, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	config = &Argon2Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &config.memory, &config.iterations, &config.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.StdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.saltLength = uint32(len(salt))

	hash, err = base64.StdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.keyLength = uint32(len(hash))

	return config, salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
