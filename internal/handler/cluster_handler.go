package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/middleware"
	"github.com/orchestralog/api/internal/service"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	"github.com/orchestralog/api/pkg/response"
)

type ClusterHandler struct {
	clusterService *service.ClusterService
	validate       *validator.Validate
}

func NewClusterHandler(clusterService *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{
		clusterService: clusterService,
		validate:       validator.New(),
	}
}

// GET /clusters
func (h *ClusterHandler) List(w http.ResponseWriter, r *http.Request) {
	p := pagination.Parse(r)
	status := r.URL.Query().Get("status")
	region := r.URL.Query().Get("region")

	clusters, meta, err := h.clusterService.List(r.Context(), status, region, p)
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSONWithMeta(w, http.StatusOK, clusters, meta)
}

// GET /clusters/:id
func (h *ClusterHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	cluster, svcErr := h.clusterService.GetByID(r.Context(), id)
	if svcErr != nil {
		if apiErr, ok := svcErr.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, cluster)
}

// POST /clusters
func (h *ClusterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}

	cluster, err := h.clusterService.Create(r.Context(), &req, middleware.GetUserID(r))
	if err != nil {
		if apiErr, ok := err.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusCreated, cluster)
}

// PUT /clusters/:id
func (h *ClusterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	var req request.UpdateClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	cluster, svcErr := h.clusterService.Update(r.Context(), id, &req)
	if svcErr != nil {
		if apiErr, ok := svcErr.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, cluster)
}

// DELETE /clusters/:id
func (h *ClusterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	if svcErr := h.clusterService.Delete(r.Context(), id); svcErr != nil {
		if apiErr, ok := svcErr.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "cluster deleted"})
}

// GET /clusters/:id/namespaces
func (h *ClusterHandler) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	clusterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	namespaces, svcErr := h.clusterService.ListNamespaces(r.Context(), clusterID)
	if svcErr != nil {
		if apiErr, ok := svcErr.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, namespaces)
}

// GET /clusters/:id/namespaces/:nsId/resources
func (h *ClusterHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	clusterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	nsID, err := uuid.Parse(chi.URLParam(r, "nsId"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}

	resources, svcErr := h.clusterService.ListResources(r.Context(), clusterID, nsID)
	if svcErr != nil {
		if apiErr, ok := svcErr.(*apierror.APIError); ok {
			response.Error(w, apiErr)
			return
		}
		response.Error(w, apierror.ErrInternal)
		return
	}

	response.JSON(w, http.StatusOK, resources)
}
