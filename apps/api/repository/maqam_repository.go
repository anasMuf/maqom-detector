package repository

import (
	"api/model"

	"gorm.io/gorm"
)

type MaqamRepository struct {
	db *gorm.DB
}

func NewMaqamRepository(db *gorm.DB) *MaqamRepository {
	return &MaqamRepository{db: db}
}

func (r *MaqamRepository) FindAll() ([]model.Maqam, error) {
	var maqams []model.Maqam
	err := r.db.Order("name_latin ASC").Find(&maqams).Error
	return maqams, err
}

func (r *MaqamRepository) FindByID(id string) (*model.Maqam, error) {
	maqam := &model.Maqam{}
	err := r.db.Where("id = ?", id).First(maqam).Error
	if err != nil {
		return nil, err
	}
	return maqam, nil
}
