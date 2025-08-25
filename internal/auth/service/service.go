package service

import (
	"errors"
	"time"

	"backend-service-internpro/internal/auth"
	"backend-service-internpro/internal/auth/repository"
	jwtpkg "backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/otp"

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

type service struct {
	repo       repository.Repository
	secrets    jwtpkg.Secrets
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(repo repository.Repository, secrets jwtpkg.Secrets) Service {
	return &service{
		repo:       repo,
		secrets:    secrets,
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
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
	u, err := s.repo.FindUserByUsernameOrEmail(uore)
	if err != nil || !checkPassword(password, u.PasswordHash) {
		return "", "", errors.New("invalid credentials")
	}
	access, err := jwtpkg.GenerateAccess(u.ID.String(), s.secrets.Access, s.accessTTL)
	if err != nil {
		return "", "", err
	}
	refresh, err := jwtpkg.GenerateRefresh(u.ID.String(), s.secrets.Refresh, s.refreshTTL)
	if err != nil {
		return "", "", err
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
		return "", "", err
	}
	return access, refresh, nil
}

func (s *service) Refresh(refreshToken, ua, ip string) (string, error) {
	// find by hash
	rt, err := s.repo.GetRefreshToken(bcryptHash(refreshToken))
	if err != nil || rt.Revoked || time.Now().After(rt.ExpiresAt) {
		return "", errors.New("invalid refresh token")
	}
	// (opsional) cek UA/IP match → mitigasi token theft
	if ua != "" && rt.UserAgent != "" && ua != rt.UserAgent {
		return "", errors.New("ua mismatch")
	}
	if ip != "" && rt.IP != "" && ip != rt.IP {
		return "", errors.New("ip mismatch")
	}

	return jwtpkg.GenerateAccess(rt.UserID.String(), s.secrets.Access, s.accessTTL)
}

func (s *service) Logout(refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(bcryptHash(refreshToken))
	if err != nil {
		return err
	}
	return s.repo.RevokeRefreshToken(rt.ID)
}

func (s *service) Forgot(email string) error {
	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}
	code, err := otp.Generate6()
	if err != nil {
		return err
	}
	o := &auth.OTP{
		ID:        uuid.New(),
		UserID:    u.ID,
		Code:      code,
		Purpose:   "forgot_password",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	return s.repo.SaveOTP(o)
}

func (s *service) VerifyOTP(email, code string) error {
	_, err := s.repo.FindValidOTP(email, code, "forgot_password", time.Now())
	return err
}

func (s *service) ResetPassword(email, code, newPassword string) error {
	o, err := s.repo.FindValidOTP(email, code, "forgot_password", time.Now())
	if err != nil {
		return errors.New("invalid or expired otp")
	}

	hash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateUserPassword(o.UserID, hash); err != nil {
		return err
	}
	return s.repo.MarkOTPUsed(o.ID)
}

// helpers
func bcryptHash(s string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	return string(h)
}
