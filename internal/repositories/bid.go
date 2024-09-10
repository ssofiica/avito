package repositories

import (
	"context"
	"zadanie-6105/internal/repositories/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Bid interface {
	Create(ctx context.Context, bid entities.Bid, creator_id uint64) (entities.Bid, error)
	GetByTender(ctx context.Context, bider_id uint64) (entities.BidList, error)
	GetByUser(ctx context.Context, creator_id uint64) (entities.BidList, error)
}

type BidRepo struct {
	db *pgxpool.Pool
}

func NewBidRepo(db *pgxpool.Pool) Bid {
	return BidRepo{db: db}
}

func (t BidRepo) Create(ctx context.Context, bid entities.Bid, creator_id uint64) (entities.Bid, error) {
	query := `
		insert into bid(name, description, status, tender_id, organiztion_id, creator_username, created_at) 
		values ($1, $2, $3, $4, $5, $6, now()) 
		returning id, name, description, status, tender_id, organization_id, created_at
	`
	var res entities.Bid
	row := t.db.QueryRow(
		ctx, query,
		bid.Name,
		bid.Description,
		entities.BidStatusCreated,
		bid.TenderID,
		bid.OrganizationId,
		creator_id,
	)
	//тут отсканировать надо creator_id
	err := row.Scan(&res.ID, &res.Name, &res.Description,
		&res.Status, &res.TenderID, &res.OrganizationId, &res.CreatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t BidRepo) GetByTender(ctx context.Context, bider_id uint64) (entities.BidList, error) {
	//7 columns
	query := `
		select id, name, description, status, tender_id, organization_id, creator_at
		from bid where tender_id=$1
	`
	var res entities.BidList
	rows, err := t.db.Query(ctx, query, bider_id)
	if err != nil {
		return entities.BidList{}, err
	}
	for rows.Next() {
		var bid entities.Bid
		//тут отсканировать надо creator_id
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description,
			&bid.Status, &bid.TenderID, &bid.OrganizationId, &bid.CreatedAt)
		if err != nil {
			return entities.BidList{}, err
		}
		res = append(res, bid)
	}
	return res, nil
}

// добавить limit offset
func (t BidRepo) GetByUser(ctx context.Context, creator_id uint64) (entities.BidList, error) {
	//8 columns
	query := `
		select id, name, description, status, tender_id, organiztion_id, created_at
		from bid where creator_id=$1
	`
	var res entities.BidList
	rows, err := t.db.Query(ctx, query, creator_id)
	if err != nil {
		return entities.BidList{}, err
	}
	for rows.Next() {
		var bid entities.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Description,
			&bid.Status, &bid.TenderID, &bid.OrganizationId, &bid.CreatedAt)
		if err != nil {
			return entities.BidList{}, err
		}
		res = append(res, bid)
	}
	return res, nil
}
