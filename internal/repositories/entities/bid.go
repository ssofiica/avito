package entities

import "time"

type BidStatus string

var (
	BidStatusCreated   BidStatus = "created"
	BidStatusCanceled  BidStatus = "canceled"
	BidStatusPublished BidStatus = "published"
	BidStatusApproved  BidStatus = "approved"
	BidStatusRejected  BidStatus = "rejected"
)

func (status *BidStatus) Scan(str string) {
	switch str {
	case "created":
		*status = BidStatusCreated
	case "canceled":
		*status = BidStatusCanceled
	case "published":
		*status = BidStatusPublished
	case "approved":
		*status = BidStatusApproved
	case "rejected":
		*status = BidStatusRejected
	default:
		*status = ""
	}
}

type Bid struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      BidStatus `json:"status"`
	AuthorType  string    `json:"author_type"`
	AuthorID    string    `json:"author_id"`
	TenderID    string    `json:"tender_id,omitempty"`
	Version     uint8     `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
}

type BidList []Bid
