package application

import (
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/domain/repository"
)

type userApps struct {
	userApp repository.UserRepository
}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (entity.User, map[string]string)
}

// implement contract user-repository
var _ UserAppInterface = &userApps{}

func (u *userApps) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.userApp.SaveUser(user)
}

func (u *userApps) GetUser(userId uint64) (*entity.User, error) {
	return u.userApp.GetUser(userId)
}

func (u *userApps) GetUsers() ([]entity.User, error) {
	return u.userApp.GetUsers()
}

func (u *userApps) GetUserByEmailAndPassword(user *entity.User) (entity.User, map[string]string) {
	return u.userApp.GetUserByEmailAndPassword(user)
}
