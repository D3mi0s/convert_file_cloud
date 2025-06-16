package repository

import (
	"file-upload-service/models"

	"gorm.io/gorm"
)

func (r *FileRepository) GetFilesByUser(userID uint) ([]models.File, error) {
	var files []models.File
	err := r.DB.Where("user_id = ?", userID).Find(&files).Error
	return files, err
}

type FileRepository struct {
	DB *gorm.DB
}

func (r *FileRepository) CreateFile(file *models.File) error {
	return r.DB.Create(file).Error
}

func (r *FileRepository) GetFileByID(id uint) (*models.File, error) {
	var file models.File
	err := r.DB.First(&file, id).Error
	return &file, err
}
