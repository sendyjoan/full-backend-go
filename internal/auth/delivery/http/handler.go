package http

import (
	"context"
	"net/http"

	"backend-service-internpro/internal/auth"
	"backend-service-internpro/internal/auth/service"
	apperrors "backend-service-internpro/internal/pkg/errors"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct{ svc service.Service }

// New registers auth routes into the Huma API.
func New(api huma.API, svc service.Service) {
	h := &Handler{svc: svc}

	// Group /v1/auth
	g := huma.NewGroup(api, "/v1/auth")

	// POST /login
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/login",
		Summary: "Login and get access/refresh tokens",
	}, func(ctx context.Context, in *struct {
		Body          auth.LoginRequest `json:"body"`
		UserAgent     string            `header:"User-Agent"`
		XForwardedFor string            `header:"X-Forwarded-For"`
	}) (*auth.LoginResponse, error) {
		ua := in.UserAgent
		ip := in.XForwardedFor

		access, refresh, err := h.svc.Login(in.Body.UsernameOrEmail, in.Body.Password, ua, ip)
		if err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return nil, appErr.ToHumaError()
			}
			return nil, huma.Error500InternalServerError("login failed")
		}
		return &auth.LoginResponse{AccessToken: access, RefreshToken: refresh}, nil
	})

	// POST /refresh
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/refresh",
		Summary: "Exchange refresh token for new access token",
	}, func(ctx context.Context, in *struct {
		Body          auth.RefreshRequest `json:"body"`
		UserAgent     string              `header:"User-Agent"`
		XForwardedFor string              `header:"X-Forwarded-For"`
	}) (*auth.RefreshResponse, error) {
		ua := in.UserAgent
		ip := in.XForwardedFor

		access, err := h.svc.Refresh(in.Body.RefreshToken, ua, ip)
		if err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return nil, appErr.ToHumaError()
			}
			return nil, huma.Error500InternalServerError("refresh failed")
		}
		return &auth.RefreshResponse{AccessToken: access}, nil
	})

	// POST /logout
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/logout",
		Summary: "Revoke refresh token (logout)",
	}, func(ctx context.Context, in *struct {
		Body auth.RefreshRequest `json:"body"`
	}) (*auth.BasicResponse, error) {
		if err := h.svc.Logout(in.Body.RefreshToken); err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return nil, appErr.ToHumaError()
			}
			return nil, huma.Error500InternalServerError("logout failed")
		}
		return &auth.BasicResponse{Message: "logout successful"}, nil
	})

	// POST /forgot
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/forgot",
		Summary: "Send OTP for password reset",
	}, func(ctx context.Context, in *struct {
		Body auth.ForgotRequest `json:"body"`
	}) (*auth.BasicResponse, error) {
		err := h.svc.Forgot(in.Body.Email)
		// Always return success message for security (prevent email enumeration)
		if err != nil {
			// Log the actual error for debugging but don't expose it
			// You can add logging here later
		}
		return &auth.BasicResponse{Message: "If the email exists, an OTP has been sent"}, nil
	})

	// POST /verify-otp
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/verify-otp",
		Summary: "Validate OTP for password reset",
	}, func(ctx context.Context, in *struct {
		Body auth.VerifyOTPRequest `json:"body"`
	}) (*auth.BasicResponse, error) {
		if err := h.svc.VerifyOTP(in.Body.Email, in.Body.OTP); err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return nil, appErr.ToHumaError()
			}
			return nil, huma.Error500InternalServerError("verification failed")
		}
		return &auth.BasicResponse{Message: "OTP verified successfully"}, nil
	})

	// POST /reset-password
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/reset-password",
		Summary: "Reset password with valid OTP",
	}, func(ctx context.Context, in *struct {
		Body auth.ResetPasswordRequest `json:"body"`
	}) (*auth.BasicResponse, error) {
		if err := h.svc.ResetPassword(in.Body.Email, in.Body.OTP, in.Body.NewPassword); err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return nil, appErr.ToHumaError()
			}
			return nil, huma.Error500InternalServerError("password reset failed")
		}
		return &auth.BasicResponse{Message: "Password reset successful"}, nil
	})
}
