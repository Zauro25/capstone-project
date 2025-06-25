package models

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan pengguna sistem
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Password string `gorm:"not null" json:"-"` // Password tidak akan di-serialize ke JSON
	Role     string `gorm:"type:varchar(50);not null" json:"role" validate:"required,oneof=AdminPerpustakaan AdminDPK Executive"`
	LastLogin *time.Time `json:"last_login"`
}

// TableName untuk model User
func (User) TableName() string {
	return "users"
}
