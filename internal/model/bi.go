package model

import (
	"time"

	"github.com/google/uuid"
)

// ── Superset ───────────────────────────────────────────────────────────────

type SupersetInstance struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClusterID uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name      string    `db:"name" json:"name"`
	Version   string    `db:"version" json:"version"`
	Namespace string    `db:"namespace" json:"namespace"`
	Status    string    `db:"status" json:"status"`
	URL       *string   `db:"url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type SupersetDashboard struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	SupersetInstanceID uuid.UUID `db:"superset_instance_id" json:"superset_instance_id"`
	Name               string    `db:"name" json:"name"`
	Slug               *string   `db:"slug" json:"slug"`
	Status             string    `db:"status" json:"status"`
	ChartsCount        int       `db:"charts_count" json:"charts_count"`
	Owners             []byte    `db:"owners" json:"owners,omitempty"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// ── Metabase ───────────────────────────────────────────────────────────────

type MetabaseInstance struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClusterID uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name      string    `db:"name" json:"name"`
	Version   string    `db:"version" json:"version"`
	Namespace string    `db:"namespace" json:"namespace"`
	Status    string    `db:"status" json:"status"`
	URL       *string   `db:"url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type MetabaseDashboard struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	MetabaseInstanceID uuid.UUID `db:"metabase_instance_id" json:"metabase_instance_id"`
	Name               string    `db:"name" json:"name"`
	Collection         *string   `db:"collection" json:"collection"`
	CardsCount         int       `db:"cards_count" json:"cards_count"`
	ViewsCount         int       `db:"views_count" json:"views_count"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// ── N8N ───────────────────────────────────────────────────────────────────

type N8NInstance struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClusterID uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name      string    `db:"name" json:"name"`
	Version   string    `db:"version" json:"version"`
	Namespace string    `db:"namespace" json:"namespace"`
	Status    string    `db:"status" json:"status"`
	URL       *string   `db:"url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type N8NWorkflow struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	N8NInstanceID    uuid.UUID  `db:"n8n_instance_id" json:"n8n_instance_id"`
	Name             string     `db:"name" json:"name"`
	Status           string     `db:"status" json:"status"`
	NodesCount       int        `db:"nodes_count" json:"nodes_count"`
	ConnectionsCount int        `db:"connections_count" json:"connections_count"`
	LastExecutionAt  *time.Time `db:"last_execution_at" json:"last_execution_at"`
	TotalExecutions  int        `db:"total_executions" json:"total_executions"`
	SuccessRate      float64    `db:"success_rate" json:"success_rate"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}
