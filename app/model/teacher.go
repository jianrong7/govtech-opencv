package model

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Email    string    `gorm:"not null" json:"email"`
	Students []Student `gorm:"many2many:teacher_students;" json:"students"`
}
