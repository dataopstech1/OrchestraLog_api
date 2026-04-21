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

type MLAIHandler struct {
	svc      *service.MLAIService
	validate *validator.Validate
}

func NewMLAIHandler(svc *service.MLAIService) *MLAIHandler {
	return &MLAIHandler{svc: svc, validate: validator.New()}
}

func (h *MLAIHandler) apiErr(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		response.Error(w, apiErr)
		return
	}
	response.Error(w, apierror.ErrInternal)
}

func (h *MLAIHandler) decodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apierror.ErrBadRequest
	}
	if err := h.validate.Struct(dst); err != nil {
		return apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════════════════
// MLFLOW
// ═══════════════════════════════════════════════════════════════════════════

func (h *MLAIHandler) ListMLFlowInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListMLFlowInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *MLAIHandler) GetMLFlowInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetMLFlowInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) CreateMLFlowInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateMLFlowInstanceRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateMLFlowInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *MLAIHandler) ListExperiments(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListExperiments(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *MLAIHandler) GetExperiment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	expID, err := parseID(r, "expId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetExperiment(r.Context(), id, expID)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListModels(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *MLAIHandler) GetModel(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	modelID, err := parseID(r, "modelId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetModel(r.Context(), id, modelID)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) UpdateModelStage(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	modelID, err := parseID(r, "modelId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateMLFlowModelStageRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.UpdateModelStage(r.Context(), id, modelID, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

// ═══════════════════════════════════════════════════════════════════════════
// FEAST
// ═══════════════════════════════════════════════════════════════════════════

func (h *MLAIHandler) ListFeastInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListFeastInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *MLAIHandler) GetFeastInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetFeastInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) CreateFeastInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateFeastInstanceRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateFeastInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *MLAIHandler) ListFeastEntities(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListFeastEntities(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *MLAIHandler) ListFeastFeatureViews(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListFeastFeatureViews(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

// ═══════════════════════════════════════════════════════════════════════════
// JUPYTERHUB
// ═══════════════════════════════════════════════════════════════════════════

func (h *MLAIHandler) ListJupyterHubInstances(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListJupyterHubInstances(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *MLAIHandler) GetJupyterHubInstance(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetJupyterHubInstance(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) CreateJupyterHubInstance(w http.ResponseWriter, r *http.Request) {
	var req request.CreateJupyterHubInstanceRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateJupyterHubInstance(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *MLAIHandler) ListNotebooks(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	items, svcErr := h.svc.ListNotebooks(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *MLAIHandler) CreateNotebook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.CreateNotebookRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateNotebook(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *MLAIHandler) UpdateNotebook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	nbID, err := parseID(r, "nbId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.UpdateNotebook(r.Context(), id, nbID, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) DeleteNotebook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	nbID, err := parseID(r, "nbId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.DeleteNotebook(r.Context(), id, nbID); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "notebook deleted"})
}

func (h *MLAIHandler) StartNotebook(w http.ResponseWriter, r *http.Request) {
	h.setNotebookStatus(w, r, "Running")
}

func (h *MLAIHandler) StopNotebook(w http.ResponseWriter, r *http.Request) {
	h.setNotebookStatus(w, r, "Stopped")
}

func (h *MLAIHandler) setNotebookStatus(w http.ResponseWriter, r *http.Request, status string) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	nbID, err := parseID(r, "nbId")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.SetNotebookStatus(r.Context(), id, nbID, status); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": status})
}

// ═══════════════════════════════════════════════════════════════════════════
// LLM
// ═══════════════════════════════════════════════════════════════════════════

func (h *MLAIHandler) ListLLMDeployments(w http.ResponseWriter, r *http.Request) {
	items, meta, err := h.svc.ListLLMDeployments(r.Context(), optionalClusterID(r), pagination.Parse(r))
	if err != nil {
		h.apiErr(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, items, meta)
}

func (h *MLAIHandler) GetLLMDeployment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.GetLLMDeployment(r.Context(), id)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) CreateLLMDeployment(w http.ResponseWriter, r *http.Request) {
	var req request.CreateLLMDeploymentRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.CreateLLMDeployment(r.Context(), &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *MLAIHandler) UpdateLLMDeployment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateLLMDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	item, svcErr := h.svc.UpdateLLMDeployment(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) DeleteLLMDeployment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.svc.DeleteLLMDeployment(r.Context(), id); svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "llm deployment deleted"})
}

func (h *MLAIHandler) ScaleLLMDeployment(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.ScaleLLMDeploymentRequest
	if err := h.decodeAndValidate(r, &req); err != nil {
		h.apiErr(w, err)
		return
	}
	item, svcErr := h.svc.ScaleLLMDeployment(r.Context(), id, &req)
	if svcErr != nil {
		h.apiErr(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *MLAIHandler) MetricsStub(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "metrics not yet implemented"})
}

