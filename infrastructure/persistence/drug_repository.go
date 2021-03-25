package persistence

import (
	"errors"
	"strings"

	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/domain/repository"
	"github.com/jinzhu/gorm"
)

type DrugRepo struct {
	db *gorm.DB
}

func NewDrugRepository(db *gorm.DB) *DrugRepo {
	return &DrugRepo{db}
}

// DrugRepo implement the contract of repository.DrugRepository interface
var _ repository.DrugRepository = &DrugRepo{}

func (repo *DrugRepo) SaveDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	dbErr := map[string]string{}
	err := repo.db.Debug().Create(&drug).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_name"] = "drug name already taken"
			return nil, dbErr
		}
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return drug, nil
}

func (repo *DrugRepo) GetDrug(drugId uint64) (*entity.Drug, error) {
	var drug entity.Drug
	err := repo.db.Debug().Where("id = ?", drugId).Take(&drug).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("drug not found")
	}
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	return &drug, nil
}

func (repo *DrugRepo) GetAllDrug() ([]entity.Drug, error) {
	var drugs []entity.Drug
	err := repo.db.Debug().Limit(100).Order("created_at desc").Find(&drugs).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("Drug not found")
	}
	return drugs, nil
}

func (repo *DrugRepo) UpdateDrug(drug *entity.Drug) (*entity.Drug, map[string]string) {
	dbErr := map[string]string{}
	err := repo.db.Save(&drug).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_name"] = "name already taken"
			return nil, dbErr
		}
	}
	return drug, nil
}

func (repo *DrugRepo) DeleteDrug(drugId uint64) error {
	var drug entity.Drug
	err := repo.db.Where("id = ?", drugId).Delete(&drug).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
