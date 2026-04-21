package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type BIRepository struct {
	db *sqlx.DB
}

func NewBIRepository(db *sqlx.DB) *BIRepository {
	return &BIRepository{db: db}
}

// ── Superset ───────────────────────────────────────────────────────────────

func (r *BIRepository) ListSupersetInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.SupersetInstance, int, error) {
	return listWithOptionalFilter[model.SupersetInstance](ctx, r.db, "superset_instances", "cluster_id", clusterID, limit, offset)
}

func (r *BIRepository) FindSupersetInstance(ctx context.Context, id uuid.UUID) (*model.SupersetInstance, error) {
	return findByID[model.SupersetInstance](ctx, r.db, "superset_instances", id)
}

func (r *BIRepository) CreateSupersetInstance(ctx context.Context, s *model.SupersetInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO superset_instances (id, cluster_id, name, version, namespace, status, url)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :url)
	`, s)
	return err
}

func (r *BIRepository) ListSupersetDashboards(ctx context.Context, instanceID uuid.UUID) ([]model.SupersetDashboard, error) {
	var items []model.SupersetDashboard
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM superset_dashboards WHERE superset_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

// ── Metabase ───────────────────────────────────────────────────────────────

func (r *BIRepository) ListMetabaseInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.MetabaseInstance, int, error) {
	return listWithOptionalFilter[model.MetabaseInstance](ctx, r.db, "metabase_instances", "cluster_id", clusterID, limit, offset)
}

func (r *BIRepository) FindMetabaseInstance(ctx context.Context, id uuid.UUID) (*model.MetabaseInstance, error) {
	return findByID[model.MetabaseInstance](ctx, r.db, "metabase_instances", id)
}

func (r *BIRepository) CreateMetabaseInstance(ctx context.Context, m *model.MetabaseInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO metabase_instances (id, cluster_id, name, version, namespace, status, url)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :url)
	`, m)
	return err
}

func (r *BIRepository) ListMetabaseDashboards(ctx context.Context, instanceID uuid.UUID) ([]model.MetabaseDashboard, error) {
	var items []model.MetabaseDashboard
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM metabase_dashboards WHERE metabase_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

// ── N8N ───────────────────────────────────────────────────────────────────

func (r *BIRepository) ListN8NInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.N8NInstance, int, error) {
	return listWithOptionalFilter[model.N8NInstance](ctx, r.db, "n8n_instances", "cluster_id", clusterID, limit, offset)
}

func (r *BIRepository) FindN8NInstance(ctx context.Context, id uuid.UUID) (*model.N8NInstance, error) {
	return findByID[model.N8NInstance](ctx, r.db, "n8n_instances", id)
}

func (r *BIRepository) CreateN8NInstance(ctx context.Context, n *model.N8NInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO n8n_instances (id, cluster_id, name, version, namespace, status, url)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :url)
	`, n)
	return err
}

func (r *BIRepository) ListN8NWorkflows(ctx context.Context, instanceID uuid.UUID) ([]model.N8NWorkflow, error) {
	var items []model.N8NWorkflow
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM n8n_workflows WHERE n8n_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

func (r *BIRepository) FindN8NWorkflow(ctx context.Context, id uuid.UUID) (*model.N8NWorkflow, error) {
	return findByID[model.N8NWorkflow](ctx, r.db, "n8n_workflows", id)
}

func (r *BIRepository) UpdateN8NWorkflowStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE n8n_workflows SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}
