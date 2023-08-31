package model

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email       string    `gorm:"not null" json:"email"`
	IsSuspended bool      `gorm:"not null" json:"isSuspended"`
	Teachers    []Teacher `gorm:"many2many:teacher_students;" json:"teachers"`
}
