package entities

import "time"

type TenderStatus string
type TenderType string

var (
	TenderStatusCreated    TenderStatus = "created"
	TenderStatusClosed     TenderStatus = "closed"
	TenderStatusPublished  TenderStatus = "published"
	TenderTypeConstruction TenderType   = "construction"
	TenderTypeDelivery     TenderType   = "delivery"
	TenderTypeManufacture  TenderType   = "manufacture"
)

func (status *TenderStatus) Scan(str string) {
	switch str {
	case "created":
		*status = TenderStatusCreated
	case "closed":
		*status = TenderStatusClosed
	case "published":
		*status = TenderStatusPublished
	default:
		*status = ""
	}
}

type Tender struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	ServiceType     TenderType   `json:"service_type,omitempty"`
	Status          TenderStatus `json:"status"`
	OrganizationID  string       `json:"organization_id,omitempty"`
	CreatorUsername string       `json:"username"`
	Version         uint8        `json:"version"`
	CreatedAt       time.Time    `json:"created_at"`
}

type TenderList []Tender
