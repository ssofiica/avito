package entities

import "time"

type TenderStatus string

var (
	TenderStatusCreated   TenderStatus = "created"
	TenderStatusClosed    TenderStatus = "closed"
	TenderStatusPublished TenderStatus = "published"
)

type Tender struct {
	ID              uint64
	Name            string
	Description     string
	ServiceType     string
	Status          TenderStatus
	OrganizationID  uint64
	CreatorUsername string
	CreatedAt       time.Time
}

type TenderList []Tender
