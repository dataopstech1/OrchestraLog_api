package request

import "github.com/google/uuid"

type KafkaBrokerConfig struct {
	Count        int    `json:"count" validate:"required,min=1"`
	CPU          string `json:"cpu"`
	Memory       string `json:"memory"`
	Storage      string `json:"storage"`
	StorageClass string `json:"storage_class"`
	JVMHeap      string `json:"jvm_heap"`
}

type KafkaZookeeperConfig struct {
	Enabled      bool   `json:"enabled"`
	ExternalConn string `json:"external_conn"`
	Replicas     int    `json:"replicas"`
	CPU          string `json:"cpu"`
	Memory       string `json:"memory"`
	Storage      string `json:"storage"`
}

type KafkaSecurityConfig struct {
	AuthEnabled   bool   `json:"auth_enabled"`
	AuthMechanism string `json:"auth_mechanism" validate:"omitempty,oneof=PLAIN SCRAM-SHA-256 SCRAM-SHA-512"`
	AuthzEnabled  bool   `json:"authz_enabled"`
	AuthzType     string `json:"authz_type" validate:"omitempty,oneof=simple opa"`
	TLSEnabled    bool   `json:"tls_enabled"`
	TLSVersion    string `json:"tls_version"`
}

type KafkaServicesConfig struct {
	SchemaRegistryEnabled  bool `json:"schema_registry_enabled"`
	SchemaRegistryReplicas int  `json:"schema_registry_replicas"`
	KafkaConnectEnabled    bool `json:"kafka_connect_enabled"`
	KafkaConnectReplicas   int  `json:"kafka_connect_replicas"`
	KSQLEnabled            bool `json:"ksql_enabled"`
	KSQLReplicas           int  `json:"ksql_replicas"`
	MonitoringEnabled      bool `json:"monitoring_enabled"`
	PrometheusEnabled      bool `json:"prometheus_enabled"`
	GrafanaEnabled         bool `json:"grafana_enabled"`
}

type CreateKafkaClusterRequest struct {
	ClusterID   uuid.UUID            `json:"cluster_id" validate:"required"`
	Name        string               `json:"name" validate:"required"`
	Version     string               `json:"version" validate:"required"`
	Environment string               `json:"environment" validate:"required,oneof=development staging production"`
	Namespace   string               `json:"namespace" validate:"required"`
	Brokers     KafkaBrokerConfig    `json:"brokers" validate:"required"`
	Zookeeper   KafkaZookeeperConfig `json:"zookeeper"`
	Security    KafkaSecurityConfig  `json:"security"`
	Services    KafkaServicesConfig  `json:"services"`
	LogRetentionHours        int     `json:"log_retention_hours"`
	DefaultReplicationFactor int     `json:"default_replication_factor"`
	MinInsyncReplicas        int     `json:"min_insync_replicas"`
}

type UpdateKafkaClusterRequest struct {
	Name      *string `json:"name"`
	Status    *string `json:"status" validate:"omitempty,oneof=pending creating running failed updating deleting"`
	Version   *string `json:"version"`
}

type CreateKafkaTopicRequest struct {
	Name              string `json:"name" validate:"required"`
	Partitions        int    `json:"partitions" validate:"required,min=1"`
	ReplicationFactor int    `json:"replication_factor" validate:"required,min=1"`
	RetentionMs       int64  `json:"retention_ms"`
}

type UpdateKafkaTopicRequest struct {
	Partitions  *int   `json:"partitions" validate:"omitempty,min=1"`
	RetentionMs *int64 `json:"retention_ms"`
}

type CreateKafkaResourceRequest struct {
	Name          string `json:"name" validate:"required"`
	Type          string `json:"type" validate:"required,oneof=topic user acl connector custom"`
	Namespace     string `json:"namespace"`
	Configuration []byte `json:"configuration"`
}

type UpdateKafkaResourceRequest struct {
	Status        *string `json:"status" validate:"omitempty,oneof=active pending failed deleting"`
	Configuration []byte  `json:"configuration"`
}
