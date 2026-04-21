package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type MLAIRepository struct {
	db *sqlx.DB
}

func NewMLAIRepository(db *sqlx.DB) *MLAIRepository {
	return &MLAIRepository{db: db}
}

// ── MLFlow ─────────────────────────────────────────────────────────────────

func (r *MLAIRepository) ListMLFlowInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.MLFlowInstance, int, error) {
	return listWithOptionalFilter[model.MLFlowInstance](ctx, r.db, "mlflow_instances", "cluster_id", clusterID, limit, offset)
}

func (r *MLAIRepository) FindMLFlowInstance(ctx context.Context, id uuid.UUID) (*model.MLFlowInstance, error) {
	return findByID[model.MLFlowInstance](ctx, r.db, "mlflow_instances", id)
}

func (r *MLAIRepository) CreateMLFlowInstance(ctx context.Context, m *model.MLFlowInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO mlflow_instances (id, cluster_id, name, version, namespace, status, tracking_url, artifact_store)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :tracking_url, :artifact_store)
	`, m)
	return err
}

func (r *MLAIRepository) ListExperiments(ctx context.Context, instanceID uuid.UUID) ([]model.MLFlowExperiment, error) {
	var items []model.MLFlowExperiment
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM mlflow_experiments WHERE mlflow_instance_id=$1 ORDER BY created_at DESC`, instanceID)
	return items, err
}

func (r *MLAIRepository) FindExperiment(ctx context.Context, id uuid.UUID) (*model.MLFlowExperiment, error) {
	return findByID[model.MLFlowExperiment](ctx, r.db, "mlflow_experiments", id)
}

func (r *MLAIRepository) ListModels(ctx context.Context, instanceID uuid.UUID) ([]model.MLFlowModel, error) {
	var items []model.MLFlowModel
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM mlflow_models WHERE mlflow_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

func (r *MLAIRepository) FindModel(ctx context.Context, id uuid.UUID) (*model.MLFlowModel, error) {
	return findByID[model.MLFlowModel](ctx, r.db, "mlflow_models", id)
}

func (r *MLAIRepository) UpdateModelStage(ctx context.Context, id uuid.UUID, stage string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE mlflow_models SET stage=$1, updated_at=NOW() WHERE id=$2`, stage, id)
	return err
}

// ── Feast ──────────────────────────────────────────────────────────────────

func (r *MLAIRepository) ListFeastInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.FeastInstance, int, error) {
	return listWithOptionalFilter[model.FeastInstance](ctx, r.db, "feast_instances", "cluster_id", clusterID, limit, offset)
}

func (r *MLAIRepository) FindFeastInstance(ctx context.Context, id uuid.UUID) (*model.FeastInstance, error) {
	return findByID[model.FeastInstance](ctx, r.db, "feast_instances", id)
}

func (r *MLAIRepository) CreateFeastInstance(ctx context.Context, fi *model.FeastInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO feast_instances (id, cluster_id, name, version, namespace, status, online_store, offline_store)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :online_store, :offline_store)
	`, fi)
	return err
}

func (r *MLAIRepository) ListFeastEntities(ctx context.Context, instanceID uuid.UUID) ([]model.FeastEntity, error) {
	var items []model.FeastEntity
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM feast_entities WHERE feast_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

func (r *MLAIRepository) ListFeastFeatureViews(ctx context.Context, instanceID uuid.UUID) ([]model.FeastFeatureView, error) {
	var items []model.FeastFeatureView
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM feast_feature_views WHERE feast_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

// ── JupyterHub ─────────────────────────────────────────────────────────────

func (r *MLAIRepository) ListJupyterHubInstances(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.JupyterHubInstance, int, error) {
	return listWithOptionalFilter[model.JupyterHubInstance](ctx, r.db, "jupyterhub_instances", "cluster_id", clusterID, limit, offset)
}

func (r *MLAIRepository) FindJupyterHubInstance(ctx context.Context, id uuid.UUID) (*model.JupyterHubInstance, error) {
	return findByID[model.JupyterHubInstance](ctx, r.db, "jupyterhub_instances", id)
}

func (r *MLAIRepository) CreateJupyterHubInstance(ctx context.Context, jh *model.JupyterHubInstance) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO jupyterhub_instances (id, cluster_id, name, version, namespace, status, hub_url, max_users)
		VALUES (:id, :cluster_id, :name, :version, :namespace, :status, :hub_url, :max_users)
	`, jh)
	return err
}

func (r *MLAIRepository) ListNotebooks(ctx context.Context, instanceID uuid.UUID) ([]model.JupyterHubNotebook, error) {
	var items []model.JupyterHubNotebook
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM jupyterhub_notebooks WHERE jupyterhub_instance_id=$1 ORDER BY name`, instanceID)
	return items, err
}

func (r *MLAIRepository) FindNotebook(ctx context.Context, id uuid.UUID) (*model.JupyterHubNotebook, error) {
	return findByID[model.JupyterHubNotebook](ctx, r.db, "jupyterhub_notebooks", id)
}

func (r *MLAIRepository) CreateNotebook(ctx context.Context, nb *model.JupyterHubNotebook) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO jupyterhub_notebooks (id, jupyterhub_instance_id, name, owner, status, image, cpu_limit, memory_limit, gpu_limit)
		VALUES (:id, :jupyterhub_instance_id, :name, :owner, :status, :image, :cpu_limit, :memory_limit, :gpu_limit)
	`, nb)
	return err
}

func (r *MLAIRepository) UpdateNotebook(ctx context.Context, nb *model.JupyterHubNotebook) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE jupyterhub_notebooks SET
			image=:image, cpu_limit=:cpu_limit, memory_limit=:memory_limit,
			gpu_limit=:gpu_limit, status=:status, updated_at=NOW()
		WHERE id=:id
	`, nb)
	return err
}

func (r *MLAIRepository) DeleteNotebook(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM jupyterhub_notebooks WHERE id=$1`, id)
	return err
}

func (r *MLAIRepository) UpdateNotebookStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE jupyterhub_notebooks SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

// ── LLM ───────────────────────────────────────────────────────────────────

func (r *MLAIRepository) ListLLMDeployments(ctx context.Context, clusterID *uuid.UUID, limit, offset int) ([]model.LLMDeployment, int, error) {
	return listWithOptionalFilter[model.LLMDeployment](ctx, r.db, "llm_deployments", "cluster_id", clusterID, limit, offset)
}

func (r *MLAIRepository) FindLLMDeployment(ctx context.Context, id uuid.UUID) (*model.LLMDeployment, error) {
	return findByID[model.LLMDeployment](ctx, r.db, "llm_deployments", id)
}

func (r *MLAIRepository) CreateLLMDeployment(ctx context.Context, d *model.LLMDeployment) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO llm_deployments (id, cluster_id, name, model_name, model_version, namespace, status, endpoint_url, replicas, gpu_count, gpu_type, max_tokens, context_window)
		VALUES (:id, :cluster_id, :name, :model_name, :model_version, :namespace, :status, :endpoint_url, :replicas, :gpu_count, :gpu_type, :max_tokens, :context_window)
	`, d)
	return err
}

func (r *MLAIRepository) UpdateLLMDeployment(ctx context.Context, d *model.LLMDeployment) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE llm_deployments SET name=:name, status=:status, replicas=:replicas, updated_at=NOW() WHERE id=:id
	`, d)
	return err
}

func (r *MLAIRepository) DeleteLLMDeployment(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM llm_deployments WHERE id=$1`, id)
	return err
}
