package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type ClusterRepository struct {
	db *sqlx.DB
}

func NewClusterRepository(db *sqlx.DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

func (r *ClusterRepository) List(ctx context.Context, status, region string, limit, offset int) ([]model.Cluster, int, error) {
	var clusters []model.Cluster
	var total int

	query := `SELECT * FROM clusters WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM clusters WHERE 1=1`
	args := []any{}
	i := 1

	if status != "" {
		query += ` AND status = $` + itoa(i)
		countQuery += ` AND status = $` + itoa(i)
		args = append(args, status)
		i++
	}
	if region != "" {
		query += ` AND region = $` + itoa(i)
		countQuery += ` AND region = $` + itoa(i)
		args = append(args, region)
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

func (r *ClusterRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Cluster, error) {
	var cluster model.Cluster
	err := r.db.GetContext(ctx, &cluster, `SELECT * FROM clusters WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &cluster, err
}

func (r *ClusterRepository) Create(ctx context.Context, c *model.Cluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO clusters (id, name, status, region, k8s_version, api_server_url, kubeconfig, nodes_total, nodes_ready, created_by)
		VALUES (:id, :name, :status, :region, :k8s_version, :api_server_url, :kubeconfig, :nodes_total, :nodes_ready, :created_by)
	`, c)
	return err
}

func (r *ClusterRepository) Update(ctx context.Context, c *model.Cluster) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE clusters SET
			name = :name, status = :status, k8s_version = :k8s_version,
			api_server_url = :api_server_url, kubeconfig = :kubeconfig,
			nodes_total = :nodes_total, nodes_ready = :nodes_ready,
			updated_at = NOW()
		WHERE id = :id
	`, c)
	return err
}

func (r *ClusterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM clusters WHERE id = $1`, id)
	return err
}

func (r *ClusterRepository) ListNamespaces(ctx context.Context, clusterID uuid.UUID) ([]model.Namespace, error) {
	var namespaces []model.Namespace
	err := r.db.SelectContext(ctx, &namespaces, `
		SELECT * FROM namespaces WHERE cluster_id = $1 ORDER BY name
	`, clusterID)
	return namespaces, err
}

func (r *ClusterRepository) ListResources(ctx context.Context, namespaceID uuid.UUID) ([]model.Resource, error) {
	var resources []model.Resource
	err := r.db.SelectContext(ctx, &resources, `
		SELECT * FROM resources WHERE namespace_id = $1 ORDER BY type, name
	`, namespaceID)
	return resources, err
}

func (r *ClusterRepository) CountNamespaces(ctx context.Context, clusterID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM namespaces WHERE cluster_id = $1`, clusterID)
	return count, err
}

func (r *ClusterRepository) CountResources(ctx context.Context, clusterID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM resources r
		JOIN namespaces n ON r.namespace_id = n.id
		WHERE n.cluster_id = $1
	`, clusterID)
	return count, err
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
