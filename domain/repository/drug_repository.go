package repository

import "github.com/akwanmaroso/ddd-drugs/domain/entity"

type FoodRepository interface {
	SaveDrug(*entity.Drug) (*entity.Drug, map[string]string)
	GetFood(uint64) (*entity.Drug, error)
	GetAllFood() ([]entity.Drug, error)
	UpdateFood(*entity.Drug) (*entity.Drug, map[string]string)
	DeleteFood(uint64) error
}
