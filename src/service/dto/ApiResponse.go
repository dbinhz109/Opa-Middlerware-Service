package dto

type ApiResponse struct {
	// Application specific error code
	ErrorCode int `json:"error_code,omitempty"`
	Error     int `json:"error,omitempty"`
	// Message is the error message that may be displayed to end users
	Message string `json:"message,omitempty"`
	// Additional data if any
	Data interface{} `json:"data,omitempty"`
	// For responses with amount semantic
	TotalCount int64 `json:"total_count,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
