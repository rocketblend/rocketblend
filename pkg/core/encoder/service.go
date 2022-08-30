package encoder

import (
	"crypto/md5"
	"encoding/base32"
	"encoding/hex"
)

type (
	Service struct {
	}
)

func NewService() *Service {
	return &Service{}
}

func (s *Service) Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (s *Service) Encode(str string) string {
	return base32.StdEncoding.EncodeToString([]byte(str))
}

func (s *Service) Decode(str string) (string, error) {
	dec, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}
