package auth

// Login
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// Forgot password
type ForgotRequest struct {
	Email string `json:"email"`
}
type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	NewPassword string `json:"new_password"`
}

type BasicResponse struct {
	Message string `json:"message"`
}
