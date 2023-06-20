package svc

import (
	"crypto/sha256"
	"fmt"
)

// Hash return sha256(s)
func Hash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
