package persistence

import (
	"fmt"

	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Repositories struct {
	User repository.UserRepository
	Drug repository.DrugRepository
	db   *gorm.DB
}

func NewRepositories(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: NewUserRepository(db),
		Drug: NewDrugRepository(db),
		db:   db,
	}, nil
}

func (repository *Repositories) Close() error {
	return repository.db.Close()
}

func (repository *Repositories) Automigrate() error {
	return repository.db.AutoMigrate(&entity.User{}, &entity.Drug{}).Error
}
