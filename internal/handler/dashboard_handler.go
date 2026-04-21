package handler

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/pkg/response"
)

type DashboardHandler struct {
	db *sqlx.DB
}

func NewDashboardHandler(db *sqlx.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GET /dashboard/summary
func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clusterStats := h.countByStatus(ctx, "clusters", []string{"healthy", "warning", "critical", "offline"})
	flowStats := h.countByStatus(ctx, "data_flows", []string{"running", "stopped", "failed", "draft", "deployed"})

	response.JSON(w, http.StatusOK, map[string]any{
		"clusters": map[string]any{
			"total":    clusterStats["total"],
			"healthy":  clusterStats["healthy"],
			"warning":  clusterStats["warning"],
			"critical": clusterStats["critical"],
		},
		"services": map[string]any{
			"kafka":      h.serviceCount(ctx, "kafka_clusters"),
			"spark":      h.serviceCount(ctx, "spark_clusters"),
			"flink":      h.serviceCount(ctx, "flink_clusters"),
			"hive":       h.serviceCount(ctx, "hive_instances"),
			"hdfs":       h.serviceCount(ctx, "hdfs_clusters"),
			"nifi":       h.serviceCount(ctx, "nifi_instances"),
			"mlflow":     h.serviceCount(ctx, "mlflow_instances"),
			"feast":      h.serviceCount(ctx, "feast_instances"),
			"jupyterhub": h.serviceCount(ctx, "jupyterhub_instances"),
			"llm":        h.serviceCount(ctx, "llm_deployments"),
			"superset":   h.serviceCount(ctx, "superset_instances"),
			"metabase":   h.serviceCount(ctx, "metabase_instances"),
			"n8n":        h.serviceCount(ctx, "n8n_instances"),
		},
		"data_flows": map[string]any{
			"total":    flowStats["total"],
			"running":  flowStats["running"],
			"stopped":  flowStats["stopped"],
			"failed":   flowStats["failed"],
		},
	})
}

// GET /dashboard/services-status
func (h *DashboardHandler) ServicesStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response.JSON(w, http.StatusOK, map[string]any{
		"kafka":      h.serviceStatusCounts(ctx, "kafka_clusters"),
		"spark":      h.serviceStatusCounts(ctx, "spark_clusters"),
		"flink":      h.serviceStatusCounts(ctx, "flink_clusters"),
		"hive":       h.serviceStatusCounts(ctx, "hive_instances"),
		"hdfs":       h.serviceStatusCounts(ctx, "hdfs_clusters"),
		"nifi":       h.serviceStatusCounts(ctx, "nifi_instances"),
		"mlflow":     h.serviceStatusCounts(ctx, "mlflow_instances"),
		"feast":      h.serviceStatusCounts(ctx, "feast_instances"),
		"jupyterhub": h.serviceStatusCounts(ctx, "jupyterhub_instances"),
		"llm":        h.serviceStatusCounts(ctx, "llm_deployments"),
		"superset":   h.serviceStatusCounts(ctx, "superset_instances"),
		"metabase":   h.serviceStatusCounts(ctx, "metabase_instances"),
		"n8n":        h.serviceStatusCounts(ctx, "n8n_instances"),
	})
}

// GET /dashboard/recent-activity
func (h *DashboardHandler) RecentActivity(w http.ResponseWriter, r *http.Request) {
	var logs []map[string]any
	rows, err := h.db.QueryxContext(r.Context(), `
		SELECT al.id, al.action, al.entity_type, al.entity_id, al.created_at,
		       u.first_name || ' ' || u.last_name AS user_name, u.email
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		ORDER BY al.created_at DESC LIMIT 20
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			row := make(map[string]any)
			if err := rows.MapScan(row); err == nil {
				logs = append(logs, row)
			}
		}
	}
	if logs == nil {
		logs = []map[string]any{}
	}
	response.JSON(w, http.StatusOK, logs)
}

// GET /dashboard/alerts
func (h *DashboardHandler) Alerts(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, []any{})
}

// ── helpers ────────────────────────────────────────────────────────────────

func (h *DashboardHandler) countByStatus(ctx context.Context, table string, statuses []string) map[string]int {
	result := map[string]int{"total": 0}
	for _, s := range statuses {
		result[s] = 0
	}
	rows, err := h.db.QueryxContext(ctx, `SELECT status, COUNT(*) as cnt FROM `+table+` GROUP BY status`)
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var cnt int
		if err := rows.Scan(&status, &cnt); err == nil {
			result[status] = cnt
			result["total"] += cnt
		}
	}
	return result
}

func (h *DashboardHandler) serviceCount(ctx context.Context, table string) map[string]int {
	var total int
	h.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM `+table)
	return map[string]int{"instances": total}
}

func (h *DashboardHandler) serviceStatusCounts(ctx context.Context, table string) map[string]int {
	result := map[string]int{"total": 0, "running": 0, "healthy": 0, "failed": 0}
	rows, err := h.db.QueryxContext(ctx, `SELECT status, COUNT(*) as cnt FROM `+table+` GROUP BY status`)
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var cnt int
		if err := rows.Scan(&status, &cnt); err == nil {
			result[status] = cnt
			result["total"] += cnt
		}
	}
	return result
}
