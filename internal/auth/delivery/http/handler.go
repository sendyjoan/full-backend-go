package http

import (
	"context"
	"net/http"

	"backend-service-internpro/internal/auth"
	"backend-service-internpro/internal/auth/service"
	"backend-service-internpro/internal/pkg/constants"
	apperrors "backend-service-internpro/internal/pkg/errors"
	"backend-service-internpro/internal/pkg/response"

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
		Tags:    []string{"Authentication"},
	}, func(ctx context.Context, in *struct {
		Body          auth.LoginRequest
		UserAgent     string `header:"User-Agent"`
		XForwardedFor string `header:"X-Forwarded-For"`
	}) (*struct {
		Body auth.LoginResponse
	}, error) {
		ua := in.UserAgent
		ip := in.XForwardedFor

		access, refresh, err := h.svc.Login(in.Body.UsernameOrEmail, in.Body.Password, ua, ip)
		if err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return &struct {
					Body auth.LoginResponse
				}{
					Body: *response.Error(appErr.Message),
				}, nil
			}
			return &struct {
				Body auth.LoginResponse
			}{
				Body: *response.Error(constants.LoginFailed),
			}, nil
		}

		loginData := auth.LoginData{
			AccessToken:  access,
			RefreshToken: refresh,
		}

		return &struct {
			Body auth.LoginResponse
		}{
			Body: *response.Success(constants.LoginSuccess, loginData),
		}, nil
	})

	// POST /refresh
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/refresh",
		Summary: "Exchange refresh token for new access token",
		Tags:    []string{"Authentication"},
	}, func(ctx context.Context, in *struct {
		Body          auth.RefreshRequest
		UserAgent     string `header:"User-Agent"`
		XForwardedFor string `header:"X-Forwarded-For"`
	}) (*struct {
		Body auth.RefreshResponse
	}, error) {
		ua := in.UserAgent
		ip := in.XForwardedFor

		access, err := h.svc.Refresh(in.Body.RefreshToken, ua, ip)
		if err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return &struct {
					Body auth.RefreshResponse
				}{
					Body: *response.Error(appErr.Message),
				}, nil
			}
			return &struct {
				Body auth.RefreshResponse
			}{
				Body: *response.Error(constants.RefreshFailed),
			}, nil
		}

		refreshData := auth.RefreshData{
			AccessToken: access,
		}

		return &struct {
			Body auth.RefreshResponse
		}{
			Body: *response.Success(constants.RefreshSuccess, refreshData),
		}, nil
	})

	// POST /logout
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/logout",
		Summary: "Revoke refresh token (logout)",
		Tags:    []string{"Authentication"},
	}, func(ctx context.Context, in *struct {
		Body auth.RefreshRequest
	}) (*struct {
		Body auth.BasicResponse
	}, error) {
		if err := h.svc.Logout(in.Body.RefreshToken); err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return &struct {
					Body auth.BasicResponse
				}{
					Body: *response.Error(appErr.Message),
				}, nil
			}
			return &struct {
				Body auth.BasicResponse
			}{
				Body: *response.Error(constants.LogoutFailed),
			}, nil
		}
		return &struct {
			Body auth.BasicResponse
		}{
			Body: *response.SuccessWithoutData(constants.LogoutSuccess),
		}, nil
	})

	// POST /forgot
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/forgot",
		Summary: "Send OTP for password reset",
		Tags:    []string{"Authentication"},
	}, func(ctx context.Context, in *struct {
		Body auth.ForgotRequest
	}) (*struct {
		Body auth.BasicResponse
	}, error) {
		err := h.svc.Forgot(in.Body.Email)
		// Always return success message for security (prevent email enumeration)
		if err != nil {
			// Log the actual error for debugging but don't expose it
			// You can add logging here later
		}
		return &struct {
			Body auth.BasicResponse
		}{
			Body: *response.SuccessWithoutData(constants.OTPSent),
		}, nil
	})

	// POST /verify-otp
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/verify-otp",
		Summary: "Validate OTP for password reset",
		Tags:    []string{"Authentication"},
	}, func(ctx context.Context, in *struct {
		Body auth.VerifyOTPRequest
	}) (*struct {
		Body auth.BasicResponse
	}, error) {
		if err := h.svc.VerifyOTP(in.Body.Email, in.Body.OTP); err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return &struct {
					Body auth.BasicResponse
				}{
					Body: *response.Error(appErr.Message),
				}, nil
			}
			return &struct {
				Body auth.BasicResponse
			}{
				Body: *response.Error(constants.OTPInvalid),
			}, nil
		}
		return &struct {
			Body auth.BasicResponse
		}{
			Body: *response.SuccessWithoutData(constants.OTPVerified),
		}, nil
	})

	// POST /reset-password
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/reset-password",
		Summary: "Reset password with valid OTP",
		Tags:    []string{"Authentication"},
	}, func(ctx context.Context, in *struct {
		Body auth.ResetPasswordRequest
	}) (*struct {
		Body auth.BasicResponse
	}, error) {
		if err := h.svc.ResetPassword(in.Body.Email, in.Body.OTP, in.Body.NewPassword); err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				return &struct {
					Body auth.BasicResponse
				}{
					Body: *response.Error(appErr.Message),
				}, nil
			}
			return &struct {
				Body auth.BasicResponse
			}{
				Body: *response.Error(constants.PasswordResetFailed),
			}, nil
		}
		return &struct {
			Body auth.BasicResponse
		}{
			Body: *response.SuccessWithoutData(constants.PasswordResetSuccess),
		}, nil
	})
}
