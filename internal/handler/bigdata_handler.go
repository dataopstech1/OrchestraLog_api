package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/service"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	"github.com/orchestralog/api/pkg/response"
)

type BigDataHandler struct {
	svc      *service.BigDataService
	validate *validator.Validate
}

func NewBigDataHandler(svc *service.BigDataService) *BigDataHandler {
	return &BigDataHandler{svc: svc, validate: validator.New()}
}

func (h *BigDataHandler) apiErr(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		response.Error(w, apiErr)
		return
	}
	response.Error(w, apierror.ErrInternal)
}

func parseID(r *http.Request, param string) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, param))
}

func parseURLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func optionalClusterID(r *http.Request) *uuid.UUID {
	raw := r.URL.Query().Get("cluster_id")
	if raw == "" {
		return nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil
	}
	return &id
}

// ═══════════════════════════════════════════════════════════════════════════
// SPARK
// ═══════════════════════════════════════════════════════════════════════════

func (h *BigDataHandler) ListSparkClusters(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListSparkClusters(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BigDataHandler) GetSparkCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetSparkCluster(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) CreateSparkCluster(w http.ResponseWriter, r *http.Request) {
	var req request.CreateSparkClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	item, svcErr := h.svc.CreateSparkCluster(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BigDataHandler) UpdateSparkCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateSparkClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.UpdateSparkCluster(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) DeleteSparkCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.DeleteSparkCluster(r.Context(), id); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "spark cluster deleted"})
}

func (h *BigDataHandler) ListSparkApplications(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	apps, svcErr := h.svc.ListSparkApplications(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, apps)
}

func (h *BigDataHandler) GetSparkApplication(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	appID, err := parseID(r, "appId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	app, svcErr := h.svc.GetSparkApplication(r.Context(), id, appID)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, app)
}

// ═══════════════════════════════════════════════════════════════════════════
// FLINK
// ═══════════════════════════════════════════════════════════════════════════

func (h *BigDataHandler) ListFlinkClusters(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListFlinkClusters(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BigDataHandler) GetFlinkCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetFlinkCluster(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) CreateFlinkCluster(w http.ResponseWriter, r *http.Request) {
	var req request.CreateFlinkClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	item, svcErr := h.svc.CreateFlinkCluster(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BigDataHandler) UpdateFlinkCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateFlinkClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.UpdateFlinkCluster(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) DeleteFlinkCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.DeleteFlinkCluster(r.Context(), id); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "flink cluster deleted"})
}

func (h *BigDataHandler) ListFlinkJobs(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	jobs, svcErr := h.svc.ListFlinkJobs(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, jobs)
}

func (h *BigDataHandler) GetFlinkJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	jobID, err := parseID(r, "jobId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	job, svcErr := h.svc.GetFlinkJob(r.Context(), id, jobID)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, job)
}

func (h *BigDataHandler) CancelFlinkJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	jobID, err := parseID(r, "jobId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.CancelFlinkJob(r.Context(), id, jobID); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "job canceled"})
}

// ═══════════════════════════════════════════════════════════════════════════
// HIVE
// ═══════════════════════════════════════════════════════════════════════════

func (h *BigDataHandler) ListHiveInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListHiveInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BigDataHandler) GetHiveInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetHiveInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) CreateHiveInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateHiveInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	item, svcErr := h.svc.CreateHiveInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BigDataHandler) ListHiveTables(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	tables, svcErr := h.svc.ListHiveTables(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, tables)
}

func (h *BigDataHandler) GetHiveTable(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	tableID, err := parseID(r, "tableId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	t, svcErr := h.svc.GetHiveTable(r.Context(), id, tableID)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, t)
}

// ═══════════════════════════════════════════════════════════════════════════
// HDFS
// ═══════════════════════════════════════════════════════════════════════════

func (h *BigDataHandler) ListHDFSClusters(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListHDFSClusters(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BigDataHandler) GetHDFSCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetHDFSCluster(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) CreateHDFSCluster(w http.ResponseWriter, r *http.Request) {
	var req request.CreateHDFSClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	item, svcErr := h.svc.CreateHDFSCluster(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

// ═══════════════════════════════════════════════════════════════════════════
// NIFI
// ═══════════════════════════════════════════════════════════════════════════

func (h *BigDataHandler) ListNiFiInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListNiFiInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BigDataHandler) GetNiFiInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetNiFiInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BigDataHandler) CreateNiFiInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateNiFiInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	item, svcErr := h.svc.CreateNiFiInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BigDataHandler) ListNiFiProcessGroups(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	groups, svcErr := h.svc.ListNiFiProcessGroups(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, groups)
}

// MetricsStub tüm servisler için placeholder metrics döner
func (h *BigDataHandler) MetricsStub(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "metrics not yet implemented"})
}
