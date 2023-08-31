package model

import "gorm.io/gorm"

// Teacher struct
type Teacher struct {
	gorm.Model
	// ID       uint      `gorm:"primaryKey" json:"id"`
	Email    string    `gorm:"not null" json:"email"`
	Students []Student `gorm:"many2many:teacher_students;" json:"students"`
}
