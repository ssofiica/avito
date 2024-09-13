package repositories

import (
	"context"
	"strings"
	"zadanie-6105/internal/delivery/operation"
	"zadanie-6105/internal/repositories/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bid interface {
	Create(ctx context.Context, bid entities.Bid) (entities.Bid, error)
	GetByTender(ctx context.Context, tender_id string, params operation.BidParams) (entities.BidList, error)
	GetUserBids(ctx context.Context, params operation.BidParams, id int) (entities.BidList, error)
	GetBidStatus(ctx context.Context, id string) (entities.BidStatus, error)
	ChangeBidStatus(ctx context.Context, status entities.BidStatus, id string) (entities.Bid, error)
	EditBid(ctx context.Context, bid entities.Bid, id string) (entities.Bid, error)
	GetTenderIDForBid(ctx context.Context, id string) (string, error)
}

type BidRepo struct {
	db *pgxpool.Pool
}

func NewBidRepo(db *pgxpool.Pool) Bid {
	return &BidRepo{db: db}
}

func (t *BidRepo) Create(ctx context.Context, bid entities.Bid) (entities.Bid, error) {
	query := `
		insert into bid(name, description, status, tender_id, author_id, author_type, version, created_at) 
		values ($1, $2, $3, $4, $5, $6, 1, now()) 
		returning id, name, status, author_type, author_id, version, created_at
	`
	var res entities.Bid
	row := t.db.QueryRow(
		ctx, query,
		bid.Name,
		bid.Description,
		entities.BidStatusCreated,
		bid.TenderID,
		bid.AuthorID,
		bid.AuthorType,
	)
	err := row.Scan(
		&res.ID,
		&res.Name,
		&res.Status,
		&res.AuthorType,
		&res.AuthorID,
		&res.Version,
		&res.CreatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *BidRepo) GetByTender(ctx context.Context, tender_id string, params operation.BidParams) (entities.BidList, error) {
	var (
		res entities.BidList
		sb  strings.Builder
	)
	sb.WriteString(
		`select id, name, status, author_type, author_id, version, created_at from bid 
		where tender_id=@id`)
	queryFilters, args := t.inQuery(params)
	sb.WriteString(queryFilters)
	namedArgs := pgx.NamedArgs{
		"limit":  args["limit"],
		"offset": args["offset"],
		"id":     tender_id,
	}
	rows, err := t.db.Query(ctx, sb.String(), namedArgs)
	if err != nil {
		return entities.BidList{}, err
	}
	for rows.Next() {
		var bid entities.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Status, &bid.AuthorType,
			&bid.AuthorID, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return entities.BidList{}, err
		}
		res = append(res, bid)
	}
	return res, nil
}

func (t *BidRepo) GetUserBids(ctx context.Context, params operation.BidParams, id int) (entities.BidList, error) {
	var (
		res entities.BidList
		sb  strings.Builder
	)
	sb.WriteString(
		`select id, name, status, authorType, author_id, version, created_at
		from bid where author_id=@id order by name ASC`)
	queryFilters, args := t.inQuery(params)
	sb.WriteString(queryFilters)
	namedArgs := pgx.NamedArgs{
		"limit":  args["limit"],
		"offset": args["offset"],
		"id":     id,
	}
	rows, err := t.db.Query(ctx, sb.String(), namedArgs)
	if err != nil {
		return entities.BidList{}, err
	}
	for rows.Next() {
		var bid entities.Bid
		err := rows.Scan(&bid.ID, &bid.Name, &bid.Status,
			&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return entities.BidList{}, err
		}
		res = append(res, bid)
	}
	return res, nil
}

func (t *BidRepo) ChangeBidStatus(ctx context.Context, status entities.BidStatus, id string) (entities.Bid, error) {
	query := `
		update bid set status=$1, updated_at=now() where id=$2
		returning id, name, status, author_type, author_id, version, created_at
	`
	var res entities.Bid
	row := t.db.QueryRow(ctx, query, status, id)
	err := row.Scan(&res.ID, &res.Name, &res.Status,
		&res.AuthorType, &res.AuthorID, &res.Version, &res.CreatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *BidRepo) GetBidStatus(ctx context.Context, id string) (entities.BidStatus, error) {
	query := `select status from bid where id=$1`
	var res entities.BidStatus
	err := t.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *BidRepo) GetTenderIDForBid(ctx context.Context, id string) (string, error) {
	query := `select tender_id from bid where id=$1`
	var res string
	err := t.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *BidRepo) inQuery(params operation.BidParams) (string, map[string]any) {
	var (
		args map[string]any = map[string]any{
			"limit":  5,
			"offset": 0,
		}
		sb strings.Builder
	)
	sb.WriteString(` limit @limit`)
	if params.Limit > 0 {
		args["limit"] = params.Limit
	}
	sb.WriteString(` offset @offset`)
	if params.Offset > 0 {
		args["offset"] = params.Offset
	}
	return sb.String(), args
}

func (t *BidRepo) EditBid(ctx context.Context, bid entities.Bid, id string) (entities.Bid, error) {
	var (
		res entities.Bid
		sb  strings.Builder
	)
	sb.WriteString(
		`update bid set version=version+1`)
	queryFilters, args := t.EditQuery(bid)
	sb.WriteString(queryFilters)
	sb.WriteString(` where id=@id returning id, name, status, author_type, author_id, version, created_at`)
	namedArgs := pgx.NamedArgs{
		"name":        args["name"],
		"description": args["description"],
		"id":          id,
	}
	row := t.db.QueryRow(ctx, sb.String(), namedArgs)
	err := row.Scan(&res.ID, &res.Name, &res.Status,
		&res.AuthorType, &res.AuthorID, &res.Version, &res.CreatedAt)
	if err != nil {
		return entities.Bid{}, err
	}
	return res, nil
}

func (t *BidRepo) EditQuery(params entities.Bid) (string, map[string]any) {
	args := map[string]any{}
	var sb strings.Builder
	if params.Name != "" {
		sb.WriteString(`, name=@name`)
		args["name"] = params.Name
	}
	if params.Description != "" {
		sb.WriteString(`, description=@description`)
		args["description"] = params.Description
	}
	return sb.String(), args
}
