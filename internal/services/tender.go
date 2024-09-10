package services

import (
	"context"
	"zadanie-6105/internal/repositories"
	"zadanie-6105/internal/repositories/entities"
)

type TenderService struct {
	repo repositories.Tender
}

func NewTenderService(repo repositories.Tender) TenderService {
	return TenderService{
		repo: repo,
	}
}

func (s TenderService) GetTenderList(ctx context.Context) (entities.TenderList, error) {
	return s.repo.GetTenderList(ctx)
}
