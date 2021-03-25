package mock

import "github.com/akwanmaroso/ddd-drugs/domain/entity"

// UserAppInterface is a mock user app interface
type UserAppInterface struct {
	SaveUserFn                  func(user *entity.User) (*entity.User, map[string]string)
	GetUsersFn                  func() ([]entity.User, error)
	GetUserFn                   func(uint64) (*entity.User, error)
	GetUserByEmailAndPasswordFn func(*entity.User) (*entity.User, map[string]string)
}

// SaveUser calls the SaveUserFn
func (u *UserAppInterface) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.SaveUserFn(user)
}

// GetUsers calls the GetUsersFn
func (u *UserAppInterface) GetUsers() ([]entity.User, error) {
	return u.GetUsersFn()
}

// GetUser calls the GetUserFn
func (u *UserAppInterface) GetUser(userId uint64) (*entity.User, error) {
	return u.GetUserFn(userId)
}

// GetUserByEmailAndPassword calls the GetUserByEmailAndPasswordFn
func (u *UserAppInterface) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.GetUserByEmailAndPasswordFn(user)
}
