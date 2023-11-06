package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/argon2"
)

const (
	saltSize    = 16
	keySize     = 32
	timeCost    = 1
	memory      = 64 * 1024
	parallelism = 2
)

type HashInterface interface {
	GenerateHash(password string) (string, error)
	ComparePassword(hash, password string) (bool, error)
}

type Hash struct {
}

func NewHash() HashInterface {
	return &Hash{}
}

func (h *Hash) GenerateHash(password string) (string, error) {

	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, timeCost, memory, parallelism, keySize)
	saltAndHash := append(salt, hash...)
	encodedSaltAndHash := base64.RawStdEncoding.EncodeToString(saltAndHash)

	return encodedSaltAndHash, nil
}

func (h *Hash) ComparePassword(hash, password string) (bool, error) {
	decodedSaltAndHash, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	salt := decodedSaltAndHash[:saltSize]
	existingHash := decodedSaltAndHash[saltSize:]

	computedHash := argon2.IDKey([]byte(password), salt, timeCost, memory, parallelism, keySize)

	if subtle.ConstantTimeCompare(existingHash, computedHash) == 1 {
		return true, nil
	}

	return false, errors.New("password mismatch")
}
