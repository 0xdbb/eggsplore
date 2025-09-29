package server

import (
	"encoding/json"
	"time"
)

type Accounts struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	ProfileUrl string    `json:"profile_url"`
	Status     string    `json:"status"`
	Role       string    `json:"role"`
	IsApproved bool      `json:"is_approved"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Eggs struct {
	InventoryID string    `json:"inventory_id"`
	Hatched     bool      `json:"hatched"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	CollectedAt time.Time `json:"collected_at"`
}

type GetEggsByPlayerResponse struct {
	InventoryID string    `json:"inventory_id"`
	Type        string    `json:"type"`
	Hatched     bool      `json:"hatched"`
	Message     string    `json:"message"`
	CollectedAt time.Time `json:"collected_at"`
}

type GetToolsByPlayerResponse struct {
	InventoryID string `json:"inventory_id"`
	Durability  int32  `json:"durability"`
	Equipped    bool   `json:"equipped"`
	Description string `json:"description"`
}

type Inventory struct {
	ID          string    `json:"id"`
	PlayerID    string    `json:"player_id"`
	ItemType    string    `json:"item_type"`
	Quantity    int32     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Players struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Coins     int64           `json:"coins"`
	Xp        int64           `json:"xp"`
	Level     int32           `json:"level"`
	Settings  json.RawMessage `json:"settings"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Session struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	RefreshToken string    `json:"refresh_token"`
	AccountAgent string    `json:"account_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Tools struct {
	InventoryID string `json:"inventory_id"`
	Durability  int32  `json:"durability"`
	Equipped    bool   `json:"equipped"`
}
