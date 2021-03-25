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

func TestGetAllDrug_Success(t *testing.T) {
	getAllDrugRepo = func() ([]entity.Drug, error) {
		return []entity.Drug{
			{
				ID:          1,
				Name:        "drug name first",
				Description: "drug description first",
				UserID:      1,
			}, {
				ID:          2,
				Name:        "drug name second",
				Description: "drug description second",
				UserID:      2,
			},
		}, nil
	}
	drug, err := drugAppFake.GetAllDrug()
	assert.Nil(t, err)
	assert.EqualValues(t, len(drug), 2)
}

func TestUpdateDrug_Success(t *testing.T) {
	updateDrugRepo = func(drug *entity.Drug) (*entity.Drug, map[string]string) {
		return &entity.Drug{
			ID:          1,
			Name:        "drug title update",
			Description: "drug description update",
			UserID:      1,
		}, nil
	}
	drug := &entity.Drug{
		ID:          1,
		Name:        "drug title update",
		Description: "drug description update",
		UserID:      1,
	}
	d, err := drugAppFake.UpdateDrug(drug)
	assert.Nil(t, err)
	assert.EqualValues(t, d.Name, "drug title update")
	assert.EqualValues(t, d.Description, "drug description update")
	assert.EqualValues(t, d.UserID, 1)
}

func TestDeleteDrug_Success(t *testing.T) {
	deleteDrugRepo = func(drugId uint64) error {
		return nil
	}

	drugId := uint64(1)
	err := drugAppFake.DeleteDrug(drugId)
	assert.Nil(t, err)
}
