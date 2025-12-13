package domain

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Details    any    `json:"details,omitempty"`
}

func NewAPIError(msg string, status int, details any) APIError {
	return APIError{
		msg,
		status,
		details,
	}
}
