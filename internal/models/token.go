package models

import (
	"time"
)

const (
	TokenTypeRefresh       = "refresh"
	TokenTypeResetPassword = "resetPassword"
	TokenTypeVerifyEmail   = "verifyEmail"
)

type Token struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Token       string    `gorm:"index;not null" json:"token"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"userId"`
	Type        string    `gorm:"not null" json:"type"`
	Expires     time.Time `gorm:"not null" json:"expires"`
	Blacklisted bool      `gorm:"default:false" json:"blacklisted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}