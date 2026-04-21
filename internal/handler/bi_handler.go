package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/service"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	"github.com/orchestralog/api/pkg/response"
)

type BIHandler struct {
	svc      *service.BIService
	validate *validator.Validate
}

func NewBIHandler(svc *service.BIService) *BIHandler {
	return &BIHandler{svc: svc, validate: validator.New()}
}

func (h *BIHandler) apiErr(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		response.Error(w, apiErr)
		return
	}
	response.Error(w, apierror.ErrInternal)
}

func (h *BIHandler) decodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apierror.ErrBadRequest
	}
	if err := h.validate.Struct(dst); err != nil {
		return apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
	}
	return nil
}

func (h *BIHandler) MetricsStub(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "metrics not yet implemented"})
}

// ═══════════════════════════════════════════════════════════════════════════
// SUPERSET
// ═══════════════════════════════════════════════════════════════════════════

func (h *BIHandler) ListSupersetInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListSupersetInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BIHandler) GetSupersetInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetSupersetInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BIHandler) CreateSupersetInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateSupersetInstanceRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateSupersetInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BIHandler) ListSupersetDashboards(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListSupersetDashboards(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

// ═══════════════════════════════════════════════════════════════════════════
// METABASE
// ═══════════════════════════════════════════════════════════════════════════

func (h *BIHandler) ListMetabaseInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListMetabaseInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BIHandler) GetMetabaseInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetMetabaseInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BIHandler) CreateMetabaseInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateMetabaseInstanceRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateMetabaseInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BIHandler) ListMetabaseDashboards(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListMetabaseDashboards(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

// ═══════════════════════════════════════════════════════════════════════════
// N8N
// ═══════════════════════════════════════════════════════════════════════════

func (h *BIHandler) ListN8NInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListN8NInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *BIHandler) GetN8NInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetN8NInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BIHandler) CreateN8NInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateN8NInstanceRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateN8NInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *BIHandler) ListN8NWorkflows(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListN8NWorkflows(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *BIHandler) GetN8NWorkflow(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	wfID, err := parseID(r, "wfId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetN8NWorkflow(r.Context(), id, wfID)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *BIHandler) ActivateN8NWorkflow(w http.ResponseWriter, r *http.Request) {
	h.setWorkflowStatus(w, r, "active")
}

func (h *BIHandler) DeactivateN8NWorkflow(w http.ResponseWriter, r *http.Request) {
	h.setWorkflowStatus(w, r, "inactive")
}

func (h *BIHandler) setWorkflowStatus(w http.ResponseWriter, r *http.Request, status string) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	wfID, err := parseID(r, "wfId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.SetN8NWorkflowStatus(r.Context(), id, wfID, status); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": status})
}
