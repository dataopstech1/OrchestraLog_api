package model

import (
	"time"

	"github.com/google/uuid"
)

type KafkaCluster struct {
	ID                       uuid.UUID  `db:"id" json:"id"`
	ClusterID                uuid.UUID  `db:"cluster_id" json:"cluster_id"`
	Name                     string     `db:"name" json:"name"`
	Version                  string     `db:"version" json:"version"`
	Environment              string     `db:"environment" json:"environment"`
	Namespace                string     `db:"namespace" json:"namespace"`
	Status                   string     `db:"status" json:"status"`
	BrokerCount              int        `db:"broker_count" json:"broker_count"`
	BrokerCPU                string     `db:"broker_cpu" json:"broker_cpu"`
	BrokerMemory             string     `db:"broker_memory" json:"broker_memory"`
	BrokerStorage            string     `db:"broker_storage" json:"broker_storage"`
	BrokerStorageClass       string     `db:"broker_storage_class" json:"broker_storage_class"`
	BrokerJVMHeap            string     `db:"broker_jvm_heap" json:"broker_jvm_heap"`
	ZKEnabled                bool       `db:"zk_enabled" json:"zk_enabled"`
	ZKExternalConn           *string    `db:"zk_external_conn" json:"zk_external_conn"`
	ZKReplicas               int        `db:"zk_replicas" json:"zk_replicas"`
	ZKCPU                    string     `db:"zk_cpu" json:"zk_cpu"`
	ZKMemory                 string     `db:"zk_memory" json:"zk_memory"`
	ZKStorage                string     `db:"zk_storage" json:"zk_storage"`
	AuthEnabled              bool       `db:"auth_enabled" json:"auth_enabled"`
	AuthMechanism            *string    `db:"auth_mechanism" json:"auth_mechanism"`
	AuthzEnabled             bool       `db:"authz_enabled" json:"authz_enabled"`
	AuthzType                *string    `db:"authz_type" json:"authz_type"`
	TLSEnabled               bool       `db:"tls_enabled" json:"tls_enabled"`
	TLSVersion               string     `db:"tls_version" json:"tls_version"`
	ServiceType              string     `db:"service_type" json:"service_type"`
	IngressEnabled           bool       `db:"ingress_enabled" json:"ingress_enabled"`
	IngressClass             *string    `db:"ingress_class" json:"ingress_class"`
	IngressAnnotations       []byte     `db:"ingress_annotations" json:"ingress_annotations,omitempty"`
	SchemaRegistryEnabled    bool       `db:"schema_registry_enabled" json:"schema_registry_enabled"`
	SchemaRegistryReplicas   int        `db:"schema_registry_replicas" json:"schema_registry_replicas"`
	KafkaConnectEnabled      bool       `db:"kafka_connect_enabled" json:"kafka_connect_enabled"`
	KafkaConnectReplicas     int        `db:"kafka_connect_replicas" json:"kafka_connect_replicas"`
	KSQLEnabled              bool       `db:"ksql_enabled" json:"ksql_enabled"`
	KSQLReplicas             int        `db:"ksql_replicas" json:"ksql_replicas"`
	MonitoringEnabled        bool       `db:"monitoring_enabled" json:"monitoring_enabled"`
	PrometheusEnabled        bool       `db:"prometheus_enabled" json:"prometheus_enabled"`
	GrafanaEnabled           bool       `db:"grafana_enabled" json:"grafana_enabled"`
	JMXExporterEnabled       bool       `db:"jmx_exporter_enabled" json:"jmx_exporter_enabled"`
	ServerProperties         []byte     `db:"server_properties" json:"server_properties,omitempty"`
	JVMOptions               []byte     `db:"jvm_options" json:"jvm_options,omitempty"`
	LogRetentionHours        int        `db:"log_retention_hours" json:"log_retention_hours"`
	DefaultReplicationFactor int        `db:"default_replication_factor" json:"default_replication_factor"`
	MinInsyncReplicas        int        `db:"min_insync_replicas" json:"min_insync_replicas"`
	CreatedBy                *uuid.UUID `db:"created_by" json:"created_by"`
	CreatedAt                time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time  `db:"updated_at" json:"updated_at"`
}

type KafkaTopic struct {
	ID                uuid.UUID `db:"id" json:"id"`
	KafkaClusterID    uuid.UUID `db:"kafka_cluster_id" json:"kafka_cluster_id"`
	Name              string    `db:"name" json:"name"`
	Partitions        int       `db:"partitions" json:"partitions"`
	ReplicationFactor int       `db:"replication_factor" json:"replication_factor"`
	RetentionMs       int64     `db:"retention_ms" json:"retention_ms"`
	Status            string    `db:"status" json:"status"`
	Configuration     []byte    `db:"configuration" json:"configuration,omitempty"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type KafkaConsumerGroup struct {
	ID              uuid.UUID `db:"id" json:"id"`
	KafkaClusterID  uuid.UUID `db:"kafka_cluster_id" json:"kafka_cluster_id"`
	Name            string    `db:"name" json:"name"`
	State           *string   `db:"state" json:"state"`
	Members         int       `db:"members" json:"members"`
	TopicsSubscribed []byte   `db:"topics_subscribed" json:"topics_subscribed,omitempty"`
	TotalLag        int64     `db:"total_lag" json:"total_lag"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type KafkaResource struct {
	ID             uuid.UUID `db:"id" json:"id"`
	KafkaClusterID uuid.UUID `db:"kafka_cluster_id" json:"kafka_cluster_id"`
	Name           string    `db:"name" json:"name"`
	Type           string    `db:"type" json:"type"`
	Namespace      *string   `db:"namespace" json:"namespace"`
	Configuration  []byte    `db:"configuration" json:"configuration,omitempty"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
