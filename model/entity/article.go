package entity

import (
	"gorm.io/gorm"
	"time"
)

//Article struct is a model that represent articles table in database
type Article struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Title string `gorm:"type:varchar(255)" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	AuthorID uint64 `gorm:"not null" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Author User `gorm:"foreignKey:AuthorID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"author"`
}