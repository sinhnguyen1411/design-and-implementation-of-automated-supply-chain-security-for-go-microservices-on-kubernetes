package kyc

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusVerified Status = "verified"
	StatusRejected Status = "rejected"
)

type DocType string

const (
	DocPassport      DocType = "passport"
	DocDriverLicense DocType = "driver_license"
	DocNationalID    DocType = "national_id"
	DocProofAddress  DocType = "proof_of_address"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type Document struct {
	ID         string     `json:"id"`
	Type       DocType    `json:"type"`
	Verified   bool       `json:"verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type Profile struct {
	UserID       string     `json:"user_id"`
	Status       Status     `json:"status"`
	RiskLevel    RiskLevel  `json:"risk_level"`
	RiskScore    int        `json:"risk_score"` // 0-100
	Jurisdiction string     `json:"jurisdiction"`
	Documents    []Document `json:"documents"`
	Reason       string     `json:"reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
}

type AMLFlag struct {
	UserID    string    `json:"user_id"`
	Code      string    `json:"code"` // "PEP", "SANCTIONS", "UNUSUAL_ACTIVITY"
	Severity  RiskLevel `json:"severity"`
	Notes     string    `json:"notes"`
	FlaggedAt time.Time `json:"flagged_at"`
}

type Registry struct {
	mu       sync.RWMutex
	profiles map[string]*Profile
	flags    []AMLFlag
}

func NewRegistry() *Registry {
	return &Registry{profiles: make(map[string]*Profile)}
}

var ErrNotFound = errors.New("kyc profile not found")

func (r *Registry) Create(userID, jurisdiction string) Profile {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := &Profile{
		UserID:       userID,
		Status:       StatusPending,
		RiskLevel:    RiskLow,
		RiskScore:    0,
		Jurisdiction: jurisdiction,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	r.profiles[userID] = p
	return *p
}

func (r *Registry) AddDocument(userID string, t DocType, expiresAt *time.Time) (Document, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.profiles[userID]
	if !ok {
		return Document{}, ErrNotFound
	}
	doc := Document{ID: uuid.NewString(), Type: t, ExpiresAt: expiresAt}
	p.Documents = append(p.Documents, doc)
	p.UpdatedAt = time.Now().UTC()
	return doc, nil
}

func (r *Registry) Approve(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.profiles[userID]
	if !ok {
		return ErrNotFound
	}
	now := time.Now().UTC()
	for i := range p.Documents {
		p.Documents[i].Verified = true
		p.Documents[i].VerifiedAt = &now
	}
	p.Status = StatusVerified
	p.ReviewedAt = &now
	p.UpdatedAt = now
	return nil
}

func (r *Registry) Reject(userID, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.profiles[userID]
	if !ok {
		return ErrNotFound
	}
	now := time.Now().UTC()
	p.Status = StatusRejected
	p.Reason = reason
	p.ReviewedAt = &now
	p.UpdatedAt = now
	return nil
}

func (r *Registry) FlagAML(userID, code, notes string, severity RiskLevel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.profiles[userID]
	if !ok {
		return ErrNotFound
	}
	r.flags = append(r.flags, AMLFlag{
		UserID:    userID,
		Code:      code,
		Severity:  severity,
		Notes:     notes,
		FlaggedAt: time.Now().UTC(),
	})
	switch severity {
	case RiskHigh:
		p.RiskScore = clamp(p.RiskScore+40, 0, 100)
	case RiskMedium:
		p.RiskScore = clamp(p.RiskScore+20, 0, 100)
	default:
		p.RiskScore = clamp(p.RiskScore+10, 0, 100)
	}
	p.RiskLevel = riskFromScore(p.RiskScore)
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *Registry) Get(userID string) (Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.profiles[userID]
	if !ok {
		return Profile{}, ErrNotFound
	}
	return *p, nil
}

func (r *Registry) FlagsFor(userID string) []AMLFlag {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []AMLFlag
	for _, f := range r.flags {
		if f.UserID == userID {
			out = append(out, f)
		}
	}
	return out
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func riskFromScore(score int) RiskLevel {
	switch {
	case score >= 60:
		return RiskHigh
	case score >= 30:
		return RiskMedium
	default:
		return RiskLow
	}
}
