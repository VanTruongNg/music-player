package domain

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Username     string  `json:"username" gorm:"uniqueIndex;not null"`
	Email        string  `json:"email" gorm:"uniqueIndex;not null"`
	Password     string  `json:"-" gorm:"not null"`
	FullName     string  `json:"fullName"`
	Avatar       string  `json:"avatar"`
	TwoFAEnabled bool    `json:"twoFAEnabled" gorm:"not null;default:false"`
	TwoFASecret  string  `json:"twoFASecret" gorm:"size:128"`
	LastLoginAt  *string `json:"lastLoginAt" gorm:"type:timestamp"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	u.BaseModel = NewBaseModel()
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
