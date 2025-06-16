package repository

import (
	"file-conversion-service/models"

	"gorm.io/gorm"
)

type FileRepository struct {
	DB *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{DB: db}
}

func (r *FileRepository) CreateFile(file *models.File) error {
	return r.DB.Create(file).Error
}

func (r *FileRepository) GetFileByID(id uint) (*models.File, error) {
	var file models.File
	err := r.DB.First(&file, id).Error
	return &file, err
}

func (r *FileRepository) UpdateFileStatus(id uint, status string) error {
	return r.DB.Model(&models.File{}).Where("id = ?", id).Update("status", status).Error
}

func (r *FileRepository) UpdateConvertedName(id uint, convertedName string) error {
	return r.DB.Model(&models.File{}).Where("id = ?", id).Update("converted_name", convertedName).Error
}
