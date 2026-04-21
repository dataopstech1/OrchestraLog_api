package model

import (
	"time"

	"github.com/google/uuid"
)

// ── Spark ──────────────────────────────────────────────────────────────────

type SparkCluster struct {
	ID           uuid.UUID `db:"id" json:"id"`
	ClusterID    uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name         string    `db:"name" json:"name"`
	Version      string    `db:"version" json:"version"`
	Namespace    string    `db:"namespace" json:"namespace"`
	Status       string    `db:"status" json:"status"`
	MasterURL    *string   `db:"master_url" json:"master_url"`
	WorkerCount  int       `db:"worker_count" json:"worker_count"`
	WorkerCPU    string    `db:"worker_cpu" json:"worker_cpu"`
	WorkerMemory string    `db:"worker_memory" json:"worker_memory"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type SparkApplication struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	SparkClusterID  uuid.UUID  `db:"spark_cluster_id" json:"spark_cluster_id"`
	Name            string     `db:"name" json:"name"`
	Status          string     `db:"status" json:"status"`
	Duration        *string    `db:"duration" json:"duration"`
	StartedAt       *time.Time `db:"started_at" json:"started_at"`
	CompletedAt     *time.Time `db:"completed_at" json:"completed_at"`
	StagesTotal     int        `db:"stages_total" json:"stages_total"`
	StagesCompleted int        `db:"stages_completed" json:"stages_completed"`
	ExecutorCount   int        `db:"executor_count" json:"executor_count"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

// ── Flink ──────────────────────────────────────────────────────────────────

type FlinkCluster struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ClusterID        uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name             string    `db:"name" json:"name"`
	Version          string    `db:"version" json:"version"`
	Namespace        string    `db:"namespace" json:"namespace"`
	Status           string    `db:"status" json:"status"`
	JobManagerURL    *string   `db:"jobmanager_url" json:"jobmanager_url"`
	TaskManagerCount int       `db:"taskmanager_count" json:"taskmanager_count"`
	SlotsPerTM       int       `db:"slots_per_tm" json:"slots_per_tm"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type FlinkJob struct {
	ID                   uuid.UUID  `db:"id" json:"id"`
	FlinkClusterID       uuid.UUID  `db:"flink_cluster_id" json:"flink_cluster_id"`
	Name                 string     `db:"name" json:"name"`
	Status               string     `db:"status" json:"status"`
	StartTime            *time.Time `db:"start_time" json:"start_time"`
	Duration             *string    `db:"duration" json:"duration"`
	Parallelism          int        `db:"parallelism" json:"parallelism"`
	CheckpointsCompleted int        `db:"checkpoints_completed" json:"checkpoints_completed"`
	BytesIn              int64      `db:"bytes_in" json:"bytes_in"`
	BytesOut             int64      `db:"bytes_out" json:"bytes_out"`
	RecordsIn            int64      `db:"records_in" json:"records_in"`
	RecordsOut           int64      `db:"records_out" json:"records_out"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
}

// ── Hive ───────────────────────────────────────────────────────────────────

type HiveInstance struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ClusterID      uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name           string    `db:"name" json:"name"`
	Version        string    `db:"version" json:"version"`
	Namespace      string    `db:"namespace" json:"namespace"`
	Status         string    `db:"status" json:"status"`
	MetastoreURL   *string   `db:"metastore_url" json:"metastore_url"`
	HiveServer2URL *string   `db:"hiveserver2_url" json:"hiveserver2_url"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type HiveTable struct {
	ID             uuid.UUID `db:"id" json:"id"`
	HiveInstanceID uuid.UUID `db:"hive_instance_id" json:"hive_instance_id"`
	DatabaseName   string    `db:"database_name" json:"database_name"`
	TableName      string    `db:"table_name" json:"table_name"`
	TableType      *string   `db:"table_type" json:"table_type"`
	Format         *string   `db:"format" json:"format"`
	Partitions     int       `db:"partitions" json:"partitions"`
	TotalSize      *string   `db:"total_size" json:"total_size"`
	RowsCount      int64     `db:"rows_count" json:"rows_count"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// ── HDFS ───────────────────────────────────────────────────────────────────

type HDFSCluster struct {
	ID                     uuid.UUID `db:"id" json:"id"`
	ClusterID              uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name                   string    `db:"name" json:"name"`
	Version                string    `db:"version" json:"version"`
	Namespace              string    `db:"namespace" json:"namespace"`
	Status                 string    `db:"status" json:"status"`
	NamenodeCount          int       `db:"namenode_count" json:"namenode_count"`
	DatanodeCount          int       `db:"datanode_count" json:"datanode_count"`
	JournalnodeCount       int       `db:"journalnode_count" json:"journalnode_count"`
	TotalCapacity          *string   `db:"total_capacity" json:"total_capacity"`
	UsedCapacity           *string   `db:"used_capacity" json:"used_capacity"`
	RemainingCapacity      *string   `db:"remaining_capacity" json:"remaining_capacity"`
	TotalBlocks            int       `db:"total_blocks" json:"total_blocks"`
	UnderReplicatedBlocks  int       `db:"under_replicated_blocks" json:"under_replicated_blocks"`
	CreatedAt              time.Time `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at"`
}

// ── NiFi ───────────────────────────────────────────────────────────────────

type NiFiInstance struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClusterID uuid.UUID `db:"cluster_id" json:"cluster_id"`
	Name      string    `db:"name" json:"name"`
	Version   string    `db:"version" json:"version"`
	Namespace string    `db:"namespace" json:"namespace"`
	Status    string    `db:"status" json:"status"`
	UIURL     *string   `db:"ui_url" json:"ui_url"`
	NodeCount int       `db:"node_count" json:"node_count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type NiFiProcessGroup struct {
	ID               uuid.UUID `db:"id" json:"id"`
	NiFiInstanceID   uuid.UUID `db:"nifi_instance_id" json:"nifi_instance_id"`
	Name             string    `db:"name" json:"name"`
	Status           string    `db:"status" json:"status"`
	ProcessorsTotal  int       `db:"processors_total" json:"processors_total"`
	ProcessorsRunning int      `db:"processors_running" json:"processors_running"`
	ProcessorsStopped int      `db:"processors_stopped" json:"processors_stopped"`
	InputBytesPerSec  *string  `db:"input_bytes_per_sec" json:"input_bytes_per_sec"`
	OutputBytesPerSec *string  `db:"output_bytes_per_sec" json:"output_bytes_per_sec"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
