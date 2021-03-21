package repository

import "github.com/akwanmaroso/ddd-drugs/domain/entity"

type DrugRepository interface {
	SaveDrug(*entity.Drug) (*entity.Drug, map[string]string)
	GetDrug(uint64) (*entity.Drug, error)
	GetAllDrug() ([]entity.Drug, error)
	UpdateDrug(*entity.Drug) (*entity.Drug, map[string]string)
	DeleteDrug(uint64) error
}
