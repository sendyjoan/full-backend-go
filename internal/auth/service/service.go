package service

import (
	"time"

	"backend-service-internpro/internal/auth"
	"backend-service-internpro/internal/auth/repository"
	apperrors "backend-service-internpro/internal/pkg/errors"
	jwtpkg "backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/otp"
	"backend-service-internpro/internal/pkg/validator"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(uore, password, ua, ip string) (access, refresh string, err error)
	Refresh(refreshToken, ua, ip string) (access string, err error)
	Logout(refreshToken string) error
	Forgot(email string) error
	VerifyOTP(email, code string) error
	ResetPassword(email, code, newPassword string) error
}

type Config struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type service struct {
	repo       repository.Repository
	secrets    jwtpkg.Secrets
	accessTTL  time.Duration
	refreshTTL time.Duration
	validator  *validator.Validator
}

func New(repo repository.Repository, secrets jwtpkg.Secrets) Service {
	return &service{
		repo:       repo,
		secrets:    secrets,
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
		validator:  validator.New(),
	}
}

func NewWithConfig(repo repository.Repository, secrets jwtpkg.Secrets, cfg Config) Service {
	return &service{
		repo:       repo,
		secrets:    secrets,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
		validator:  validator.New(),
	}
}

func hashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), err
}
func checkPassword(p, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func (s *service) Login(uore, password, ua, ip string) (string, string, error) {
	// Validate input
	if ok, msg := s.validator.IsRequired(uore, "username/email"); !ok {
		return "", "", apperrors.ValidationFailed(msg)
	}
	if ok, msg := s.validator.IsRequired(password, "password"); !ok {
		return "", "", apperrors.ValidationFailed(msg)
	}

	u, err := s.repo.FindUserByUsernameOrEmail(uore)
	if err != nil || !checkPassword(password, u.PasswordHash) {
		return "", "", apperrors.InvalidCredentials()
	}

	access, err := jwtpkg.GenerateAccess(u.ID.String(), s.secrets.Access, s.accessTTL)
	if err != nil {
		return "", "", apperrors.InternalServer("failed to generate access token")
	}

	refresh, err := jwtpkg.GenerateRefresh(u.ID.String(), s.secrets.Refresh, s.refreshTTL)
	if err != nil {
		return "", "", apperrors.InternalServer("failed to generate refresh token")
	}

	// store refresh token hash
	rt := &auth.RefreshToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		TokenHash: bcryptHash(refresh),
		UserAgent: ua,
		IP:        ip,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.repo.CreateRefreshToken(rt); err != nil {
		return "", "", apperrors.InternalServer("failed to store refresh token")
	}
	return access, refresh, nil
}

func (s *service) Refresh(refreshToken, ua, ip string) (string, error) {
	if ok, msg := s.validator.IsRequired(refreshToken, "refresh token"); !ok {
		return "", apperrors.ValidationFailed(msg)
	}

	// find by hash
	rt, err := s.repo.GetRefreshToken(bcryptHash(refreshToken))
	if err != nil || rt.Revoked || time.Now().After(rt.ExpiresAt) {
		return "", apperrors.InvalidRefreshToken()
	}

	// (opsional) cek UA/IP match → mitigasi token theft
	if ua != "" && rt.UserAgent != "" && ua != rt.UserAgent {
		return "", apperrors.Unauthorized().WithDetails("user agent mismatch")
	}
	if ip != "" && rt.IP != "" && ip != rt.IP {
		return "", apperrors.Unauthorized().WithDetails("ip address mismatch")
	}

	access, err := jwtpkg.GenerateAccess(rt.UserID.String(), s.secrets.Access, s.accessTTL)
	if err != nil {
		return "", apperrors.InternalServer("failed to generate access token")
	}

	return access, nil
}

func (s *service) Logout(refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(bcryptHash(refreshToken))
	if err != nil {
		return err
	}
	return s.repo.RevokeRefreshToken(rt.ID)
}

func (s *service) Forgot(email string) error {
	if ok, msg := s.validator.IsRequired(email, "email"); !ok {
		return apperrors.ValidationFailed(msg)
	}
	if !s.validator.IsValidEmail(email) {
		return apperrors.ValidationFailed("invalid email format")
	}

	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return apperrors.EmailNotFound()
	}

	code, err := otp.Generate6()
	if err != nil {
		return apperrors.InternalServer("failed to generate OTP")
	}

	o := &auth.OTP{
		ID:        uuid.New(),
		UserID:    u.ID,
		Code:      code,
		Purpose:   "forgot_password",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.repo.SaveOTP(o); err != nil {
		return apperrors.InternalServer("failed to save OTP")
	}

	return nil
}

func (s *service) VerifyOTP(email, code string) error {
	if ok, msg := s.validator.IsRequired(email, "email"); !ok {
		return apperrors.ValidationFailed(msg)
	}
	if ok, msg := s.validator.IsValidOTP(code); !ok {
		return apperrors.ValidationFailed(msg)
	}

	_, err := s.repo.FindValidOTP(email, code, "forgot_password", time.Now())
	if err != nil {
		return apperrors.InvalidOTP()
	}
	return nil
}

func (s *service) ResetPassword(email, code, newPassword string) error {
	if ok, msg := s.validator.IsRequired(email, "email"); !ok {
		return apperrors.ValidationFailed(msg)
	}
	if ok, msg := s.validator.IsValidOTP(code); !ok {
		return apperrors.ValidationFailed(msg)
	}
	if ok, msg := s.validator.IsValidPassword(newPassword); !ok {
		return apperrors.ValidationFailed(msg)
	}

	o, err := s.repo.FindValidOTP(email, code, "forgot_password", time.Now())
	if err != nil {
		return apperrors.InvalidOTP()
	}

	hash, err := hashPassword(newPassword)
	if err != nil {
		return apperrors.InternalServer("failed to hash password")
	}

	if err := s.repo.UpdateUserPassword(o.UserID, hash); err != nil {
		return apperrors.InternalServer("failed to update password")
	}

	if err := s.repo.MarkOTPUsed(o.ID); err != nil {
		return apperrors.InternalServer("failed to mark OTP as used")
	}

	return nil
}

// helpers
func bcryptHash(s string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	return string(h)
}
