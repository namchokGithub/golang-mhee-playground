package rsa

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) ToHex(str string, isWrite bool) (string, error) {
	// Convert the whole PEM text (including newlines) to hex
	hexStr := hex.EncodeToString([]byte(str))

	if isWrite {
		// Write hex string to .txt file // Write to ./etc
		err := os.WriteFile("./etc/prv_key_hex.txt", []byte(hexStr), 0644)
		if err != nil {
			return "", fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return hexStr, nil
}
