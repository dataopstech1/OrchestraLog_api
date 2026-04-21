package model

import (
	"time"

	"github.com/google/uuid"
)

type Cluster struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Name         string     `db:"name" json:"name"`
	Status       string     `db:"status" json:"status"`
	Region       string     `db:"region" json:"region"`
	K8sVersion   string     `db:"k8s_version" json:"k8s_version"`
	APIServerURL *string    `db:"api_server_url" json:"api_server_url"`
	Kubeconfig   *string    `db:"kubeconfig" json:"-"`
	NodesTotal   int        `db:"nodes_total" json:"nodes_total"`
	NodesReady   int        `db:"nodes_ready" json:"nodes_ready"`
	CreatedBy    *uuid.UUID `db:"created_by" json:"created_by"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type Namespace struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ClusterID   uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name        string    `db:"name" json:"name"`
	Status      string    `db:"status" json:"status"`
	CPUQuota    *string   `db:"cpu_quota" json:"cpu_quota"`
	MemoryQuota *string   `db:"memory_quota" json:"memory_quota"`
	PodsQuota   *string   `db:"pods_quota" json:"pods_quota"`
	CPUUsed     *string   `db:"cpu_used" json:"cpu_used"`
	MemoryUsed  *string   `db:"memory_used" json:"memory_used"`
	PodsUsed    int       `db:"pods_used" json:"pods_used"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Resource struct {
	ID             uuid.UUID `db:"id" json:"id"`
	NamespaceID    uuid.UUID `db:"namespace_id" json:"namespace_id"`
	Name           string    `db:"name" json:"name"`
	Type           string    `db:"type" json:"type"`
	Status         string    `db:"status" json:"status"`
	ReplicasReady  *int      `db:"replicas_ready" json:"replicas_ready"`
	ReplicasTotal  *int      `db:"replicas_total" json:"replicas_total"`
	CPUUsage       *string   `db:"cpu_usage" json:"cpu_usage"`
	MemoryUsage    *string   `db:"memory_usage" json:"memory_usage"`
	Age            *string   `db:"age" json:"age"`
	Labels         []byte    `db:"labels" json:"labels,omitempty"`
	Annotations    []byte    `db:"annotations" json:"annotations,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
