package model

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewShortenResponse(result string) *ShortenResponse {
	return &ShortenResponse{Result: result}
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{Message: message}
}
