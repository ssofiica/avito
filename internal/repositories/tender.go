package repositories

import (
	"context"
	"errors"
	"strings"
	"zadanie-6105/internal/delivery/operation"
	"zadanie-6105/internal/repositories/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tender interface {
	Create(ctx context.Context, tender entities.Tender) (entities.Tender, error)
	GetTenderList(ctx context.Context, params operation.TenderListParams) (entities.TenderList, error)
	GetByUser(ctx context.Context, creator string, params operation.TenderListParams) (entities.TenderList, error)
	ChangeTenderStatus(ctx context.Context, status entities.TenderStatus, id string) (entities.Tender, error)
	GetTenderStatus(ctx context.Context, id string) (entities.TenderStatus, error)
	EditTender(ctx context.Context, tender entities.Tender, id string) (entities.Tender, error)
	CheckTenderOrganization(ctx context.Context, id string) (string, error)
}

type TenderRepo struct {
	db *pgxpool.Pool
}

func NewTenderRepo(db *pgxpool.Pool) Tender {
	return &TenderRepo{db: db}
}

func (t *TenderRepo) Create(ctx context.Context, tender entities.Tender) (entities.Tender, error) {
	query := `
		insert into tender(name, description, service_type, status, organization_id, creator_username, version, created_at) 
		values ($1, $2, $3, $4, $5, $6, 1, now()) 
		returning id, name, description, service_type, status, version, created_at
	`
	var res entities.Tender
	row := t.db.QueryRow(
		ctx, query,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		entities.TenderStatusCreated,
		tender.OrganizationID,
		tender.CreatorUsername,
	)
	err := row.Scan(&res.ID, &res.Name, &res.Description,
		&res.ServiceType, &res.Status, &res.Version, &res.CreatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *TenderRepo) GetTenderList(ctx context.Context, params operation.TenderListParams) (entities.TenderList, error) {
	var (
		res entities.TenderList
		sb  strings.Builder
	)
	sb.WriteString(
		`select id, name, description, service_type, status, version, created_at
		from tender order by name ASC`)
	queryFilters, args := t.inQuery(params)
	sb.WriteString(queryFilters)
	namedArgs := pgx.NamedArgs{
		"limit":  args["limit"],
		"offset": args["offset"],
	}
	rows, err := t.db.Query(ctx, sb.String(), namedArgs)
	if err != nil {
		return entities.TenderList{}, err
	}
	for rows.Next() {
		var tend entities.Tender
		err := rows.Scan(&tend.ID, &tend.Name, &tend.Description,
			&tend.ServiceType, &tend.Status, &tend.Version, &tend.CreatedAt)
		if err != nil {
			return entities.TenderList{}, err
		}
		res = append(res, tend)
	}
	return res, nil
}

func (t *TenderRepo) GetByUser(ctx context.Context, creator string, params operation.TenderListParams) (entities.TenderList, error) {
	var (
		res entities.TenderList
		sb  strings.Builder
	)
	sb.WriteString(
		`select id, name, description, service_type, status, version, created_at
		from tender where creator_username=@creator`)
	queryFilters, args := t.inQuery(params)
	sb.WriteString(queryFilters)
	namedArgs := pgx.NamedArgs{
		"limit":   args["limit"],
		"offset":  args["offset"],
		"creator": creator,
	}
	rows, err := t.db.Query(ctx, sb.String(), namedArgs)
	if err != nil {
		return entities.TenderList{}, err
	}
	for rows.Next() {
		var tend entities.Tender
		err := rows.Scan(&tend.ID, &tend.Name, &tend.Description,
			&tend.ServiceType, &tend.Status, &tend.Version, &tend.CreatedAt)
		if err != nil {
			return entities.TenderList{}, err
		}
		res = append(res, tend)
	}
	return res, nil
}

func (t *TenderRepo) ChangeTenderStatus(ctx context.Context, status entities.TenderStatus, id string) (entities.Tender, error) {
	query := `
		update tender set status=$1, updated_at=now() where id=$2
		returning id, name, description, service_type, status, version, created_at
	`
	var res entities.Tender
	row := t.db.QueryRow(ctx, query, status, id)
	err := row.Scan(&res.ID, &res.Name, &res.Description,
		&res.ServiceType, &res.Status, &res.Version, &res.CreatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *TenderRepo) GetTenderStatus(ctx context.Context, id string) (entities.TenderStatus, error) {
	query := `select status from tender where id=$1`
	var res entities.TenderStatus
	err := t.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *TenderRepo) EditTender(ctx context.Context, tender entities.Tender, id string) (entities.Tender, error) {
	var (
		res entities.Tender
		sb  strings.Builder
	)
	sb.WriteString(
		`update tender set version=version+1`)
	queryFilters, args := t.EditQuery(tender)
	sb.WriteString(queryFilters)
	sb.WriteString(` where id=@id returning id, name, description, service_type, status, version, created_at`)
	namedArgs := pgx.NamedArgs{
		"name":        args["name"],
		"description": args["description"],
		"service":     args["service"],
		"id":          id,
	}
	row := t.db.QueryRow(ctx, sb.String(), namedArgs)
	err := row.Scan(&res.ID, &res.Name, &res.Description,
		&res.ServiceType, &res.Status, &res.Version, &res.CreatedAt)
	if err != nil {
		return entities.Tender{}, err
	}
	return res, nil
}

func (t *TenderRepo) CheckTenderOrganization(ctx context.Context, id string) (string, error) {
	query := `select organization_id from tender where id=$1`
	var res string
	err := t.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return res, nil
}

func (t *TenderRepo) EditQuery(params entities.Tender) (string, map[string]any) {
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
	if params.ServiceType != "" {
		sb.WriteString(`, service_type=@service`)
		args["service"] = params.ServiceType
	}
	return sb.String(), args
}

func (t *TenderRepo) inQuery(params operation.TenderListParams) (string, map[string]any) {
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
	if params.Offset > 0 {
		sb.WriteString(` offset @offset`)
		args["offset"] = params.Offset
	}
	return sb.String(), args
}
