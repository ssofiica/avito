package repositories

import (
	"context"
	"zadanie-6105/internal/repositories/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tender interface {
	Create(ctx context.Context, tender entities.Tender, creator_id uint64) (entities.Tender, error)
	GetTenderList(ctx context.Context) (entities.TenderList, error)
	GetByUser(ctx context.Context, creator_id uint64) (entities.TenderList, error)
	ChangeStatus(ctx context.Context, status entities.TenderStatus, id uint64) (entities.Tender, error)
	GetTenderStatus(ctx context.Context, id uint64) (entities.TenderStatus, error)
}

type TenderRepo struct {
	db *pgxpool.Pool
}

func NewTenderRepo(db *pgxpool.Pool) TenderRepo {
	return TenderRepo{db: db}
}

func (t TenderRepo) Create(ctx context.Context, tender entities.Tender, creator_id uint64) (entities.Tender, error) {
	query := `
		insert into tender(name, description, service_type, status, organization_id, creator_username, created_at) 
		values ($1, $2, $3, $4, $5, $6, now()) 
		returning id, name, description, service_type, status, organizitaion_id
	`
	var res entities.Tender
	row := t.db.QueryRow(
		ctx, query,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		entities.TenderStatusCreated,
		tender.OrganizationID,
		creator_id,
	)
	//тут отсканировать надо creator_id
	err := row.Scan(&res.ID, &res.Name, &res.Description,
		&res.ServiceType, &res.Status, &res.OrganizationID, &res.CreatorUsername)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t TenderRepo) GetTenderList(ctx context.Context) (entities.TenderList, error) {
	//7 columns
	query := `
		select id, name, description, service_type, status, organization_id, creator_username
		from tender
	`
	var res entities.TenderList
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return entities.TenderList{}, err
	}
	for rows.Next() {
		var tend entities.Tender
		err := rows.Scan(&tend.ID, &tend.Name, &tend.Description,
			&tend.ServiceType, &tend.Status, &tend.OrganizationID, &tend.CreatorUsername)
		if err != nil {
			return entities.TenderList{}, err
		}
		res = append(res, tend)
	}
	return res, nil
}

// добавить limit offset
func (t TenderRepo) GetByUser(ctx context.Context, creator_id uint64) (entities.TenderList, error) {
	//8 columns
	query := `
		select id, name, description, service_type, status, organization_id, creator_username, created_at
		from tender where creator_id=$1
	`
	var res entities.TenderList
	rows, err := t.db.Query(ctx, query, creator_id)
	if err != nil {
		return entities.TenderList{}, err
	}
	for rows.Next() {
		var tend entities.Tender
		//тут отсканировать надо creator_id
		err := rows.Scan(&tend.ID, &tend.Name, &tend.Description,
			&tend.ServiceType, &tend.Status, &tend.OrganizationID, &tend.CreatorUsername, &tend.CreatedAt)
		if err != nil {
			return entities.TenderList{}, err
		}
		res = append(res, tend)
	}
	return res, nil
}

func (t TenderRepo) ChangeStatus(ctx context.Context, status entities.TenderStatus, id uint64) (entities.Tender, error) {
	// добавить получение юзера
	query := `
		update tender set status=$1 where id=$2
		returning id, name, description, service_type, status, organization_id, created_at
	`
	var res entities.Tender
	row := t.db.QueryRow(
		ctx, query,
		status, id,
	)
	//тут отсканировать надо creator_id
	err := row.Scan(&res.ID, &res.Name, &res.Description,
		&res.ServiceType, &res.Status, &res.OrganizationID, &res.CreatedAt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t TenderRepo) GetTenderStatus(ctx context.Context, id uint64) (entities.TenderStatus, error) {
	query := `select status from tender where id=$2`
	var res entities.TenderStatus
	err := t.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}
