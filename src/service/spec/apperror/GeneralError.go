package apperror

// AppError is the default error struct containing detailed information about the error
//
// swagger:model
type AppError struct {
	// Application specific error code
	ErrorCode int `json:"errorCode,omitempty"`
	// Message is the error message that may be displayed to end users
	Message *string `json:"message,omitempty"`
	// Additional data if any
	Data interface{} `json:"data,omitempty"`
}

// New generates an application error
func New(code int, msg string, data interface{}) *AppError {
	return &AppError{ErrorCode: code, Message: &msg, Data: data}
}

func NewClone(src *AppError, data interface{}) *AppError {
	return &AppError{ErrorCode: src.ErrorCode, Message: src.Message, Data: data}
}

// Error returns the error message.
func (e AppError) Error() string {
	if e.Message == nil {
		return ""
	}
	return *e.Message
}
