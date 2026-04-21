package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/model"
)

type DataFlowRepository struct {
	db *sqlx.DB
}

func NewDataFlowRepository(db *sqlx.DB) *DataFlowRepository {
	return &DataFlowRepository{db: db}
}

func (r *DataFlowRepository) List(ctx context.Context, status string, limit, offset int) ([]model.DataFlow, int, error) {
	var items []model.DataFlow
	var total int

	query := `SELECT * FROM data_flows WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM data_flows WHERE 1=1`
	args := []any{}
	i := 1

	if status != "" {
		query += ` AND status=$` + itoa(i)
		countQuery += ` AND status=$` + itoa(i)
		args = append(args, status)
		i++
	}

	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY updated_at DESC LIMIT $` + itoa(i) + ` OFFSET $` + itoa(i+1)
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, total, err
}

func (r *DataFlowRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.DataFlow, error) {
	return findByID[model.DataFlow](ctx, r.db, "data_flows", id)
}

func (r *DataFlowRepository) Create(ctx context.Context, df *model.DataFlow) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO data_flows (id, name, description, status, template_id, layout, created_by)
		VALUES (:id, :name, :description, :status, :template_id, :layout, :created_by)
	`, df)
	return err
}

func (r *DataFlowRepository) Update(ctx context.Context, df *model.DataFlow) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE data_flows SET name=:name, description=:description, status=:status, updated_at=NOW()
		WHERE id=:id
	`, df)
	return err
}

func (r *DataFlowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM data_flows WHERE id=$1`, id)
	return err
}

func (r *DataFlowRepository) ListNodes(ctx context.Context, flowID uuid.UUID) ([]model.FlowNode, error) {
	var items []model.FlowNode
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM flow_nodes WHERE data_flow_id=$1 ORDER BY position_x`, flowID)
	return items, err
}

func (r *DataFlowRepository) ListEdges(ctx context.Context, flowID uuid.UUID) ([]model.FlowEdge, error) {
	var items []model.FlowEdge
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM flow_edges WHERE data_flow_id=$1`, flowID)
	return items, err
}

func (r *DataFlowRepository) ReplaceNodes(ctx context.Context, flowID uuid.UUID, nodes []model.FlowNode) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM flow_nodes WHERE data_flow_id=$1`, flowID); err != nil {
		return err
	}
	for _, n := range nodes {
		if _, err := tx.NamedExecContext(ctx, `
			INSERT INTO flow_nodes (id, data_flow_id, node_id, label, resource_type, cluster, namespace, status, position_x, position_y, config)
			VALUES (:id, :data_flow_id, :node_id, :label, :resource_type, :cluster, :namespace, :status, :position_x, :position_y, :config)
		`, n); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *DataFlowRepository) ReplaceEdges(ctx context.Context, flowID uuid.UUID, edges []model.FlowEdge) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM flow_edges WHERE data_flow_id=$1`, flowID); err != nil {
		return err
	}
	for _, e := range edges {
		if _, err := tx.NamedExecContext(ctx, `
			INSERT INTO flow_edges (id, data_flow_id, edge_id, source_node_id, target_node_id, flow_type, label, animated, config)
			VALUES (:id, :data_flow_id, :edge_id, :source_node_id, :target_node_id, :flow_type, :label, :animated, :config)
		`, e); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *DataFlowRepository) ListTemplates(ctx context.Context) ([]model.FlowTemplate, error) {
	var items []model.FlowTemplate
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM flow_templates ORDER BY name`)
	return items, err
}
