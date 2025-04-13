package identityapp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go_short/domain/identity/entity"
	"go_short/domain/identity/repository"
	"go_short/domain/identity/service"

	"github.com/golang-jwt/jwt/v5"
)

// ErrUserNotFound 等自訂錯誤
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("username or email already exists")
var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrInternal = errors.New("internal server error")
var ErrTokenGeneration = errors.New("failed to generate token")

// App 是 Identity 領域的應用服務
type App struct {
	userRepo        repository.UserRepository
	identityService service.IdentityService
	jwtSecret       []byte
	jwtExpiration   time.Duration
}

// NewApp 創建 Identity 應用服務實例
func NewApp(userRepo repository.UserRepository, identityService service.IdentityService) *App {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("Warning: JWT_SECRET environment variable not set. Using default insecure key.")
		secret = "a_very_insecure_default_secret_key_change_me"
	}

	expStr := os.Getenv("JWT_EXPIRATION_HOURS")
	expHours, err := time.ParseDuration(expStr + "h")
	if err != nil || expHours <= 0 {
		log.Printf("Warning: Invalid or missing JWT_EXPIRATION_HOURS. Using default 24 hours. Error: %v", err)
		expHours = 24 * time.Hour
	}

	return &App{
		userRepo:        userRepo,
		identityService: identityService,
		jwtSecret:       []byte(secret),
		jwtExpiration:   expHours,
	}
}

// RegisterUser 處理使用者註冊的用例
func (a *App) RegisterUser(ctx context.Context, username, email, password string) (*entity.User, error) {
	existingUser, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("Error finding user by username %s: %v", username, err)
		return nil, ErrInternal
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}
	existingUser, err = a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("Error finding user by email %s: %v", email, err)
		return nil, ErrInternal
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		IsActive: true,
	}
	if err := user.SetPassword(password); err != nil {
		log.Printf("Error hashing password for user %s: %v", username, err)
		return nil, ErrInternal
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		log.Printf("Error creating user %s: %v", username, err)
		return nil, ErrInternal
	}

	return user, nil
}

// AuthenticateUser 處理使用者登入認證的用例，並返回 JWT
func (a *App) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("Error finding user by username %s during auth: %v", username, err)
		return "", ErrAuthenticationFailed
	}
	if user == nil {
		return "", ErrAuthenticationFailed
	}
	if !user.IsActive {
		log.Printf("User %s is inactive", username)
		return "", ErrAuthenticationFailed
	}

	if !user.CheckPassword(password) {
		return "", ErrAuthenticationFailed
	}

	now := time.Now()
	user.LastLogin = &now
	if err := a.userRepo.Update(ctx, user); err != nil {
		log.Printf("Error updating last login for user %s: %v", username, err)
	}

	tokenString, err := a.generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT for user %s: %v", username, err)
		return "", ErrTokenGeneration
	}

	return tokenString, nil
}

func (a *App) generateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"usn": user.Username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(a.jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (a *App) ActivateUser(ctx context.Context, userID uint) error {
	return a.identityService.ActivateUser(ctx, userID)
}

func (a *App) DeactivateUser(ctx context.Context, userID uint) error {
	return a.identityService.DeactivateUser(ctx, userID)
}

// --- 可以添加其他用例，如 GetUserProfile, ChangePassword 等 ---
