package utils

import (
	crypto "crypto/rand"
	"math/big"
)

func NewCryptoRand(size int64) int64 {
	safeNum, err := crypto.Int(crypto.Reader, big.NewInt(size))
	if err != nil {
		panic(err)
	}
	return safeNum.Int64()
}
