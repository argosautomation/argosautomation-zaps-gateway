package db

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents an organization using the platform
type Tenant struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	MonthlyQuota   int       `json:"monthly_quota" db:"monthly_quota"`
	CurrentUsage   int       `json:"current_usage" db:"current_usage"`
	QuotaResetAt   time.Time `json:"quota_reset_at" db:"quota_reset_at"`
	OverageAllowed bool      `json:"overage_allowed" db:"overage_allowed"`
	Metadata       JSONBMap  `json:"metadata,omitempty" db:"metadata"`
}



// APIKey represents the metadata for a gateway API key
type APIKey struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	TenantID     uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	KeyHash      string     `json:"-" db:"key_hash"`            // Never expose
	KeyPrefix    string     `json:"key_prefix" db:"key_prefix"` // For display (e.g., "gk_abc123...")
	Name         string     `json:"name" db:"name"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	LastUsed     *time.Time `json:"last_used,omitempty" db:"last_used"`
	RequestCount int64      `json:"request_count" db:"request_count"`
	Enabled      bool       `json:"enabled" db:"enabled"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	Notes        *string    `json:"notes,omitempty" db:"notes"`
}

// UsageLog represents aggregated hourly usage statistics
type UsageLog struct {
	ID                   int64      `json:"id" db:"id"`
	TenantID             uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	APIKeyID             *uuid.UUID `json:"api_key_id,omitempty" db:"api_key_id"`
	HourBucket           time.Time  `json:"hour_bucket" db:"hour_bucket"`
	RequestCount         int        `json:"request_count" db:"request_count"`
	ErrorCount           int        `json:"error_count" db:"error_count"`
	AvgLatencyMs         *int       `json:"avg_latency_ms,omitempty" db:"avg_latency_ms"`
	TotalTokensProcessed int64      `json:"total_tokens_processed" db:"total_tokens_processed"`
	PIIEvents            JSONBMap   `json:"pii_events,omitempty" db:"pii_events"`
}



// AuditLog represents a security/audit event
type AuditLog struct {
	ID        int64      `json:"id" db:"id"`
	TenantID  *uuid.UUID `json:"tenant_id,omitempty" db:"tenant_id"`
	EventType string     `json:"event_type" db:"event_type"`
	EventData JSONBMap   `json:"event_data,omitempty" db:"event_data"`
	IPAddress *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string    `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// ProviderKey represents an encrypted API key for an upstream LLM provider
type ProviderKey struct {
	ID           uuid.UUID `json:"id" db:"id"`
	TenantID     uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Provider     string    `json:"provider" db:"provider"`
	EncryptedKey string    `json:"-" db:"encrypted_key"` // Never expose
	Enabled      bool      `json:"enabled" db:"enabled"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// JSONBMap is a helper type for PostgreSQL JSONB columns
type JSONBMap map[string]interface{}
