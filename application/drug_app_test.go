package application

import (
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeDrugRepo struct{}

var (
	saveDrugRepo   func(drug *entity.Drug) (*entity.Drug, map[string]string)
	getDrugRepo    func(uint64) (*entity.Drug, error)
	getAllDrugRepo func() ([]entity.Drug, error)
	updateDrugRepo func(*entity.Drug) (*entity.Drug, map[string]string)
	deleteDrugRepo func(uint64) error
)

func (f *fakeDrugRepo) SaveDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	return saveDrugRepo(drug)
}

func (f *fakeDrugRepo) GetAllDrug() ([]entity.Drug, error) {
	return getAllDrugRepo()
}

func (f *fakeDrugRepo) GetDrug(drugId uint64) (*entity.Drug, error) {
	return getDrugRepo(drugId)
}

func (f *fakeDrugRepo) UpdateDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	return updateDrugRepo(drug)
}

func (f *fakeDrugRepo) DeleteDrug(drugId uint64) error {
	return deleteDrugRepo(drugId)
}

var drugAppFake DrugAppInterface = &fakeDrugRepo{}

func TestSaveDrug_Success(t *testing.T) {
	saveDrugRepo = func(drug *entity.Drug) (*entity.Drug, map[string]string) {
		return &entity.Drug{
			ID:          1,
			Name:        "drug title",
			Description: "drug description",
			UserID:      1,
		}, nil
	}
	drug := &entity.Drug{
		ID:          1,
		Name:        "drug title",
		Description: "drug description",
		UserID:      1,
	}
	d, err := drugAppFake.SaveDrug(drug)
	assert.Nil(t, err)
	assert.EqualValues(t, d.Name, "drug title")
	assert.EqualValues(t, d.Description, "drug description")
	assert.EqualValues(t, d.UserID, 1)
}
