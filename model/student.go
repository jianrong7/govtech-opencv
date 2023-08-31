package model

import "gorm.io/gorm"

// Student struct
type Student struct {
	gorm.Model
	// ID          uint      `gorm:"primaryKey" json:"id"`
	Email       string    `gorm:"not null" json:"email"`
	IsSuspended bool      `gorm:"not null" json:"isSuspended"`
	Teachers    []Teacher `gorm:"many2many:teacher_students;" json:"teachers"`
}
