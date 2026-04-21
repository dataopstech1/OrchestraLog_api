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

type KafkaHandler struct {
	kafkaService *service.KafkaService
	validate     *validator.Validate
}

func NewKafkaHandler(kafkaService *service.KafkaService) *KafkaHandler {
	return &KafkaHandler{kafkaService: kafkaService, validate: validator.New()}
}

func (h *KafkaHandler) respondError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*apierror.APIError); ok {
		response.Error(w, apiErr)
		return
	}
	response.Error(w, apierror.ErrInternal)
}

// GET /kafka/clusters
func (h *KafkaHandler) ListClusters(w http.ResponseWriter, r *http.Request) {
	p := pagination.Parse(r)
	var clusterID *uuid.UUID
	if raw := r.URL.Query().Get("cluster_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			response.Error(w, apierror.ErrBadRequest)
			return
		}
		clusterID = &id
	}

	clusters, meta, err := h.kafkaService.ListClusters(r.Context(), clusterID, p)
	if err != nil {
		h.respondError(w, err)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, clusters, meta)
}

// GET /kafka/clusters/:id
func (h *KafkaHandler) GetCluster(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	kc, svcErr := h.kafkaService.GetClusterByID(r.Context(), id)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, kc)
}

// POST /kafka/clusters
func (h *KafkaHandler) CreateCluster(w http.ResponseWriter, r *http.Request) {
	var req request.CreateKafkaClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	kc, err := h.kafkaService.CreateCluster(r.Context(), &req, middleware.GetUserID(r))
	if err != nil {
		h.respondError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, kc)
}

// PUT /kafka/clusters/:id
func (h *KafkaHandler) UpdateCluster(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateKafkaClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	kc, svcErr := h.kafkaService.UpdateCluster(r.Context(), id, &req)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, kc)
}

// DELETE /kafka/clusters/:id
func (h *KafkaHandler) DeleteCluster(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.kafkaService.DeleteCluster(r.Context(), id); svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "kafka cluster deleted"})
}

// GET /kafka/clusters/:id/metrics
func (h *KafkaHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	metrics, svcErr := h.kafkaService.GetMetrics(r.Context(), id)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, metrics)
}

// GET /kafka/clusters/:id/topics
func (h *KafkaHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	p := pagination.Parse(r)
	topics, meta, svcErr := h.kafkaService.ListTopics(r.Context(), id, p)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSONWithMeta(w, http.StatusOK, topics, meta)
}

// POST /kafka/clusters/:id/topics
func (h *KafkaHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.CreateKafkaTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	t, svcErr := h.kafkaService.CreateTopic(r.Context(), id, &req)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, t)
}

// PUT /kafka/clusters/:id/topics/:topicId
func (h *KafkaHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	topicID, err := uuid.Parse(chi.URLParam(r, "topicId"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateKafkaTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	t, svcErr := h.kafkaService.UpdateTopic(r.Context(), id, topicID, &req)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, t)
}

// DELETE /kafka/clusters/:id/topics/:topicId
func (h *KafkaHandler) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	topicID, err := uuid.Parse(chi.URLParam(r, "topicId"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.kafkaService.DeleteTopic(r.Context(), id, topicID); svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "topic deleted"})
}

// GET /kafka/clusters/:id/consumers
func (h *KafkaHandler) ListConsumers(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	groups, svcErr := h.kafkaService.ListConsumerGroups(r.Context(), id)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, groups)
}

// GET /kafka/clusters/:id/consumers/:groupId
func (h *KafkaHandler) GetConsumer(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	g, svcErr := h.kafkaService.GetConsumerGroup(r.Context(), id, groupID)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, g)
}

// GET /kafka/clusters/:id/resources
func (h *KafkaHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	resources, svcErr := h.kafkaService.ListResources(r.Context(), id)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, resources)
}

// POST /kafka/clusters/:id/resources
func (h *KafkaHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.CreateKafkaResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apierror.NewWithDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error()))
		return
	}
	res, svcErr := h.kafkaService.CreateResource(r.Context(), id, &req)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

// PUT /kafka/clusters/:id/resources/:resId
func (h *KafkaHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	resID, err := uuid.Parse(chi.URLParam(r, "resId"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	var req request.UpdateKafkaResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	res, svcErr := h.kafkaService.UpdateResource(r.Context(), id, resID, &req)
	if svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, res)
}

// DELETE /kafka/clusters/:id/resources/:resId
func (h *KafkaHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	resID, err := uuid.Parse(chi.URLParam(r, "resId"))
	if err != nil {
		response.Error(w, apierror.ErrBadRequest)
		return
	}
	if svcErr := h.kafkaService.DeleteResource(r.Context(), id, resID); svcErr != nil {
		h.respondError(w, svcErr)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "resource deleted"})
}
