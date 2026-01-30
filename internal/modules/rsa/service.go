package rsa

import "encoding/hex"

func ToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}
