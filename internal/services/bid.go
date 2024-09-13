package services

import (
	"context"
	"zadanie-6105/internal/delivery/operation"
	"zadanie-6105/internal/repositories"
	"zadanie-6105/internal/repositories/entities"
)

type Bid interface {
	GetUserBids(ctx context.Context, params operation.BidParams, username string) (entities.BidList, error)
	CreateBid(ctx context.Context, bid entities.Bid) (entities.Bid, error)
	GetBidsForTender(ctx context.Context, tendor_id string, params operation.BidParams) (entities.BidList, error)
	GetBidStatus(ctx context.Context, id string) (entities.BidStatus, error)
	ChangeBidStatus(ctx context.Context, status entities.BidStatus, id string) (entities.Bid, error)
	SubmitBid(ctx context.Context, decision entities.BidStatus, bid_id string, username string) (entities.Bid, error)
	EditBid(ctx context.Context, bid entities.Bid, id string) (entities.Bid, error)
}

type BidService struct {
	repo   repositories.Bid
	user   repositories.User
	tender repositories.Tender
}

func NewBidService(repo repositories.Bid, user repositories.User, tender repositories.Tender) Bid {
	return &BidService{
		repo:   repo,
		user:   user,
		tender: tender,
	}
}

func (s *BidService) GetUserBids(ctx context.Context, params operation.BidParams, username string) (entities.BidList, error) {
	id, err := s.user.GetUserIDByUsername(ctx, username)
	if err != nil {
		return entities.BidList{}, err
	}
	return s.repo.GetUserBids(ctx, params, id)
}

func (s *BidService) CreateBid(ctx context.Context, bid entities.Bid) (entities.Bid, error) {
	return s.repo.Create(ctx, bid)
}

func (s *BidService) GetBidsForTender(ctx context.Context, tendor_id string, params operation.BidParams) (entities.BidList, error) {
	return s.repo.GetByTender(ctx, tendor_id, params)
}

func (s *BidService) GetBidStatus(ctx context.Context, id string) (entities.BidStatus, error) {
	return s.repo.GetBidStatus(ctx, id)
}

func (s *BidService) ChangeBidStatus(ctx context.Context, status entities.BidStatus, id string) (entities.Bid, error) {
	return s.repo.ChangeBidStatus(ctx, status, id)
}

// username это логин автора тендера!
func (s *BidService) SubmitBid(ctx context.Context, decision entities.BidStatus, bid_id string, username string) (entities.Bid, error) {
	organizationId, err := s.user.IsResponsible(ctx, username) //0 - значит юзер не остветсвенен за организацию
	if err != nil {
		return entities.Bid{}, err
	}
	if organizationId == "" {
		return entities.Bid{}, ErrNotResponsible
	}
	id, err := s.repo.GetTenderIDForBid(ctx, bid_id)
	if err != nil {
		return entities.Bid{}, err
	}
	tenderOrganizationId, err := s.tender.CheckTenderOrganization(ctx, id)
	if err != nil {
		return entities.Bid{}, err
	}
	if tenderOrganizationId != organizationId {
		return entities.Bid{}, ErrNoAccess
	}
	bid, err := s.repo.ChangeBidStatus(ctx, entities.BidStatusApproved, bid_id)
	if err != nil {
		return entities.Bid{}, err
	}
	_, err = s.tender.ChangeTenderStatus(ctx, entities.TenderStatusClosed, id)
	if err != nil {
		return entities.Bid{}, err
	}
	return bid, nil
}

func (s *BidService) EditBid(ctx context.Context, bid entities.Bid, id string) (entities.Bid, error) {
	return s.repo.EditBid(ctx, bid, id)
}
