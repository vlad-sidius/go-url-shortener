package service

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

type RandomHashGenerator struct{}

func NewRandomHashGenerator() *RandomHashGenerator {
	return &RandomHashGenerator{}
}

func (g *RandomHashGenerator) Generate() (string, error) {
	bytes := make([]byte, 6) // 6 байт → ~8 символов в base64
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:8], nil
}
