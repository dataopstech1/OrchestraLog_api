package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type BigDataRepository struct {
	db *sqlx.DB
}

func NewBigDataRepository(db *sqlx.DB) *BigDataRepository {
	return &BigDataRepository{db: db}
}

// ── Spark ──────────────────────────────────────────────────────────────────

func (r *BigDataRepository) ListSparkClusters(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.SparkCluster, int, error) {
	return listWithOptionalFilter[model.SparkCluster](ctx, r.db, "spark_clusters", "cluster_id", clusterID, limit, offset)
}

func (r *BigDataRepository) FindSparkCluster(ctx context.Context, id uuid.UUID) (*model.SparkCluster, error) {
	return findByID[model.SparkCluster](ctx, r.db, "spark_clusters", id)
}

func (r *BigDataRepository) CreateSparkCluster(ctx context.Context, sc *model.SparkCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO spark_clusters (id, cluster_id, name, version, namespace, status, master_url, worker_count, worker_cpu, worker_memory)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :master_url, :worker_count, :worker_cpu, :worker_memory)
	`, sc)
	return err
}

func (r *BigDataRepository) UpdateSparkCluster(ctx context.Context, sc *model.SparkCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE spark_clusters SET name=:name, status=:status, worker_count=:worker_count, updated_at=NOW() WHERE id=:id
	`, sc)
	return err
}

func (r *BigDataRepository) DeleteSparkCluster(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM spark_clusters WHERE id=$1`, id)
	return err
}

func (r *BigDataRepository) ListSparkApplications(ctx context.Context, sparkClusterID uuid.UUID) ([]model.SparkApplication, error) {
	var apps []model.SparkApplication
	err := r.db.SelectContext(ctx, &apps, `SELECT * FROM spark_applications WHERE spark_cluster_id=$1 ORDER BY created_at DESC`, sparkClusterID)
	return apps, err
}

func (r *BigDataRepository) FindSparkApplication(ctx context.Context, id uuid.UUID) (*model.SparkApplication, error) {
	return findByID[model.SparkApplication](ctx, r.db, "spark_applications", id)
}

// ── Flink ──────────────────────────────────────────────────────────────────

func (r *BigDataRepository) ListFlinkClusters(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.FlinkCluster, int, error) {
	return listWithOptionalFilter[model.FlinkCluster](ctx, r.db, "flink_clusters", "cluster_id", clusterID, limit, offset)
}

func (r *BigDataRepository) FindFlinkCluster(ctx context.Context, id uuid.UUID) (*model.FlinkCluster, error) {
	return findByID[model.FlinkCluster](ctx, r.db, "flink_clusters", id)
}

func (r *BigDataRepository) CreateFlinkCluster(ctx context.Context, fc *model.FlinkCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO flink_clusters (id, cluster_id, name, version, namespace, status, jobmanager_url, taskmanager_count, slots_per_tm)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :jobmanager_url, :taskmanager_count, :slots_per_tm)
	`, fc)
	return err
}

func (r *BigDataRepository) UpdateFlinkCluster(ctx context.Context, fc *model.FlinkCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE flink_clusters SET name=:name, status=:status, updated_at=NOW() WHERE id=:id
	`, fc)
	return err
}

func (r *BigDataRepository) DeleteFlinkCluster(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM flink_clusters WHERE id=$1`, id)
	return err
}

func (r *BigDataRepository) ListFlinkJobs(ctx context.Context, flinkClusterID uuid.UUID) ([]model.FlinkJob, error) {
	var jobs []model.FlinkJob
	err := r.db.SelectContext(ctx, &jobs, `SELECT * FROM flink_jobs WHERE flink_cluster_id=$1 ORDER BY created_at DESC`, flinkClusterID)
	return jobs, err
}

func (r *BigDataRepository) FindFlinkJob(ctx context.Context, id uuid.UUID) (*model.FlinkJob, error) {
	return findByID[model.FlinkJob](ctx, r.db, "flink_jobs", id)
}

func (r *BigDataRepository) UpdateFlinkJobStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE flink_jobs SET status=$1 WHERE id=$2`, status, id)
	return err
}

// ── Hive ───────────────────────────────────────────────────────────────────

func (r *BigDataRepository) ListHiveInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.HiveInstance, int, error) {
	return listWithOptionalFilter[model.HiveInstance](ctx, r.db, "hive_instances", "cluster_id", clusterID, limit, offset)
}

func (r *BigDataRepository) FindHiveInstance(ctx context.Context, id uuid.UUID) (*model.HiveInstance, error) {
	return findByID[model.HiveInstance](ctx, r.db, "hive_instances", id)
}

func (r *BigDataRepository) CreateHiveInstance(ctx context.Context, hi *model.HiveInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO hive_instances (id, cluster_id, name, version, namespace, status, metastore_url, hiveserver2_url)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :metastore_url, :hiveserver2_url)
	`, hi)
	return err
}

func (r *BigDataRepository) ListHiveTables(ctx context.Context, hiveInstanceID uuid.UUID) ([]model.HiveTable, error) {
	var tables []model.HiveTable
	err := r.db.SelectContext(ctx, &tables, `SELECT * FROM hive_tables WHERE hive_instance_id=$1 ORDER BY database_name, table_name`, hiveInstanceID)
	return tables, err
}

func (r *BigDataRepository) FindHiveTable(ctx context.Context, id uuid.UUID) (*model.HiveTable, error) {
	return findByID[model.HiveTable](ctx, r.db, "hive_tables", id)
}

// ── HDFS ───────────────────────────────────────────────────────────────────

func (r *BigDataRepository) ListHDFSClusters(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.HDFSCluster, int, error) {
	return listWithOptionalFilter[model.HDFSCluster](ctx, r.db, "hdfs_clusters", "cluster_id", clusterID, limit, offset)
}

func (r *BigDataRepository) FindHDFSCluster(ctx context.Context, id uuid.UUID) (*model.HDFSCluster, error) {
	return findByID[model.HDFSCluster](ctx, r.db, "hdfs_clusters", id)
}

func (r *BigDataRepository) CreateHDFSCluster(ctx context.Context, hc *model.HDFSCluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO hdfs_clusters (id, cluster_id, name, version, namespace, status, namenode_count, datanode_count, journalnode_count)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :namenode_count, :datanode_count, :journalnode_count)
	`, hc)
	return err
}

// ── NiFi ───────────────────────────────────────────────────────────────────

func (r *BigDataRepository) ListNiFiInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.NiFiInstance, int, error) {
	return listWithOptionalFilter[model.NiFiInstance](ctx, r.db, "nifi_instances", "cluster_id", clusterID, limit, offset)
}

func (r *BigDataRepository) FindNiFiInstance(ctx context.Context, id uuid.UUID) (*model.NiFiInstance, error) {
	return findByID[model.NiFiInstance](ctx, r.db, "nifi_instances", id)
}

func (r *BigDataRepository) CreateNiFiInstance(ctx context.Context, ni *model.NiFiInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO nifi_instances (id, cluster_id, name, version, namespace, status, ui_url, node_count)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :ui_url, :node_count)
	`, ni)
	return err
}

func (r *BigDataRepository) ListNiFiProcessGroups(ctx context.Context, nifiInstanceID uuid.UUID) ([]model.NiFiProcessGroup, error) {
	var groups []model.NiFiProcessGroup
	err := r.db.SelectContext(ctx, &groups, `SELECT * FROM nifi_process_groups WHERE nifi_instance_id=$1 ORDER BY name`, nifiInstanceID)
	return groups, err
}

// ── Generic Helpers ────────────────────────────────────────────────────────

func listWithOptionalFilter[T any](ctx context.Context, db *sqlx.DB, table, filterCol string, filterVal *uuid.UUID, limit, offset int) ([]T, int, error) {
	var items []T
	var total int

	query := `SELECT * FROM ` + table + ` WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM ` + table + ` WHERE 1=1`
	args := []any{}
	i := 1

	if filterVal != nil {
		cond := ` AND ` + filterCol + ` = $` + itoa(i)
		query += cond
		countQuery += cond
		args = append(args, *filterVal)
		i++
	}

	if err := db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY created_at DESC LIMIT $` + itoa(i) + ` OFFSET $` + itoa(i+1)
	args = append(args, limit, offset)

	err := db.SelectContext(ctx, &items, query, args...)
	return items, total, err
}

func findByID[T any](ctx context.Context, db *sqlx.DB, table string, id uuid.UUID) (*T, error) {
	var item T
	err := db.GetContext(ctx, &item, `SELECT * FROM `+table+` WHERE id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &item, err
}
