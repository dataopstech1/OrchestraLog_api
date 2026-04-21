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

type DataFlowService struct {
	repo *repository.DataFlowRepository
}

func NewDataFlowService(repo *repository.DataFlowRepository) *DataFlowService {
	return &DataFlowService{repo: repo}
}

func (s *DataFlowService) List(ctx context.Context, status string, p pagination.Params) ([]model.DataFlow, *response.Meta, error) {
	items, total, err := s.repo.List(ctx, status, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *DataFlowService) GetByID(ctx context.Context, id uuid.UUID) (*model.DataFlowWithGraph, error) {
	df, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if df == nil {
		return nil, apierror.ErrNotFound
	}

	nodes, _ := s.repo.ListNodes(ctx, id)
	edges, _ := s.repo.ListEdges(ctx, id)

	return &model.DataFlowWithGraph{DataFlow: *df, Nodes: nodes, Edges: edges}, nil
}

func (s *DataFlowService) Create(ctx context.Context, req *request.CreateDataFlowRequest, createdBy string) (*model.DataFlowWithGraph, error) {
	creatorID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, apierror.ErrBadRequest
	}

	flowID := uuid.New()
	df := &model.DataFlow{
		ID:          flowID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "draft",
		Layout:      "auto",
		CreatedBy:   &creatorID,
	}

	if err := s.repo.Create(ctx, df); err != nil {
		return nil, apierror.ErrInternal
	}

	nodes := make([]model.FlowNode, 0, len(req.Nodes))
	for _, n := range req.Nodes {
		nodes = append(nodes, model.FlowNode{
			ID:           uuid.New(),
			DataFlowID:   flowID,
			NodeID:       n.NodeID,
			Label:        n.Label,
			ResourceType: n.ResourceType,
			Cluster:      n.Cluster,
			Namespace:    n.Namespace,
			Status:       "healthy",
			PositionX:    n.PositionX,
			PositionY:    n.PositionY,
			Config:       n.Config,
		})
	}

	edges := make([]model.FlowEdge, 0, len(req.Edges))
	for _, e := range req.Edges {
		edges = append(edges, model.FlowEdge{
			ID:           uuid.New(),
			DataFlowID:   flowID,
			EdgeID:       e.EdgeID,
			SourceNodeID: e.SourceNodeID,
			TargetNodeID: e.TargetNodeID,
			FlowType:     e.FlowType,
			Label:        e.Label,
			Animated:     e.Animated,
			Config:       e.Config,
		})
	}

	if len(nodes) > 0 {
		if err := s.repo.ReplaceNodes(ctx, flowID, nodes); err != nil {
			return nil, apierror.ErrInternal
		}
	}
	if len(edges) > 0 {
		if err := s.repo.ReplaceEdges(ctx, flowID, edges); err != nil {
			return nil, apierror.ErrInternal
		}
	}

	return &model.DataFlowWithGraph{DataFlow: *df, Nodes: nodes, Edges: edges}, nil
}

func (s *DataFlowService) Update(ctx context.Context, id uuid.UUID, req *request.UpdateDataFlowRequest) (*model.DataFlow, error) {
	df, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if df == nil {
		return nil, apierror.ErrNotFound
	}
	if req.Name != nil {
		df.Name = *req.Name
	}
	if req.Description != nil {
		df.Description = req.Description
	}
	if req.Status != nil {
		df.Status = *req.Status
	}
	if err := s.repo.Update(ctx, df); err != nil {
		return nil, apierror.ErrInternal
	}
	return df, nil
}

func (s *DataFlowService) Delete(ctx context.Context, id uuid.UUID) error {
	df, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apierror.ErrInternal
	}
	if df == nil {
		return apierror.ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *DataFlowService) UpdateNodes(ctx context.Context, id uuid.UUID, req *request.UpdateFlowNodesRequest) ([]model.FlowNode, error) {
	df, err := s.repo.FindByID(ctx, id)
	if err != nil || df == nil {
		return nil, apierror.ErrNotFound
	}
	nodes := make([]model.FlowNode, 0, len(req.Nodes))
	for _, n := range req.Nodes {
		nodes = append(nodes, model.FlowNode{
			ID:           uuid.New(),
			DataFlowID:   id,
			NodeID:       n.NodeID,
			Label:        n.Label,
			ResourceType: n.ResourceType,
			Cluster:      n.Cluster,
			Namespace:    n.Namespace,
			Status:       "healthy",
			PositionX:    n.PositionX,
			PositionY:    n.PositionY,
			Config:       n.Config,
		})
	}
	if err := s.repo.ReplaceNodes(ctx, id, nodes); err != nil {
		return nil, apierror.ErrInternal
	}
	return nodes, nil
}

func (s *DataFlowService) UpdateEdges(ctx context.Context, id uuid.UUID, req *request.UpdateFlowEdgesRequest) ([]model.FlowEdge, error) {
	df, err := s.repo.FindByID(ctx, id)
	if err != nil || df == nil {
		return nil, apierror.ErrNotFound
	}
	edges := make([]model.FlowEdge, 0, len(req.Edges))
	for _, e := range req.Edges {
		edges = append(edges, model.FlowEdge{
			ID:           uuid.New(),
			DataFlowID:   id,
			EdgeID:       e.EdgeID,
			SourceNodeID: e.SourceNodeID,
			TargetNodeID: e.TargetNodeID,
			FlowType:     e.FlowType,
			Label:        e.Label,
			Animated:     e.Animated,
			Config:       e.Config,
		})
	}
	if err := s.repo.ReplaceEdges(ctx, id, edges); err != nil {
		return nil, apierror.ErrInternal
	}
	return edges, nil
}

func (s *DataFlowService) SetStatus(ctx context.Context, id uuid.UUID, status string) (*model.DataFlow, error) {
	return s.Update(ctx, id, &request.UpdateDataFlowRequest{Status: &status})
}

func (s *DataFlowService) ListTemplates(ctx context.Context) ([]model.FlowTemplate, error) {
	items, err := s.repo.ListTemplates(ctx)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return items, nil
}
