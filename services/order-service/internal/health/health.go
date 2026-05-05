package health

import (
	"crypto/sha256"
	"encoding/hex"
)

func Status() string {
	sum := sha256.Sum256([]byte("order-service"))
	return "ok:" + hex.EncodeToString(sum[:8])
}