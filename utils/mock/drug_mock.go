package mock

import "github.com/akwanmaroso/ddd-drugs/domain/entity"

// DrugAppInterface is a mock drug app interface
type DrugAppInterface struct {
	SaveDrugFn   func(*entity.Drug) (*entity.Drug, map[string]string)
	GetAllDrugFn func() ([]entity.Drug, error)
	GetDrugFn    func(uint64) (*entity.Drug, error)
	UpdateDrugFn func(*entity.Drug) (*entity.Drug, map[string]string)
	DeleteDrugFn func(uint64) error
}

// SaveDrug calls the SaveDrugFn
func (d *DrugAppInterface) SaveDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	return d.SaveDrugFn(drug)
}

// GetAllDrug calls the GetAllDrugFn
func (d *DrugAppInterface) GetAllDrug() ([]entity.Drug, error) {
	return d.GetAllDrugFn()
}

// GetDrug calls the GetDrugFn
func (d *DrugAppInterface) GetDrug(drugId uint64) (*entity.Drug, error) {
	return d.GetDrugFn(drugId)
}

// UpdateDrug calls the UpdateDrugFn
func (d *DrugAppInterface) UpdateDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	return d.UpdateDrugFn(drug)
}

// DeleteDrug calls the DeleteDrugFn
func (d *DrugAppInterface) DeleteDrug(drugId uint64) error {
	return d.DeleteDrugFn(drugId)
}
