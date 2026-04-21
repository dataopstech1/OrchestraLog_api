package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orchestralog/api/internal/dto/request"
	"github.com/orchestralog/api/internal/model"
	"github.com/orchestralog/api/internal/repository"
	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/pagination"
	"github.com/orchestralog/api/pkg/response"
)

type MLAIService struct {
	repo *repository.MLAIRepository
}

func NewMLAIService(repo *repository.MLAIRepository) *MLAIService {
	return &MLAIService{repo: repo}
}

// ── MLFlow ─────────────────────────────────────────────────────────────────

func (s *MLAIService) ListMLFlowInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.MLFlowInstance, *response.Meta, error) {
	items, total, err := s.repo.ListMLFlowInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *MLAIService) GetMLFlowInstance(ctx context.Context, id uuid.UUID) (*model.MLFlowInstance, error) {
	item, err := s.repo.FindMLFlowInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *MLAIService) CreateMLFlowInstance(ctx context.Context, req *request.CreateMLFlowInstanceRequest) (*model.MLFlowInstance, error) {
	m := &model.MLFlowInstance{
		ID:            uuid.New(),
		ClusterID:     req.ClusterID,
		Name:          req.Name,
		Version:       req.Version,
		Namespace:     req.Namespace,
		Status:        "pending",
		TrackingURL:   req.TrackingURL,
		ArtifactStore: req.ArtifactStore,
	}
	if err := s.repo.CreateMLFlowInstance(ctx, m); err != nil {
		return nil, apierror.ErrInternal
	}
	return m, nil
}

func (s *MLAIService) ListExperiments(ctx context.Context, instanceID uuid.UUID) ([]model.MLFlowExperiment, error) {
	if _, err := s.GetMLFlowInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListExperiments(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

func (s *MLAIService) GetExperiment(ctx context.Context, instanceID, expID uuid.UUID) (*model.MLFlowExperiment, error) {
	exp, err := s.repo.FindExperiment(ctx, expID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if exp == nil || exp.MLFlowInstanceID != instanceID {
		return nil, apierror.ErrNotFound
	}
	return exp, nil
}

func (s *MLAIService) ListModels(ctx context.Context, instanceID uuid.UUID) ([]model.MLFlowModel, error) {
	if _, err := s.GetMLFlowInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListModels(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

func (s *MLAIService) GetModel(ctx context.Context, instanceID, modelID uuid.UUID) (*model.MLFlowModel, error) {
	m, err := s.repo.FindModel(ctx, modelID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if m == nil || m.MLFlowInstanceID != instanceID {
		return nil, apierror.ErrNotFound
	}
	return m, nil
}

func (s *MLAIService) UpdateModelStage(ctx context.Context, instanceID, modelID uuid.UUID, req *request.UpdateMLFlowModelStageRequest) (*model.MLFlowModel, error) {
	m, err := s.GetModel(ctx, instanceID, modelID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateModelStage(ctx, modelID, req.Stage); err != nil {
		return nil, apierror.ErrInternal
	}
	m.Stage = &req.Stage
	return m, nil
}

// ── Feast ──────────────────────────────────────────────────────────────────

func (s *MLAIService) ListFeastInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.FeastInstance, *response.Meta, error) {
	items, total, err := s.repo.ListFeastInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *MLAIService) GetFeastInstance(ctx context.Context, id uuid.UUID) (*model.FeastInstance, error) {
	item, err := s.repo.FindFeastInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *MLAIService) CreateFeastInstance(ctx context.Context, req *request.CreateFeastInstanceRequest) (*model.FeastInstance, error) {
	fi := &model.FeastInstance{
		ID:           uuid.New(),
		ClusterID:    req.ClusterID,
		Name:         req.Name,
		Version:      req.Version,
		Namespace:    req.Namespace,
		Status:       "pending",
		OnlineStore:  req.OnlineStore,
		OfflineStore: req.OfflineStore,
	}
	if err := s.repo.CreateFeastInstance(ctx, fi); err != nil {
		return nil, apierror.ErrInternal
	}
	return fi, nil
}

func (s *MLAIService) ListFeastEntities(ctx context.Context, instanceID uuid.UUID) ([]model.FeastEntity, error) {
	if _, err := s.GetFeastInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListFeastEntities(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

func (s *MLAIService) ListFeastFeatureViews(ctx context.Context, instanceID uuid.UUID) ([]model.FeastFeatureView, error) {
	if _, err := s.GetFeastInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListFeastFeatureViews(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

// ── JupyterHub ─────────────────────────────────────────────────────────────

func (s *MLAIService) ListJupyterHubInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.JupyterHubInstance, *response.Meta, error) {
	items, total, err := s.repo.ListJupyterHubInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *MLAIService) GetJupyterHubInstance(ctx context.Context, id uuid.UUID) (*model.JupyterHubInstance, error) {
	item, err := s.repo.FindJupyterHubInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *MLAIService) CreateJupyterHubInstance(ctx context.Context, req *request.CreateJupyterHubInstanceRequest) (*model.JupyterHubInstance, error) {
	maxUsers := req.MaxUsers
	if maxUsers == 0 {
		maxUsers = 100
	}
	jh := &model.JupyterHubInstance{
		ID:        uuid.New(),
		ClusterID: req.ClusterID,
		Name:      req.Name,
		Version:   req.Version,
		Namespace: req.Namespace,
		Status:    "pending",
		HubURL:    req.HubURL,
		MaxUsers:  maxUsers,
	}
	if err := s.repo.CreateJupyterHubInstance(ctx, jh); err != nil {
		return nil, apierror.ErrInternal
	}
	return jh, nil
}

func (s *MLAIService) ListNotebooks(ctx context.Context, instanceID uuid.UUID) ([]model.JupyterHubNotebook, error) {
	if _, err := s.GetJupyterHubInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListNotebooks(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

func (s *MLAIService) CreateNotebook(ctx context.Context, instanceID uuid.UUID, req *request.CreateNotebookRequest) (*model.JupyterHubNotebook, error) {
	if _, err := s.GetJupyterHubInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	nb := &model.JupyterHubNotebook{
		ID:                   uuid.New(),
		JupyterHubInstanceID: instanceID,
		Name:                 req.Name,
		Owner:                req.Owner,
		Status:               "Stopped",
		Image:                req.Image,
		CPULimit:             req.CPULimit,
		MemoryLimit:          req.MemoryLimit,
		GPULimit:             req.GPULimit,
	}
	if err := s.repo.CreateNotebook(ctx, nb); err != nil {
		return nil, apierror.ErrInternal
	}
	return nb, nil
}

func (s *MLAIService) UpdateNotebook(ctx context.Context, instanceID, nbID uuid.UUID, req *request.UpdateNotebookRequest) (*model.JupyterHubNotebook, error) {
	nb, err := s.repo.FindNotebook(ctx, nbID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if nb == nil || nb.JupyterHubInstanceID != instanceID {
		return nil, apierror.ErrNotFound
	}
	if req.Image != nil {
		nb.Image = req.Image
	}
	if req.CPULimit != nil {
		nb.CPULimit = req.CPULimit
	}
	if req.MemoryLimit != nil {
		nb.MemoryLimit = req.MemoryLimit
	}
	if req.GPULimit != nil {
		nb.GPULimit = *req.GPULimit
	}
	if err := s.repo.UpdateNotebook(ctx, nb); err != nil {
		return nil, apierror.ErrInternal
	}
	return nb, nil
}

func (s *MLAIService) DeleteNotebook(ctx context.Context, instanceID, nbID uuid.UUID) error {
	nb, err := s.repo.FindNotebook(ctx, nbID)
	if err != nil {
		return apierror.ErrInternal
	}
	if nb == nil || nb.JupyterHubInstanceID != instanceID {
		return apierror.ErrNotFound
	}
	return s.repo.DeleteNotebook(ctx, nbID)
}

func (s *MLAIService) SetNotebookStatus(ctx context.Context, instanceID, nbID uuid.UUID, status string) error {
	nb, err := s.repo.FindNotebook(ctx, nbID)
	if err != nil {
		return apierror.ErrInternal
	}
	if nb == nil || nb.JupyterHubInstanceID != instanceID {
		return apierror.ErrNotFound
	}
	return s.repo.UpdateNotebookStatus(ctx, nbID, status)
}

// ── LLM ───────────────────────────────────────────────────────────────────

func (s *MLAIService) ListLLMDeployments(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.LLMDeployment, *response.Meta, error) {
	items, total, err := s.repo.ListLLMDeployments(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *MLAIService) GetLLMDeployment(ctx context.Context, id uuid.UUID) (*model.LLMDeployment, error) {
	item, err := s.repo.FindLLMDeployment(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *MLAIService) CreateLLMDeployment(ctx context.Context, req *request.CreateLLMDeploymentRequest) (*model.LLMDeployment, error) {
	replicas := req.Replicas
	if replicas == 0 {
		replicas = 1
	}
	d := &model.LLMDeployment{
		ID:            uuid.New(),
		ClusterID:     req.ClusterID,
		Name:          req.Name,
		ModelName:     req.ModelName,
		ModelVersion:  req.ModelVersion,
		Namespace:     req.Namespace,
		Status:        "pending",
		Replicas:      replicas,
		GPUCount:      req.GPUCount,
		GPUType:       req.GPUType,
		MaxTokens:     req.MaxTokens,
		ContextWindow: req.ContextWindow,
	}
	if err := s.repo.CreateLLMDeployment(ctx, d); err != nil {
		return nil, apierror.ErrInternal
	}
	return d, nil
}

func (s *MLAIService) UpdateLLMDeployment(ctx context.Context, id uuid.UUID, req *request.UpdateLLMDeploymentRequest) (*model.LLMDeployment, error) {
	d, err := s.GetLLMDeployment(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		d.Name = *req.Name
	}
	if req.Status != nil {
		d.Status = *req.Status
	}
	if err := s.repo.UpdateLLMDeployment(ctx, d); err != nil {
		return nil, apierror.ErrInternal
	}
	return d, nil
}

func (s *MLAIService) DeleteLLMDeployment(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetLLMDeployment(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteLLMDeployment(ctx, id)
}

func (s *MLAIService) ScaleLLMDeployment(ctx context.Context, id uuid.UUID, req *request.ScaleLLMDeploymentRequest) (*model.LLMDeployment, error) {
	d, err := s.GetLLMDeployment(ctx, id)
	if err != nil {
		return nil, err
	}
	d.Replicas = req.Replicas
	if err := s.repo.UpdateLLMDeployment(ctx, d); err != nil {
		return nil, apierror.ErrInternal
	}
	return d, nil
}
