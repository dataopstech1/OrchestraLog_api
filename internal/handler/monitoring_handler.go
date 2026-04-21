package handler

import (
	"net/http"

	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/response"
)

type MonitoringHandler struct{}

func NewMonitoringHandler() *MonitoringHandler { return &MonitoringHandler{} }

func (h *MonitoringHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"cluster_id": id,
		"cpu":        map[string]any{"current": 62.5, "average": 58.3, "peak": 85.2, "history": []any{}},
		"memory":     map[string]any{"current": 71.3, "average": 68.9, "peak": 89.5, "history": []any{}},
		"storage":    map[string]any{"total_bytes": 1099511627776, "used_bytes": 503316480000, "utilization": 45.8},
		"pods":       map[string]any{"running": 85, "pending": 3, "failed": 1, "succeeded": 12, "total": 101},
		"network":    map[string]any{"bytes_in_per_sec": 125000000, "bytes_out_per_sec": 98000000},
	})
}

func (h *MonitoringHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"cluster_id": id, "nodes": []any{}})
}

func (h *MonitoringHandler) GetNode(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MonitoringHandler) GetResources(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MonitoringHandler) ListPods(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"pods": []any{}})
}

func (h *MonitoringHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"events": []any{}})
}

func (h *MonitoringHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"alerts": []any{}})
}
