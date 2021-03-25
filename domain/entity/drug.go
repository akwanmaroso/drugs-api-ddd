package entity

import (
	"html"
	"strings"
	"time"
)

type Drug struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64     `gorm:"size:100;not null" json:"user_id"`
	Name        string     `gorm:"size:100;not null;unique" json:"name"`
	Description string     `gorm:"text;not null" json:"description"`
	DrugImage   string     `gorm:"size:255;null;unique" json:"drug_image"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// BeforeSave is responsible for memorizing the password before saving it to db
func (drug *Drug) BeforeSave() {
	drug.Name = html.EscapeString(strings.TrimSpace(drug.Name))
}

// Prepare is in charge of ensuring the data entered into the db is correct
func (drug *Drug) Prepare() {
	drug.Name = html.EscapeString(strings.TrimSpace(drug.Name))
	drug.CreatedAt = time.Now()
	drug.UpdatedAt = time.Now()
}

// Validate is responsible for validating the data you want to save to db
func (drug *Drug) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if drug.Name == "" || drug.Name == "null" {
			errorMessages["name_required"] = "name is required"
		}
		if drug.Description == "" || drug.Description == "null" {
			errorMessages["description_required"] = "description is required"
		}
		if drug.DrugImage == "" || drug.DrugImage == "null" {
			errorMessages["drug_image_required"] = "drug image is required"
		}
	default:
		if drug.Name == "" || drug.Name == "null" {
			errorMessages["name_required"] = "name is required"
		}
		if drug.Description == "" || drug.Description == "null" {
			errorMessages["description_required"] = "description is required"
		}
		if drug.DrugImage == "" || drug.DrugImage == "null" {
			errorMessages["drug_image_required"] = "drug image is required"
		}
	}
	return errorMessages
}
