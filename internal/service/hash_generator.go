package service

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomSlugGenerator struct{}

func NewRandomSlugGenerator() *RandomSlugGenerator {
	return &RandomSlugGenerator{}
}

func (g *RandomSlugGenerator) Generate() string {
	bytes := make([]byte, 6) // 6 байт → ~8 символов в base64
	rand.Read(bytes)         // never returns an error, according docs
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}
