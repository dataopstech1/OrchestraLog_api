package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type KafkaRepository struct {
	db *sqlx.DB
}

func NewKafkaRepository(db *sqlx.DB) *KafkaRepository {
	return &KafkaRepository{db: db}
}

// ── Kafka Clusters ─────────────────────────────────────────────────────────

func (r *KafkaRepository) ListClusters(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.KafkaCluster, int, error) {
	var clusters []model.KafkaCluster
	var total int

	query := `SELECT * FROM kafka_clusters WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM kafka_clusters WHERE 1=1`
	args := []any{}
	i := 1

	if clusterID != nil {
		query += ` AND cluster_id = $` + itoa(i)
		countQuery += ` AND cluster_id = $` + itoa(i)
		args = append(args, *clusterID)
		i++
	}

	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY created_at DESC LIMIT $` + itoa(i) + ` OFFSET $` + itoa(i+1)
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &clusters, query, args...)
	return clusters, total, err
}

func (r *KafkaRepository) FindClusterByID(ctx context.Context, id uuid.UUID) (*model.KafkaCluster, error) {
	var kc model.KafkaCluster
	err := r.db.GetContext(ctx, &kc, `SELECT * FROM kafka_clusters WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &kc, err
}

func (r *KafkaRepository) CreateCluster(ctx context.Context, kc *model.KafkaCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO kafka_clusters (
			id, cluster_id, name, version, environment, namespace, status,
			broker_count, broker_cpu, broker_memory, broker_storage, broker_storage_class, broker_jvm_heap,
			zk_enabled, zk_external_conn, zk_replicas, zk_cpu, zk_memory, zk_storage,
			auth_enabled, auth_mechanism, authz_enabled, authz_type,
			tls_enabled, tls_version, service_type, ingress_enabled, ingress_class,
			schema_registry_enabled, schema_registry_replicas,
			kafka_connect_enabled, kafka_connect_replicas,
			ksql_enabled, ksql_replicas,
			monitoring_enabled, prometheus_enabled, grafana_enabled, jmx_exporter_enabled,
			log_retention_hours, default_replication_factor, min_insync_replicas, created_by
		) VALUES (
			:id, :cluster_id, :name, :version, :environment, :namespace, :status,
			:broker_count, :broker_cpu, :broker_memory, :broker_storage, :broker_storage_class, :broker_jvm_heap,
			:zk_enabled, :zk_external_conn, :zk_replicas, :zk_cpu, :zk_memory, :zk_storage,
			:auth_enabled, :auth_mechanism, :authz_enabled, :authz_type,
			:tls_enabled, :tls_version, :service_type, :ingress_enabled, :ingress_class,
			:schema_registry_enabled, :schema_registry_replicas,
			:kafka_connect_enabled, :kafka_connect_replicas,
			:ksql_enabled, :ksql_replicas,
			:monitoring_enabled, :prometheus_enabled, :grafana_enabled, :jmx_exporter_enabled,
			:log_retention_hours, :default_replication_factor, :min_insync_replicas, :created_by
		)
	`, kc)
	return err
}

func (r *KafkaRepository) UpdateCluster(ctx context.Context, kc *model.KafkaCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE kafka_clusters SET
			name = :name, status = :status, version = :version, updated_at = NOW()
		WHERE id = :id
	`, kc)
	return err
}

func (r *KafkaRepository) DeleteCluster(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM kafka_clusters WHERE id = $1`, id)
	return err
}

// ── Topics ─────────────────────────────────────────────────────────────────

func (r *KafkaRepository) ListTopics(ctx context.Context, kafkaClusterID uuid.UUID, limit, offset int) ([]model.KafkaTopic, int, error) {
	var topics []model.KafkaTopic
	var total int

	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM kafka_topics WHERE kafka_cluster_id = $1`, kafkaClusterID); err != nil {
		return nil, 0, err
	}
	err := r.db.SelectContext(ctx, &topics, `
		SELECT * FROM kafka_topics WHERE kafka_cluster_id = $1
		ORDER BY name LIMIT $2 OFFSET $3
	`, kafkaClusterID, limit, offset)
	return topics, total, err
}

func (r *KafkaRepository) FindTopicByID(ctx context.Context, id uuid.UUID) (*model.KafkaTopic, error) {
	var t model.KafkaTopic
	err := r.db.GetContext(ctx, &t, `SELECT * FROM kafka_topics WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &t, err
}

func (r *KafkaRepository) CreateTopic(ctx context.Context, t *model.KafkaTopic) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO kafka_topics (id, kafka_cluster_id, name, partitions, replication_factor, retention_ms, status)
		VALUES (:id, :kafka_cluster_id, :name, :partitions, :replication_factor, :retention_ms, :status)
	`, t)
	return err
}

func (r *KafkaRepository) UpdateTopic(ctx context.Context, t *model.KafkaTopic) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE kafka_topics SET
			partitions = :partitions, retention_ms = :retention_ms, updated_at = NOW()
		WHERE id = :id
	`, t)
	return err
}

func (r *KafkaRepository) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM kafka_topics WHERE id = $1`, id)
	return err
}

// ── Consumer Groups ────────────────────────────────────────────────────────

func (r *KafkaRepository) ListConsumerGroups(ctx context.Context, kafkaClusterID uuid.UUID) ([]model.KafkaConsumerGroup, error) {
	var groups []model.KafkaConsumerGroup
	err := r.db.SelectContext(ctx, &groups, `
		SELECT * FROM kafka_consumer_groups WHERE kafka_cluster_id = $1 ORDER BY name
	`, kafkaClusterID)
	return groups, err
}

func (r *KafkaRepository) FindConsumerGroupByID(ctx context.Context, id uuid.UUID) (*model.KafkaConsumerGroup, error) {
	var g model.KafkaConsumerGroup
	err := r.db.GetContext(ctx, &g, `SELECT * FROM kafka_consumer_groups WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &g, err
}

// ── Resources ──────────────────────────────────────────────────────────────

func (r *KafkaRepository) ListResources(ctx context.Context, kafkaClusterID uuid.UUID) ([]model.KafkaResource, error) {
	var resources []model.KafkaResource
	err := r.db.SelectContext(ctx, &resources, `
		SELECT * FROM kafka_resources WHERE kafka_cluster_id = $1 ORDER BY type, name
	`, kafkaClusterID)
	return resources, err
}

func (r *KafkaRepository) FindResourceByID(ctx context.Context, id uuid.UUID) (*model.KafkaResource, error) {
	var res model.KafkaResource
	err := r.db.GetContext(ctx, &res, `SELECT * FROM kafka_resources WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &res, err
}

func (r *KafkaRepository) CreateResource(ctx context.Context, res *model.KafkaResource) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO kafka_resources (id, kafka_cluster_id, name, type, namespace, configuration, status)
		VALUES (:id, :kafka_cluster_id, :name, :type, :namespace, :configuration, :status)
	`, res)
	return err
}

func (r *KafkaRepository) UpdateResource(ctx context.Context, res *model.KafkaResource) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE kafka_resources SET status = :status, configuration = :configuration, updated_at = NOW()
		WHERE id = :id
	`, res)
	return err
}

func (r *KafkaRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM kafka_resources WHERE id = $1`, id)
	return err
}

// ── Stats ──────────────────────────────────────────────────────────────────

func (r *KafkaRepository) CountTopics(ctx context.Context, kafkaClusterID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM kafka_topics WHERE kafka_cluster_id = $1`, kafkaClusterID)
	return count, err
}

func (r *KafkaRepository) CountConsumerGroups(ctx context.Context, kafkaClusterID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM kafka_consumer_groups WHERE kafka_cluster_id = $1`, kafkaClusterID)
	return count, err
}

func (r *KafkaRepository) TotalConsumerLag(ctx context.Context, kafkaClusterID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.GetContext(ctx, &total, `
		SELECT COALESCE(SUM(total_lag), 0) FROM kafka_consumer_groups WHERE kafka_cluster_id = $1
	`, kafkaClusterID)
	return total, err
}
