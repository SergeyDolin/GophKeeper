package service

import (
	"errors"
	"main/internal/models"
	"main/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService - logic of auth
type AuthService struct {
	store     *storage.MemoryStorage
	jwtSecret []byte
}

func NewAuthService(store *storage.MemoryStorage) *AuthService {
	return &AuthService{
		store:     store,
		jwtSecret: []byte("secret"),
	}
}

// Register - registration for user
func (s *AuthService) Register(login, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	return s.store.SaveUser(models.User{
		Login:        login,
		PasswordHash: hash,
	})
}

// Login - login of user
func (s *AuthService) Login(login, password string) (string, error) {
	u, err := s.store.GetUser(login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}
