package account

import "github.com/google/uuid"

type KYCStatus string

const (
	KYCPending  KYCStatus = "pending"
	KYCVerified KYCStatus = "verified"
	KYCRejected KYCStatus = "rejected"
)

type Account struct {
	ID      string    `json:"id"`
	OwnerID string    `json:"owner_id"`
	Status  KYCStatus `json:"status"`
	Tier    string    `json:"tier"`
}

func New(ownerID string) Account {
	return Account{ID: uuid.NewString(), OwnerID: ownerID, Status: KYCPending, Tier: "standard"}
}

func (a *Account) Verify() { a.Status = KYCVerified }
func (a *Account) Reject() { a.Status = KYCRejected }
