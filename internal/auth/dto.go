package auth

import "backend-service-internpro/internal/pkg/response"

// Login
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" form:"username_or_email"`
	Password        string `json:"password" form:"password"`
}

type LoginData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type LoginResponse = response.ApiResponse

// Refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type RefreshData struct {
	AccessToken string `json:"access_token"`
}

type RefreshResponse = response.ApiResponse

// Forgot password
type ForgotRequest struct {
	Email string `json:"email" form:"email"`
}
type VerifyOTPRequest struct {
	Email string `json:"email" form:"email"`
	OTP   string `json:"otp" form:"otp"`
}
type ResetPasswordRequest struct {
	Email       string `json:"email" form:"email"`
	OTP         string `json:"otp" form:"otp"`
	NewPassword string `json:"new_password" form:"new_password"`
}

type BasicResponse = response.ApiResponse
