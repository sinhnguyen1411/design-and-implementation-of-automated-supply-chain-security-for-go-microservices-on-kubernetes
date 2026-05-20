package apikey

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Scope string

const (
	ScopeRead       Scope = "read"
	ScopeTrade      Scope = "trade"
	ScopeAdmin      Scope = "admin"
	ScopeMarketData Scope = "market_data"
)

type Tier string

const (
	TierFree Tier = "free"
	TierPro  Tier = "pro"
	TierAlgo Tier = "algo"
)

type APIKey struct {
	ID        string     `json:"id"`
	OwnerID   string     `json:"owner_id"`
	Name      string     `json:"name"`
	Key       string     `json:"key,omitempty"` // exposed only on Issue/Rotate
	Scopes    []Scope    `json:"scopes"`
	Tier      Tier       `json:"tier"`
	RateLimit int        `json:"rate_limit"` // req/min
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type Store struct {
	mu    sync.RWMutex
	keys  map[string]*APIKey // id → key
	byRaw map[string]string  // raw key → id
}

func NewStore() *Store {
	return &Store{
		keys:  make(map[string]*APIKey),
		byRaw: make(map[string]string),
	}
}

var (
	ErrKeyNotFound       = errors.New("api key not found")
	ErrKeyRevoked        = errors.New("api key has been revoked")
	ErrKeyExpired        = errors.New("api key has expired")
	ErrScopeInsufficient = errors.New("insufficient scope")
)

func generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(b), nil
}

func (s *Store) Issue(ownerID, name string, scopes []Scope, tier Tier, rateLimit int, ttl *time.Duration) (APIKey, error) {
	raw, err := generate()
	if err != nil {
		return APIKey{}, err
	}
	var expiresAt *time.Time
	if ttl != nil {
		t := time.Now().UTC().Add(*ttl)
		expiresAt = &t
	}
	ak := &APIKey{
		ID:        uuid.NewString(),
		OwnerID:   ownerID,
		Name:      name,
		Key:       raw,
		Scopes:    scopes,
		Tier:      tier,
		RateLimit: rateLimit,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
	}
	s.mu.Lock()
	s.keys[ak.ID] = ak
	s.byRaw[raw] = ak.ID
	s.mu.Unlock()
	return *ak, nil
}

func (s *Store) Validate(rawKey string, required Scope) (APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.byRaw[rawKey]
	if !ok {
		return APIKey{}, ErrKeyNotFound
	}
	ak := s.keys[id]
	if !ak.Active {
		return APIKey{}, ErrKeyRevoked
	}
	if ak.ExpiresAt != nil && time.Now().After(*ak.ExpiresAt) {
		return APIKey{}, ErrKeyExpired
	}
	hasScope := false
	for _, sc := range ak.Scopes {
		if sc == required {
			hasScope = true
			break
		}
	}
	if !hasScope {
		return APIKey{}, ErrScopeInsufficient
	}
	now := time.Now().UTC()
	ak.LastUsed = &now
	out := *ak
	out.Key = ""
	return out, nil
}

func (s *Store) Revoke(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ak, ok := s.keys[id]
	if !ok {
		return ErrKeyNotFound
	}
	ak.Active = false
	return nil
}

func (s *Store) Rotate(id string) (APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.keys[id]
	if !ok {
		return APIKey{}, ErrKeyNotFound
	}
	delete(s.byRaw, old.Key)
	old.Active = false

	raw, err := generate()
	if err != nil {
		return APIKey{}, err
	}
	newAK := &APIKey{
		ID:        uuid.NewString(),
		OwnerID:   old.OwnerID,
		Name:      old.Name,
		Key:       raw,
		Scopes:    old.Scopes,
		Tier:      old.Tier,
		RateLimit: old.RateLimit,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: old.ExpiresAt,
	}
	s.keys[newAK.ID] = newAK
	s.byRaw[raw] = newAK.ID
	return *newAK, nil
}

func (s *Store) ListByOwner(ownerID string) []APIKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []APIKey
	for _, ak := range s.keys {
		if ak.OwnerID == ownerID {
			safe := *ak
			safe.Key = ""
			result = append(result, safe)
		}
	}
	return result
}
