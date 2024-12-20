package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// generate random 6 digits code
func GenerateVerificationCode() (string, error) {
	max := big.NewInt(1000000) // 上限值
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	// format to 6 digits
	return fmt.Sprintf("%06d", n.Int64()), nil 
}
