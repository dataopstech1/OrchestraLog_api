package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/middleware"
	"github.com/orchestralog/api/internal/service"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	"github.com/orchestralog/api/pkg/response"
)

type DataFlowHandler struct {
	svc      *service.DataFlowService
	validate *validator.Validate
}

func NewDataFlowHandler(svc *service.DataFlowService) *DataFlowHandler {
	return &DataFlowHandler{svc: svc, validate: validator.New()}
}

func (h *DataFlowHandler) apiErr(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		response.Error(w, apiErr)
		return
	}
	response.Error(w, apierror.ErrInternal)
}

// GET /data-flows
func (h *DataFlowHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	items, meta, err := h.svc.List(r.Context(), status, pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

// GET /data-flows/:id
func (h *DataFlowHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetByID(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

// POST /data-flows
func (h *DataFlowHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateDataFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	item, svcErr := h.svc.Create(r.Context(), &req, middleware.GetUserID(r))
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

// PUT /data-flows/:id
func (h *DataFlowHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateDataFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.Update(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

// DELETE /data-flows/:id
func (h *DataFlowHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.Delete(r.Context(), id); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "data flow deleted"})
}

// PUT /data-flows/:id/nodes
func (h *DataFlowHandler) UpdateNodes(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateFlowNodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	nodes, svcErr := h.svc.UpdateNodes(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, nodes)
}

// PUT /data-flows/:id/edges
func (h *DataFlowHandler) UpdateEdges(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateFlowEdgesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	edges, svcErr := h.svc.UpdateEdges(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, edges)
}

// POST /data-flows/:id/deploy
func (h *DataFlowHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.SetStatus(r.Context(), id, "deployed")
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

// POST /data-flows/:id/stop
func (h *DataFlowHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.SetStatus(r.Context(), id, "stopped")
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

// GET /data-flows/templates
func (h *DataFlowHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	items, svcErr := h.svc.ListTemplates(r.Context())
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}
