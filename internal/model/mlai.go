package model

import (
	"time"

	"github.com/google/uuid"
)

// ── MLFlow ─────────────────────────────────────────────────────────────────

type MLFlowInstance struct {
	ID            uuid.UUID `db:"id" json:"id"`
	ClusterID     uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name          string    `db:"name" json:"name"`
	Version       string    `db:"version" json:"version"`
	Namespace     string    `db:"namespace" json:"namespace"`
	Status        string    `db:"status" json:"status"`
	TrackingURL   *string   `db:"tracking_url" json:"tracking_url"`
	ArtifactStore *string   `db:"artifact_store" json:"artifact_store"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type MLFlowExperiment struct {
	ID               uuid.UUID `db:"id" json:"id"`
	MLFlowInstanceID uuid.UUID `db:"mlflow_instance_id" json:"mlflow_instance_id"`
	Name             string    `db:"name" json:"name"`
	Status           string    `db:"status" json:"status"`
	RunsTotal        int       `db:"runs_total" json:"runs_total"`
	BestMetricName   *string   `db:"best_metric_name" json:"best_metric_name"`
	BestMetricValue  *float64  `db:"best_metric_value" json:"best_metric_value"`
	CreatedBy        *string   `db:"created_by" json:"created_by"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type MLFlowModel struct {
	ID               uuid.UUID `db:"id" json:"id"`
	MLFlowInstanceID uuid.UUID `db:"mlflow_instance_id" json:"mlflow_instance_id"`
	Name             string    `db:"name" json:"name"`
	LatestVersion    int       `db:"latest_version" json:"latest_version"`
	Stage            *string   `db:"stage" json:"stage"`
	Description      *string   `db:"description" json:"description"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// ── Feast ──────────────────────────────────────────────────────────────────

type FeastInstance struct {
	ID           uuid.UUID `db:"id" json:"id"`
	ClusterID    uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name         string    `db:"name" json:"name"`
	Version      string    `db:"version" json:"version"`
	Namespace    string    `db:"namespace" json:"namespace"`
	Status       string    `db:"status" json:"status"`
	OnlineStore  *string   `db:"online_store" json:"online_store"`
	OfflineStore *string   `db:"offline_store" json:"offline_store"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type FeastEntity struct {
	ID               uuid.UUID `db:"id" json:"id"`
	FeastInstanceID  uuid.UUID `db:"feast_instance_id" json:"feast_instance_id"`
	Name             string    `db:"name" json:"name"`
	ValueType        string    `db:"value_type" json:"value_type"`
	Description      *string   `db:"description" json:"description"`
	JoinKeys         []byte    `db:"join_keys" json:"join_keys,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type FeastFeatureView struct {
	ID              uuid.UUID `db:"id" json:"id"`
	FeastInstanceID uuid.UUID `db:"feast_instance_id" json:"feast_instance_id"`
	Name            string    `db:"name" json:"name"`
	Entities        []byte    `db:"entities" json:"entities,omitempty"`
	Features        []byte    `db:"features" json:"features,omitempty"`
	TTL             *string   `db:"ttl" json:"ttl"`
	Source          *string   `db:"source" json:"source"`
	Online          bool      `db:"online" json:"online"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// ── JupyterHub ─────────────────────────────────────────────────────────────

type JupyterHubInstance struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClusterID uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name      string    `db:"name" json:"name"`
	Version   string    `db:"version" json:"version"`
	Namespace string    `db:"namespace" json:"namespace"`
	Status    string    `db:"status" json:"status"`
	HubURL    *string   `db:"hub_url" json:"hub_url"`
	MaxUsers  int       `db:"max_users" json:"max_users"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type JupyterHubNotebook struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	JupyterHubInstanceID  uuid.UUID  `db:"jupyterhub_instance_id" json:"jupyterhub_instance_id"`
	Name                  string     `db:"name" json:"name"`
	Owner                 string     `db:"owner" json:"owner"`
	Status                string     `db:"status" json:"status"`
	Image                 *string    `db:"image" json:"image"`
	CPULimit              *string    `db:"cpu_limit" json:"cpu_limit"`
	MemoryLimit           *string    `db:"memory_limit" json:"memory_limit"`
	GPULimit              int        `db:"gpu_limit" json:"gpu_limit"`
	LastActivity          *time.Time `db:"last_activity" json:"last_activity"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updated_at"`
}

// ── LLM ───────────────────────────────────────────────────────────────────

type LLMDeployment struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ClusterID      uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name           string    `db:"name" json:"name"`
	ModelName      string    `db:"model_name" json:"model_name"`
	ModelVersion   *string   `db:"model_version" json:"model_version"`
	Namespace      string    `db:"namespace" json:"namespace"`
	Status         string    `db:"status" json:"status"`
	EndpointURL    *string   `db:"endpoint_url" json:"endpoint_url"`
	Replicas       int       `db:"replicas" json:"replicas"`
	GPUCount       int       `db:"gpu_count" json:"gpu_count"`
	GPUType        *string   `db:"gpu_type" json:"gpu_type"`
	MaxTokens      *int      `db:"max_tokens" json:"max_tokens"`
	ContextWindow  *int      `db:"context_window" json:"context_window"`
	RequestsPerMin int       `db:"requests_per_min" json:"requests_per_min"`
	AvgLatencyMs   float64   `db:"avg_latency_ms" json:"avg_latency_ms"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
