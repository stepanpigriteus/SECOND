package externalfunc

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}


