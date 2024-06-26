package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func CalculateFromFile(piece io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, piece); err != nil {
		return "", err
	}
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func CalculateFromHashes(hashes []string) (string, error) {
	hasher := sha256.New()
	for _, hash := range hashes {
		if _, err := hasher.Write([]byte(hash)); err != nil {
			return "", err
		}
	}
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}
