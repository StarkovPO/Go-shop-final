package service

import (
	"crypto/rand"
	"encoding/hex"
)

func generateUID() string {

	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	uid := hex.EncodeToString(randBytes)

	uid = uid[:8] + "-" + uid[8:12] + "-" + uid[12:16] + "-" + uid[16:20] + "-" + uid[20:]

	return uid
}
