package request

import "github.com/google/uuid"

// ── Spark ──────────────────────────────────────────────────────────────────

type CreateSparkClusterRequest struct {
	ClusterID    uuid.UUID `json:"cluster_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Version      string    `json:"version" validate:"required"`
	Namespace    string    `json:"namespace" validate:"required"`
	WorkerCount  int       `json:"worker_count"`
	WorkerCPU    string    `json:"worker_cpu"`
	WorkerMemory string    `json:"worker_memory"`
}

type UpdateSparkClusterRequest struct {
	Name        *string `json:"name"`
	Status      *string `json:"status"`
	WorkerCount *int    `json:"worker_count"`
}

// ── Flink ──────────────────────────────────────────────────────────────────

type CreateFlinkClusterRequest struct {
	ClusterID        uuid.UUID `json:"cluster_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	Version          string    `json:"version" validate:"required"`
	Namespace        string    `json:"namespace" validate:"required"`
	TaskManagerCount int       `json:"taskmanager_count"`
	SlotsPerTM       int       `json:"slots_per_tm"`
}

type UpdateFlinkClusterRequest struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
}

// ── Hive ───────────────────────────────────────────────────────────────────

type CreateHiveInstanceRequest struct {
	ClusterID      uuid.UUID `json:"cluster_id" validate:"required"`
	Name           string    `json:"name" validate:"required"`
	Version        string    `json:"version" validate:"required"`
	Namespace      string    `json:"namespace" validate:"required"`
	MetastoreURL   *string   `json:"metastore_url"`
	HiveServer2URL *string   `json:"hiveserver2_url"`
}

// ── HDFS ───────────────────────────────────────────────────────────────────

type CreateHDFSClusterRequest struct {
	ClusterID        uuid.UUID `json:"cluster_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	Version          string    `json:"version" validate:"required"`
	Namespace        string    `json:"namespace" validate:"required"`
	NamenodeCount    int       `json:"namenode_count"`
	DatanodeCount    int       `json:"datanode_count"`
	JournalnodeCount int       `json:"journalnode_count"`
}

// ── NiFi ───────────────────────────────────────────────────────────────────

type CreateNiFiInstanceRequest struct {
	ClusterID uuid.UUID `json:"cluster_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Version   string    `json:"version" validate:"required"`
	Namespace string    `json:"namespace" validate:"required"`
	NodeCount int       `json:"node_count"`
}
