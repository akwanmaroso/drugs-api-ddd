package application

import (
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/domain/repository"
)

type drugApps struct {
	drugApp repository.DrugRepository
}

type DrugAppInterface interface {
	SaveDrug(*entity.Drug) (*entity.Drug, map[string]string)
	GetAllDrug() ([]entity.Drug, error)
	GetDrug(uint64) (*entity.Drug, error)
	UpdateDrug(*entity.Drug) (*entity.Drug, map[string]string)
	DeleteDrug(uint64) error
}

// implement contract drug-repository
var _ DrugAppInterface = &drugApps{}

func (d *drugApps) SaveDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	return d.drugApp.SaveDrug(drug)
}

func (d *drugApps) GetAllDrug() ([]entity.Drug, error) {
	return d.drugApp.GetAllDrug()
}

func (d *drugApps) GetDrug(drugId uint64) (*entity.Drug, error) {
	return d.drugApp.GetDrug(drugId)
}

func (d *drugApps) UpdateDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	return d.drugApp.UpdateDrug(drug)
}

func (d *drugApps) DeleteDrug(drugId uint64) error {
	return d.DeleteDrug(drugId)
}
