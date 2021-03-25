package entity

import (
	"html"
	"strings"
	"time"

	"github.com/akwanmaroso/ddd-drugs/infrastructure/security"
	"github.com/badoux/checkmail"
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

// BeforeSave is responsible for memorizing the password before saving it to db
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

// Prepare is in charge of ensuring the data entered into the db is correct
func (user *User) Prepare() {
	user.FullName = html.EscapeString(strings.TrimSpace(user.FullName))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

// Validate is responsible for validating the data you want to save to db
func (user *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	case "login":
		if user.Password == "" {
			errorMessages["password_required"] = "password required"
		}
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	case "forgetpassword":
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	default:
		if user.FullName == "" {
			errorMessages["fullname_required"] = "fullname is required"
		}
		if user.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if user.Password == "" && len(user.Password) < 6 {
			errorMessages["invalid_password"] = "password should be at least 6 characters"
		}
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	}
	return errorMessages
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = user.PublicUser()
	}
	return result
}
