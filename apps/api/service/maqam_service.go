package service

import (
	"api/model"
	"api/repository"
)

type MaqamService struct {
	maqamRepo *repository.MaqamRepository
}

func NewMaqamService(maqamRepo *repository.MaqamRepository) *MaqamService {
	return &MaqamService{maqamRepo: maqamRepo}
}

func (s *MaqamService) GetAll() ([]model.Maqam, error) {
	return s.maqamRepo.FindAll()
}

func (s *MaqamService) GetByID(id string) (*model.Maqam, error) {
	return s.maqamRepo.FindByID(id)
}
