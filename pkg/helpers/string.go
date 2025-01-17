package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

func GenerateUniqueID() (string, error) {
	b := make([]byte, 12)

	timestamp := uint32(time.Now().Unix())
	b[0] = byte(timestamp >> 24)
	b[1] = byte(timestamp >> 16)
	b[2] = byte(timestamp >> 8)
	b[3] = byte(timestamp)

	if _, err := rand.Read(b[4:]); err != nil {
		return "", err
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}
