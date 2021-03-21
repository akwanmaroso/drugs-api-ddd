package persistence

import (
	"errors"
	"strings"

	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/domain/repository"
	"github.com/akwanmaroso/ddd-drugs/infrastructure/security"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

// UserRepo implments the contract of repository.UserRepository interface
var _ repository.UserRepository = &UserRepo{}

func (repo *UserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	dbErr := map[string]string{}
	err := repo.db.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already exists"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (repo *UserRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := repo.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *UserRepo) GetUsers() ([]entity.User, error) {
	var users []entity.User
	err := repo.db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (repo *UserRepo) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
	var user entity.User
	dbErr := map[string]string{}
	err := repo.db.Where("email = ?", u.Email).Take(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	// Verify Password
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["incorrect_password"] = "incorrect password"
		return nil, dbErr
	}
	return &user, nil
}
