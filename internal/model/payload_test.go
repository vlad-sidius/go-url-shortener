package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenRequestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    ShortenRequest
		wantErr bool
	}{
		{
			name: "valid URL",
			json: `{"url": "https://example.com"}`,
			want: ShortenRequest{URL: "https://example.com"},
		},
		{
			name: "empty URL",
			json: `{"url": ""}`,
			want: ShortenRequest{URL: ""},
		},
		{
			name:    "missing URL field",
			json:    `{}`,
			want:    ShortenRequest{},
			wantErr: false, // zero value is valid
		},
		{
			name:    "invalid JSON",
			json:    `{"url":`,
			wantErr: true,
		},
		{
			name: "extra fields ignored",
			json: `{"url": "https://test.com", "extra": "field"}`,
			want: ShortenRequest{URL: "https://test.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req ShortenRequest
			err := json.Unmarshal([]byte(tt.json), &req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, req)
			}
		})
	}
}

func TestShortenResponseMarshalJSON(t *testing.T) {
	resp := NewShortenResponse("http://short.url/abc123")
	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	expected := `{"result":"http://short.url/abc123"}`
	assert.JSONEq(t, expected, string(data))
}

func TestErrorResponseMarshalJSON(t *testing.T) {
	errResp := NewErrorResponse("Invalid URL format")
	data, err := json.Marshal(errResp)
	assert.NoError(t, err)

	expected := `{"message":"Invalid URL format"}`
	assert.JSONEq(t, expected, string(data))
}
