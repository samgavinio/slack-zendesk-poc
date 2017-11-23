package models

import "time"

type Model struct {
	ID uint `json:"id" gorm:"primary_key" validate:"len=0"`
	CreatedAt time.Time `json:"created_at" sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" sql:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}
