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

type ClusterService struct {
	clusterRepo *repository.ClusterRepository
}

func NewClusterService(clusterRepo *repository.ClusterRepository) *ClusterService {
	return &ClusterService{clusterRepo: clusterRepo}
}

func (s *ClusterService) List(ctx context.Context, status, region string, p pagination.Params) ([]model.Cluster, *response.Meta, error) {
	clusters, total, err := s.clusterRepo.List(ctx, status, region, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	meta := &response.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, p.PerPage),
	}
	return clusters, meta, nil
}

func (s *ClusterService) GetByID(ctx context.Context, id uuid.UUID) (*model.Cluster, error) {
	cluster, err := s.clusterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if cluster == nil {
		return nil, apierror.ErrNotFound
	}
	return cluster, nil
}

func (s *ClusterService) Create(ctx context.Context, req *request.CreateClusterRequest, createdBy string) (*model.Cluster, error) {
	creatorID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, apierror.ErrBadRequest
	}

	cluster := &model.Cluster{
		ID:           uuid.New(),
		Name:         req.Name,
		Status:       "healthy",
		Region:       req.Region,
		K8sVersion:   req.K8sVersion,
		APIServerURL: req.APIServerURL,
		Kubeconfig:   req.Kubeconfig,
		CreatedBy:    &creatorID,
	}

	if err := s.clusterRepo.Create(ctx, cluster); err != nil {
		return nil, apierror.ErrInternal
	}
	return cluster, nil
}

func (s *ClusterService) Update(ctx context.Context, id uuid.UUID, req *request.UpdateClusterRequest) (*model.Cluster, error) {
	cluster, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		cluster.Name = *req.Name
	}
	if req.Status != nil {
		cluster.Status = *req.Status
	}
	if req.K8sVersion != nil {
		cluster.K8sVersion = *req.K8sVersion
	}
	if req.APIServerURL != nil {
		cluster.APIServerURL = req.APIServerURL
	}
	if req.Kubeconfig != nil {
		cluster.Kubeconfig = req.Kubeconfig
	}
	if req.NodesTotal != nil {
		cluster.NodesTotal = *req.NodesTotal
	}
	if req.NodesReady != nil {
		cluster.NodesReady = *req.NodesReady
	}

	if err := s.clusterRepo.Update(ctx, cluster); err != nil {
		return nil, apierror.ErrInternal
	}
	return cluster, nil
}

func (s *ClusterService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.clusterRepo.Delete(ctx, id); err != nil {
		return apierror.ErrInternal
	}
	return nil
}

func (s *ClusterService) ListNamespaces(ctx context.Context, clusterID uuid.UUID) ([]model.Namespace, error) {
	if _, err := s.GetByID(ctx, clusterID); err != nil {
		return nil, err
	}
	namespaces, err := s.clusterRepo.ListNamespaces(ctx, clusterID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return namespaces, nil
}

func (s *ClusterService) ListResources(ctx context.Context, clusterID, namespaceID uuid.UUID) ([]model.Resource, error) {
	if _, err := s.GetByID(ctx, clusterID); err != nil {
		return nil, err
	}
	resources, err := s.clusterRepo.ListResources(ctx, namespaceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return resources, nil
}
