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

type BIService struct {
	repo *repository.BIRepository
}

func NewBIService(repo *repository.BIRepository) *BIService {
	return &BIService{repo: repo}
}

// ── Superset ───────────────────────────────────────────────────────────────

func (s *BIService) ListSupersetInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.SupersetInstance, *response.Meta, error) {
	items, total, err := s.repo.ListSupersetInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BIService) GetSupersetInstance(ctx context.Context, id uuid.UUID) (*model.SupersetInstance, error) {
	item, err := s.repo.FindSupersetInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BIService) CreateSupersetInstance(ctx context.Context, req *request.CreateSupersetInstanceRequest) (*model.SupersetInstance, error) {
	inst := &model.SupersetInstance{
		ID:        uuid.New(),
		ClusterID: req.ClusterID,
		Name:      req.Name,
		Version:   req.Version,
		Namespace: req.Namespace,
		Status:    "pending",
		URL:       req.URL,
	}
	if err := s.repo.CreateSupersetInstance(ctx, inst); err != nil {
		return nil, apierror.ErrInternal
	}
	return inst, nil
}

func (s *BIService) ListSupersetDashboards(ctx context.Context, instanceID uuid.UUID) ([]model.SupersetDashboard, error) {
	if _, err := s.GetSupersetInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListSupersetDashboards(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

// ── Metabase ───────────────────────────────────────────────────────────────

func (s *BIService) ListMetabaseInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.MetabaseInstance, *response.Meta, error) {
	items, total, err := s.repo.ListMetabaseInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BIService) GetMetabaseInstance(ctx context.Context, id uuid.UUID) (*model.MetabaseInstance, error) {
	item, err := s.repo.FindMetabaseInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BIService) CreateMetabaseInstance(ctx context.Context, req *request.CreateMetabaseInstanceRequest) (*model.MetabaseInstance, error) {
	inst := &model.MetabaseInstance{
		ID:        uuid.New(),
		ClusterID: req.ClusterID,
		Name:      req.Name,
		Version:   req.Version,
		Namespace: req.Namespace,
		Status:    "pending",
		URL:       req.URL,
	}
	if err := s.repo.CreateMetabaseInstance(ctx, inst); err != nil {
		return nil, apierror.ErrInternal
	}
	return inst, nil
}

func (s *BIService) ListMetabaseDashboards(ctx context.Context, instanceID uuid.UUID) ([]model.MetabaseDashboard, error) {
	if _, err := s.GetMetabaseInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListMetabaseDashboards(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

// ── N8N ───────────────────────────────────────────────────────────────────

func (s *BIService) ListN8NInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.N8NInstance, *response.Meta, error) {
	items, total, err := s.repo.ListN8NInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BIService) GetN8NInstance(ctx context.Context, id uuid.UUID) (*model.N8NInstance, error) {
	item, err := s.repo.FindN8NInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BIService) CreateN8NInstance(ctx context.Context, req *request.CreateN8NInstanceRequest) (*model.N8NInstance, error) {
	inst := &model.N8NInstance{
		ID:        uuid.New(),
		ClusterID: req.ClusterID,
		Name:      req.Name,
		Version:   req.Version,
		Namespace: req.Namespace,
		Status:    "pending",
		URL:       req.URL,
	}
	if err := s.repo.CreateN8NInstance(ctx, inst); err != nil {
		return nil, apierror.ErrInternal
	}
	return inst, nil
}

func (s *BIService) ListN8NWorkflows(ctx context.Context, instanceID uuid.UUID) ([]model.N8NWorkflow, error) {
	if _, err := s.GetN8NInstance(ctx, instanceID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListN8NWorkflows(ctx, instanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}

func (s *BIService) GetN8NWorkflow(ctx context.Context, instanceID, wfID uuid.UUID) (*model.N8NWorkflow, error) {
	wf, err := s.repo.FindN8NWorkflow(ctx, wfID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if wf == nil || wf.N8NInstanceID != instanceID {
		return nil, apierror.ErrNotFound
	}
	return wf, nil
}

func (s *BIService) SetN8NWorkflowStatus(ctx context.Context, instanceID, wfID uuid.UUID, status string) error {
	if _, err := s.GetN8NWorkflow(ctx, instanceID, wfID); err != nil {
		return err
	}
	return s.repo.UpdateN8NWorkflowStatus(ctx, wfID, status)
}
