package request

import "github.com/google/uuid"

// ── MLFlow ─────────────────────────────────────────────────────────────────

type CreateMLFlowInstanceRequest struct {
	ClusterID     uuid.UUID `json:"cluster_id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Version       string    `json:"version" validate:"required"`
	Namespace     string    `json:"namespace" validate:"required"`
	TrackingURL   *string   `json:"tracking_url"`
	ArtifactStore *string   `json:"artifact_store"`
}

type UpdateMLFlowModelStageRequest struct {
	Stage string `json:"stage" validate:"required,oneof=None Staging Production Archived"`
}

// ── Feast ──────────────────────────────────────────────────────────────────

type CreateFeastInstanceRequest struct {
	ClusterID    uuid.UUID `json:"cluster_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Version      string    `json:"version" validate:"required"`
	Namespace    string    `json:"namespace" validate:"required"`
	OnlineStore  *string   `json:"online_store"`
	OfflineStore *string   `json:"offline_store"`
}

// ── JupyterHub ─────────────────────────────────────────────────────────────

type CreateJupyterHubInstanceRequest struct {
	ClusterID uuid.UUID `json:"cluster_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Version   string    `json:"version" validate:"required"`
	Namespace string    `json:"namespace" validate:"required"`
	HubURL    *string   `json:"hub_url"`
	MaxUsers  int       `json:"max_users"`
}

type CreateNotebookRequest struct {
	Name        string  `json:"name" validate:"required"`
	Owner       string  `json:"owner" validate:"required"`
	Image       *string `json:"image"`
	CPULimit    *string `json:"cpu_limit"`
	MemoryLimit *string `json:"memory_limit"`
	GPULimit    int     `json:"gpu_limit"`
}

type UpdateNotebookRequest struct {
	Image       *string `json:"image"`
	CPULimit    *string `json:"cpu_limit"`
	MemoryLimit *string `json:"memory_limit"`
	GPULimit    *int    `json:"gpu_limit"`
}

// ── LLM ───────────────────────────────────────────────────────────────────

type CreateLLMDeploymentRequest struct {
	ClusterID     uuid.UUID `json:"cluster_id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	ModelName     string    `json:"model_name" validate:"required"`
	ModelVersion  *string   `json:"model_version"`
	Namespace     string    `json:"namespace" validate:"required"`
	Replicas      int       `json:"replicas"`
	GPUCount      int       `json:"gpu_count"`
	GPUType       *string   `json:"gpu_type"`
	MaxTokens     *int      `json:"max_tokens"`
	ContextWindow *int      `json:"context_window"`
}

type UpdateLLMDeploymentRequest struct {
	Name    *string `json:"name"`
	Status  *string `json:"status"`
}

type ScaleLLMDeploymentRequest struct {
	Replicas int `json:"replicas" validate:"required,min=1"`
}
