package identityapp

import (
	"context"
	"errors" // 引入 errors

	"go_short/domain/identity/entity"
	"go_short/domain/identity/repository"
	// 可能需要引入 domain/identity/service (如果有的話)
)

// ErrUserNotFound 等自訂錯誤
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("username or email already exists")
var ErrAuthenticationFailed = errors.New("authentication failed")

// App 是 Identity 領域的應用服務
type App struct {
	userRepo repository.UserRepository
	// 可能注入其他依賴，如密碼雜湊服務、權杖生成服務等
}

// NewApp 創建 Identity 應用服務實例
func NewApp(userRepo repository.UserRepository) *App {
	return &App{
		userRepo: userRepo,
	}
}

// RegisterUser 處理使用者註冊的用例
func (a *App) RegisterUser(ctx context.Context, username, email, password string) (*entity.User, error) {
	// 1. 檢查使用者名稱或郵箱是否已存在
	existingUser, _ := a.userRepo.FindByUsername(ctx, username)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}
	existingUser, _ = a.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 2. 創建 User 實體並雜湊密碼
	user := &entity.User{
		Username: username,
		Email:    email,
		IsActive: true, // 預設啟用
	}
	if err := user.SetPassword(password); err != nil {
		// 應記錄內部錯誤
		return nil, errors.New("failed to hash password")
	}

	// 3. 保存到儲存庫
	if err := a.userRepo.Create(ctx, user); err != nil {
		// 應記錄內部錯誤
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// AuthenticateUser 處理使用者登入認證的用例
func (a *App) AuthenticateUser(ctx context.Context, username, password string) (*entity.User, error) {
	// 1. 查找使用者
	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		// 應記錄內部錯誤
		return nil, ErrAuthenticationFailed
	}
	if user == nil {
		return nil, ErrAuthenticationFailed // 統一一返回認證失敗，避免洩漏用戶是否存在
	}

	// 2. 驗證密碼
	if !user.CheckPassword(password) {
		return nil, ErrAuthenticationFailed
	}

	// 3. (可選) 更新 LastLogin
	// now := time.Now()
	// user.LastLogin = &now
	// a.userRepo.Update(ctx, user) // 注意錯誤處理

	// 認證成功
	return user, nil
}

// --- 可以添加其他用例，如 GetUserProfile, ChangePassword 等 ---
