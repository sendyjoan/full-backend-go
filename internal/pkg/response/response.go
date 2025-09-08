package response

// ApiResponse represents the standard API response format
type ApiResponse struct {
	Status  bool        `json:"status" doc:"Response status (true for success, false for error)"`
	Message string      `json:"message" doc:"Response message"`
	Data    interface{} `json:"data,omitempty" doc:"Response data"`
}

// Success creates a successful API response
func Success(message string, data interface{}) *ApiResponse {
	return &ApiResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

// Error creates an error API response
func Error(message string) *ApiResponse {
	return &ApiResponse{
		Status:  false,
		Message: message,
		Data:    nil,
	}
}

// SuccessWithoutData creates a successful API response without data
func SuccessWithoutData(message string) *ApiResponse {
	return &ApiResponse{
		Status:  true,
		Message: message,
		Data:    nil,
	}
}
