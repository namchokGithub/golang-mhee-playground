package rsa

import "encoding/hex"

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) ToHex(str string) string {
	return hex.EncodeToString([]byte(str))
}
