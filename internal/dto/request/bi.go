package request

import "github.com/google/uuid"

type CreateSupersetInstanceRequest struct {
	ClusterID uuid.UUID `json:"cluster_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Version   string    `json:"version" validate:"required"`
	Namespace string    `json:"namespace" validate:"required"`
	URL       *string   `json:"url"`
}

type CreateMetabaseInstanceRequest struct {
	ClusterID uuid.UUID `json:"cluster_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Version   string    `json:"version" validate:"required"`
	Namespace string    `json:"namespace" validate:"required"`
	URL       *string   `json:"url"`
}

type CreateN8NInstanceRequest struct {
	ClusterID uuid.UUID `json:"cluster_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Version   string    `json:"version" validate:"required"`
	Namespace string    `json:"namespace" validate:"required"`
	URL       *string   `json:"url"`
}
