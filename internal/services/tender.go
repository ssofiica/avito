package services

import (
	"context"
	"errors"
	"zadanie-6105/internal/delivery/operation"
	"zadanie-6105/internal/repositories"
	"zadanie-6105/internal/repositories/entities"
)

var (
	ErrNotResponsible = errors.New("пользователь не связан с организацией")
	ErrNoAccess       = errors.New("нет доступа, пользователь не отвественен за тендер")
)

type Tender interface {
	GetTenderList(ctx context.Context, params operation.TenderListParams) (entities.TenderList, error)
	CreateTender(ctx context.Context, tender entities.Tender) (entities.Tender, error)
	GetTenderByUser(ctx context.Context, creator string, params operation.TenderListParams) (entities.TenderList, error)
	GetTenderStatus(ctx context.Context, id string, username string) (entities.TenderStatus, error)
	ChangeTenderStatus(ctx context.Context, status entities.TenderStatus, id string, username string) (entities.Tender, error)
	EditTender(ctx context.Context, tender entities.Tender, id string, username string) (entities.Tender, error)
}

type TenderService struct {
	repo repositories.Tender
	user repositories.User
}

func NewTenderService(repo repositories.Tender, user repositories.User) Tender {
	return &TenderService{
		repo: repo,
		user: user,
	}
}

func (s *TenderService) GetTenderList(ctx context.Context, params operation.TenderListParams) (entities.TenderList, error) {
	return s.repo.GetTenderList(ctx, params)
}

func (s *TenderService) CreateTender(ctx context.Context, tender entities.Tender) (entities.Tender, error) {
	organizationId, err := s.user.IsResponsible(ctx, tender.CreatorUsername) //"" - значит юзер не остветсвенен за организацию
	if err != nil {
		return entities.Tender{}, err
	}
	if organizationId == "" {
		return entities.Tender{}, ErrNotResponsible
	}
	if organizationId != tender.OrganizationID {
		return entities.Tender{}, ErrNoAccess
	}
	return s.repo.Create(ctx, tender)
}

func (s *TenderService) GetTenderByUser(ctx context.Context, creator string, params operation.TenderListParams) (entities.TenderList, error) {
	return s.repo.GetByUser(ctx, creator, params)
}

func (s *TenderService) GetTenderStatus(ctx context.Context, id string, username string) (entities.TenderStatus, error) {
	return s.repo.GetTenderStatus(ctx, id)
}

func (s *TenderService) ChangeTenderStatus(ctx context.Context, status entities.TenderStatus, id string, username string) (entities.Tender, error) {
	organizationId, err := s.user.IsResponsible(ctx, username) //"" - значит юзер не остветсвенен за организацию
	if err != nil {
		return entities.Tender{}, err
	}
	if organizationId == "" {
		return entities.Tender{}, ErrNotResponsible
	}
	tenderOrganizationId, err := s.repo.CheckTenderOrganization(ctx, id)
	if err != nil {
		return entities.Tender{}, err
	}
	if tenderOrganizationId != organizationId {
		return entities.Tender{}, ErrNoAccess
	}
	return s.repo.ChangeTenderStatus(ctx, status, id)
}

func (s *TenderService) EditTender(ctx context.Context, tender entities.Tender, id string, username string) (entities.Tender, error) {
	organizationId, err := s.user.IsResponsible(ctx, username) //0 - значит юзер не остветсвенен за организацию
	if err != nil {
		return entities.Tender{}, err
	}
	if organizationId == "" {
		return entities.Tender{}, ErrNotResponsible
	}
	tenderOrganizationId, err := s.repo.CheckTenderOrganization(ctx, id)
	if err != nil {
		return entities.Tender{}, err
	}
	if tenderOrganizationId != organizationId {
		return entities.Tender{}, ErrNoAccess
	}
	return s.repo.EditTender(ctx, tender, id)
}
