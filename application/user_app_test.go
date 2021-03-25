package application

import (
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	saveUserRepo                func(user *entity.User) (*entity.User, map[string]string)
	getUserRepo                 func(userId uint64) (*entity.User, error)
	getUsersRepo                func() ([]entity.User, error)
	getUserEmailAndPasswordRepo func(user *entity.User) (*entity.User, map[string]string)
)

type fakeUserRepo struct{}

func (u *fakeUserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return saveUserRepo(user)
}

func (u *fakeUserRepo) GetUser(userId uint64) (*entity.User, error) {
	return getUserRepo(userId)
}

func (u *fakeUserRepo) GetUsers() ([]entity.User, error) {
	return getUsersRepo()
}

func (u *fakeUserRepo) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return getUserEmailAndPasswordRepo(user)
}

var userFakeApp UserAppInterface = &fakeUserRepo{}

func TestSaveUser_Success(t *testing.T) {
	saveUserRepo = func(user *entity.User) (*entity.User, map[string]string) {
		return &entity.User{
			ID:       1,
			FullName: "john",
			Email:    "john@example.com",
			Password: "123456",
		}, nil
	}

	user := &entity.User{
		ID:       1,
		FullName: "john",
		Email:    "john@example.com",
		Password: "123456",
	}
	u, err := userFakeApp.SaveUser(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FullName, "john")
	assert.EqualValues(t, u.Email, "john@example.com")
}

func TestGetUser_Success(t *testing.T) {
	getUserRepo = func(userId uint64) (*entity.User, error) {
		return &entity.User{
			ID:       1,
			FullName: "john",
			Email:    "john@example.com",
			Password: "123456",
		}, nil
	}
	userId := uint64(1)
	u, err := userFakeApp.GetUser(userId)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FullName, "john")
	assert.EqualValues(t, u.Email, "john@example.com")
	assert.EqualValues(t, u.Password, "123456")
}
