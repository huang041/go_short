package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt" // 引入 bcrypt 用於密碼雜湊
	"gorm.io/gorm"
)

// User 代表系統中的使用者
type User struct {
	gorm.Model              // 包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Username     string     `gorm:"type:varchar(100);uniqueIndex;not null"` // 使用者名稱，唯一且不為空
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null"` // 電子郵件，唯一且不為空
	PasswordHash string     `gorm:"type:varchar(255);not null"`             // 存儲雜湊後的密碼
	IsActive     bool       `gorm:"default:true"`                           // 帳號是否啟用
	LastLogin    *time.Time // 最後登入時間
}

// TableName 指定 User 實體的資料表名稱
func (User) TableName() string {
	return "users" // 資料表名建議用複數
}

// NewUser 創建一個新的 User 實例 (通常在應用層或服務層呼叫)
// 注意：密碼雜湊應該在這裡或服務層處理，這裡僅作結構示例
func NewUser(username, email, password string) (*User, error) {
	// 在實際應用中，密碼應在服務層進行雜湊處理
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	//     return nil, err
	// }
	return &User{
		Username: username,
		Email:    email,
		// PasswordHash: string(hashedPassword), // 實際應存儲雜湊值
		PasswordHash: "placeholder_hash", // 暫時使用佔位符
		IsActive:     true,
	}, nil
}

// SetPassword 雜湊並設定使用者密碼 (可以在領域服務中實現此邏輯)
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword 驗證提供的密碼是否與儲存的雜湊匹配 (可以在領域服務中實現)
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
