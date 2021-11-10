package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"type varchar(50); not null" json:"name"`
	Price       float64 `gorm:"float; not null" json:"price"`
	Description *string `gorm:"type varchar(100)" json:"description"`
	Img         *string `gorm:"type varchar(100)" json:"img"`
}

type User struct {
	gorm.Model
	Email    string `gorm:"type varchar(100); not null; unique" json:"email"`
	Password string `gorm:"type varchar(64); not null" json:"password"`
}
