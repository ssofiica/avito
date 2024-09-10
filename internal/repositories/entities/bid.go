package entities

import "time"

type BidStatus string

var (
	BidStatusCreated   BidStatus = "created"
	BidStatusCanceled  BidStatus = "canceled"
	BidStatusPublished BidStatus = "published"
)

type Bid struct {
	ID              uint64
	Name            string
	Description     string
	Status          BidStatus
	TenderID        uint64
	OrganizationId  uint64
	CreatorUsername string
	CreatedAt       time.Time
}

type BidList []Bid
