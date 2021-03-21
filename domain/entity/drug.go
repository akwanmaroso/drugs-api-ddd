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

func (drug *Drug) BeforeSave() {
	drug.Name = html.EscapeString(strings.TrimSpace(drug.Name))
}

func (drug *Drug) Prepare() {
	drug.Name = html.EscapeString(strings.TrimSpace(drug.Name))
	drug.CreatedAt = time.Now()
	drug.UpdatedAt = time.Now()
}
