package model

import "github.com/google/uuid"

type ShortURLModel struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewShortURLModel(shortURL, originalURL string) *ShortURLModel {
	return &ShortURLModel{
		UUID:        uuid.New().String(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
