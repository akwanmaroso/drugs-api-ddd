package entity

import (
	"html"
	"strings"
	"time"

	"github.com/akwanmaroso/ddd-drugs/infrastructure/security"
)

type User struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	FullName  string     `gorm:"size:100;not null" json:"fullname"`
	Email     string     `gorm:"size:100;not null;unique" json:"email"`
	Password  string     `gorm:"size:100;not null" json:"password"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PublicUser struct {
	ID       uint64 `gorm:"primary_key;auto_increment" json:"id"`
	FullName string `gorm:"size:100;not null" json:"fullname"`
}

func (user *User) BeforeSave() error {
	hashPassword, err := security.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return nil
}

func (user *User) PublicUser() interface{} {
	return &PublicUser{
		ID:       user.ID,
		FullName: user.FullName,
	}
}

func (user *User) Prepare() {
	user.FullName = html.EscapeString(strings.TrimSpace(user.FullName))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = user.PublicUser()
	}
	return result
}
