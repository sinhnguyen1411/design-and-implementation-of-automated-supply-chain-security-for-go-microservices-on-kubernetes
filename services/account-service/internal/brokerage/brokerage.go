package brokerage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountType string
type AccountStatus string

const (
	TypeCash       AccountType = "cash"
	TypeMargin     AccountType = "margin"
	TypeRetirement AccountType = "retirement"

	StatusActive    AccountStatus = "active"
	StatusSuspended AccountStatus = "suspended"
	StatusClosed    AccountStatus = "closed"
)

type BrokerageAccount struct {
	ID            string        `json:"id"`
	OwnerID       string        `json:"owner_id"`
	Type          AccountType   `json:"type"`
	Status        AccountStatus `json:"status"`
	CashBalance   float64       `json:"cash_balance"`
	MarginBalance float64       `json:"margin_balance"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Transaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Kind      string    `json:"kind"`
	Amount    float64   `json:"amount"`
	Balance   float64   `json:"balance_after"`
	CreatedAt time.Time `json:"created_at"`
}

func New(ownerID string, t AccountType) BrokerageAccount {
	now := time.Now()
	return BrokerageAccount{
		ID:          uuid.NewString(),
		OwnerID:     ownerID,
		Type:        t,
		Status:      StatusActive,
		CashBalance: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (a *BrokerageAccount) Deposit(amount float64) (Transaction, error) {
	if amount <= 0 {
		return Transaction{}, errors.New("deposit amount must be positive")
	}
	if a.Status != StatusActive {
		return Transaction{}, errors.New("account not active")
	}
	a.CashBalance += amount
	a.UpdatedAt = time.Now()
	return Transaction{ID: uuid.NewString(), AccountID: a.ID, Kind: "deposit", Amount: amount, Balance: a.CashBalance, CreatedAt: a.UpdatedAt}, nil
}

func (a *BrokerageAccount) Withdraw(amount float64) (Transaction, error) {
	if amount <= 0 {
		return Transaction{}, errors.New("withdrawal amount must be positive")
	}
	if a.Status != StatusActive {
		return Transaction{}, errors.New("account not active")
	}
	if amount > a.CashBalance {
		return Transaction{}, errors.New("insufficient funds")
	}
	a.CashBalance -= amount
	a.UpdatedAt = time.Now()
	return Transaction{ID: uuid.NewString(), AccountID: a.ID, Kind: "withdrawal", Amount: amount, Balance: a.CashBalance, CreatedAt: a.UpdatedAt}, nil
}

func (a *BrokerageAccount) BuyingPower() float64 {
	if a.Type == TypeMargin {
		return a.CashBalance + a.MarginBalance
	}
	return a.CashBalance
}

func (a *BrokerageAccount) Suspend() { a.Status = StatusSuspended; a.UpdatedAt = time.Now() }
func (a *BrokerageAccount) Close()   { a.Status = StatusClosed; a.UpdatedAt = time.Now() }
