package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	UserID       uint   `gorm:"not null"`
	OriginalName string `gorm:"not null"`
	StoredName   string `gorm:"not null;unique"`
	Size         int64  `gorm:"not null"`
	MimeType     string `gorm:"not null"`
	Status       string `gorm:"default:'pending'"`
}
