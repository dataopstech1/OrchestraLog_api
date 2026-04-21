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

type BigDataService struct {
	repo *repository.BigDataRepository
}

func NewBigDataService(repo *repository.BigDataRepository) *BigDataService {
	return &BigDataService{repo: repo}
}

// ── Spark ──────────────────────────────────────────────────────────────────

func (s *BigDataService) ListSparkClusters(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.SparkCluster, *response.Meta, error) {
	items, total, err := s.repo.ListSparkClusters(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BigDataService) GetSparkCluster(ctx context.Context, id uuid.UUID) (*model.SparkCluster, error) {
	item, err := s.repo.FindSparkCluster(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BigDataService) CreateSparkCluster(ctx context.Context, req *request.CreateSparkClusterRequest) (*model.SparkCluster, error) {
	workerCount := req.WorkerCount
	if workerCount == 0 {
		workerCount = 4
	}
	cpu := req.WorkerCPU
	if cpu == "" {
		cpu = "4"
	}
	mem := req.WorkerMemory
	if mem == "" {
		mem = "8Gi"
	}
	sc := &model.SparkCluster{
		ID:           uuid.New(),
		ClusterID:    req.ClusterID,
		Name:         req.Name,
		Version:      req.Version,
		Namespace:    req.Namespace,
		Status:       "pending",
		WorkerCount:  workerCount,
		WorkerCPU:    cpu,
		WorkerMemory: mem,
	}
	if err := s.repo.CreateSparkCluster(ctx, sc); err != nil {
		return nil, apierror.ErrInternal
	}
	return sc, nil
}

func (s *BigDataService) UpdateSparkCluster(ctx context.Context, id uuid.UUID, req *request.UpdateSparkClusterRequest) (*model.SparkCluster, error) {
	sc, err := s.GetSparkCluster(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		sc.Name = *req.Name
	}
	if req.Status != nil {
		sc.Status = *req.Status
	}
	if req.WorkerCount != nil {
		sc.WorkerCount = *req.WorkerCount
	}
	if err := s.repo.UpdateSparkCluster(ctx, sc); err != nil {
		return nil, apierror.ErrInternal
	}
	return sc, nil
}

func (s *BigDataService) DeleteSparkCluster(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetSparkCluster(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteSparkCluster(ctx, id)
}

func (s *BigDataService) ListSparkApplications(ctx context.Context, sparkClusterID uuid.UUID) ([]model.SparkApplication, error) {
	if _, err := s.GetSparkCluster(ctx, sparkClusterID); err != nil {
		return nil, err
	}
	apps, err := s.repo.ListSparkApplications(ctx, sparkClusterID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return apps, nil
}

func (s *BigDataService) GetSparkApplication(ctx context.Context, sparkClusterID, appID uuid.UUID) (*model.SparkApplication, error) {
	app, err := s.repo.FindSparkApplication(ctx, appID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if app == nil || app.SparkClusterID != sparkClusterID {
		return nil, apierror.ErrNotFound
	}
	return app, nil
}

// ── Flink ──────────────────────────────────────────────────────────────────

func (s *BigDataService) ListFlinkClusters(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.FlinkCluster, *response.Meta, error) {
	items, total, err := s.repo.ListFlinkClusters(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BigDataService) GetFlinkCluster(ctx context.Context, id uuid.UUID) (*model.FlinkCluster, error) {
	item, err := s.repo.FindFlinkCluster(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BigDataService) CreateFlinkCluster(ctx context.Context, req *request.CreateFlinkClusterRequest) (*model.FlinkCluster, error) {
	tmCount := req.TaskManagerCount
	if tmCount == 0 {
		tmCount = 4
	}
	slots := req.SlotsPerTM
	if slots == 0 {
		slots = 4
	}
	fc := &model.FlinkCluster{
		ID:               uuid.New(),
		ClusterID:        req.ClusterID,
		Name:             req.Name,
		Version:          req.Version,
		Namespace:        req.Namespace,
		Status:           "pending",
		TaskManagerCount: tmCount,
		SlotsPerTM:       slots,
	}
	if err := s.repo.CreateFlinkCluster(ctx, fc); err != nil {
		return nil, apierror.ErrInternal
	}
	return fc, nil
}

func (s *BigDataService) UpdateFlinkCluster(ctx context.Context, id uuid.UUID, req *request.UpdateFlinkClusterRequest) (*model.FlinkCluster, error) {
	fc, err := s.GetFlinkCluster(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		fc.Name = *req.Name
	}
	if req.Status != nil {
		fc.Status = *req.Status
	}
	if err := s.repo.UpdateFlinkCluster(ctx, fc); err != nil {
		return nil, apierror.ErrInternal
	}
	return fc, nil
}

func (s *BigDataService) DeleteFlinkCluster(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetFlinkCluster(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteFlinkCluster(ctx, id)
}

func (s *BigDataService) ListFlinkJobs(ctx context.Context, flinkClusterID uuid.UUID) ([]model.FlinkJob, error) {
	if _, err := s.GetFlinkCluster(ctx, flinkClusterID); err != nil {
		return nil, err
	}
	jobs, err := s.repo.ListFlinkJobs(ctx, flinkClusterID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return jobs, nil
}

func (s *BigDataService) GetFlinkJob(ctx context.Context, flinkClusterID, jobID uuid.UUID) (*model.FlinkJob, error) {
	job, err := s.repo.FindFlinkJob(ctx, jobID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if job == nil || job.FlinkClusterID != flinkClusterID {
		return nil, apierror.ErrNotFound
	}
	return job, nil
}

func (s *BigDataService) CancelFlinkJob(ctx context.Context, flinkClusterID, jobID uuid.UUID) error {
	if _, err := s.GetFlinkJob(ctx, flinkClusterID, jobID); err != nil {
		return err
	}
	return s.repo.UpdateFlinkJobStatus(ctx, jobID, "Canceled")
}

// ── Hive ───────────────────────────────────────────────────────────────────

func (s *BigDataService) ListHiveInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.HiveInstance, *response.Meta, error) {
	items, total, err := s.repo.ListHiveInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BigDataService) GetHiveInstance(ctx context.Context, id uuid.UUID) (*model.HiveInstance, error) {
	item, err := s.repo.FindHiveInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BigDataService) CreateHiveInstance(ctx context.Context, req *request.CreateHiveInstanceRequest) (*model.HiveInstance, error) {
	hi := &model.HiveInstance{
		ID:             uuid.New(),
		ClusterID:      req.ClusterID,
		Name:           req.Name,
		Version:        req.Version,
		Namespace:      req.Namespace,
		Status:         "pending",
		MetastoreURL:   req.MetastoreURL,
		HiveServer2URL: req.HiveServer2URL,
	}
	if err := s.repo.CreateHiveInstance(ctx, hi); err != nil {
		return nil, apierror.ErrInternal
	}
	return hi, nil
}

func (s *BigDataService) ListHiveTables(ctx context.Context, hiveInstanceID uuid.UUID) ([]model.HiveTable, error) {
	if _, err := s.GetHiveInstance(ctx, hiveInstanceID); err != nil {
		return nil, err
	}
	tables, err := s.repo.ListHiveTables(ctx, hiveInstanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return tables, nil
}

func (s *BigDataService) GetHiveTable(ctx context.Context, hiveInstanceID, tableID uuid.UUID) (*model.HiveTable, error) {
	t, err := s.repo.FindHiveTable(ctx, tableID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if t == nil || t.HiveInstanceID != hiveInstanceID {
		return nil, apierror.ErrNotFound
	}
	return t, nil
}

// ── HDFS ───────────────────────────────────────────────────────────────────

func (s *BigDataService) ListHDFSClusters(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.HDFSCluster, *response.Meta, error) {
	items, total, err := s.repo.ListHDFSClusters(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BigDataService) GetHDFSCluster(ctx context.Context, id uuid.UUID) (*model.HDFSCluster, error) {
	item, err := s.repo.FindHDFSCluster(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BigDataService) CreateHDFSCluster(ctx context.Context, req *request.CreateHDFSClusterRequest) (*model.HDFSCluster, error) {
	nn := req.NamenodeCount
	if nn == 0 {
		nn = 2
	}
	dn := req.DatanodeCount
	if dn == 0 {
		dn = 6
	}
	jn := req.JournalnodeCount
	if jn == 0 {
		jn = 3
	}
	hc := &model.HDFSCluster{
		ID:               uuid.New(),
		ClusterID:        req.ClusterID,
		Name:             req.Name,
		Version:          req.Version,
		Namespace:        req.Namespace,
		Status:           "pending",
		NamenodeCount:    nn,
		DatanodeCount:    dn,
		JournalnodeCount: jn,
	}
	if err := s.repo.CreateHDFSCluster(ctx, hc); err != nil {
		return nil, apierror.ErrInternal
	}
	return hc, nil
}

// ── NiFi ───────────────────────────────────────────────────────────────────

func (s *BigDataService) ListNiFiInstances(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.NiFiInstance, *response.Meta, error) {
	items, total, err := s.repo.ListNiFiInstances(ctx, clusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	return items, buildMeta(p, total), nil
}

func (s *BigDataService) GetNiFiInstance(ctx context.Context, id uuid.UUID) (*model.NiFiInstance, error) {
	item, err := s.repo.FindNiFiInstance(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if item == nil {
		return nil, apierror.ErrNotFound
	}
	return item, nil
}

func (s *BigDataService) CreateNiFiInstance(ctx context.Context, req *request.CreateNiFiInstanceRequest) (*model.NiFiInstance, error) {
	nodeCount := req.NodeCount
	if nodeCount == 0 {
		nodeCount = 3
	}
	ni := &model.NiFiInstance{
		ID:        uuid.New(),
		ClusterID: req.ClusterID,
		Name:      req.Name,
		Version:   req.Version,
		Namespace: req.Namespace,
		Status:    "pending",
		NodeCount: nodeCount,
	}
	if err := s.repo.CreateNiFiInstance(ctx, ni); err != nil {
		return nil, apierror.ErrInternal
	}
	return ni, nil
}

func (s *BigDataService) ListNiFiProcessGroups(ctx context.Context, nifiInstanceID uuid.UUID) ([]model.NiFiProcessGroup, error) {
	if _, err := s.GetNiFiInstance(ctx, nifiInstanceID); err != nil {
		return nil, err
	}
	groups, err := s.repo.ListNiFiProcessGroups(ctx, nifiInstanceID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return groups, nil
}

// ── Helpers ────────────────────────────────────────────────────────────────

func buildMeta(p pagination.Params, total int) *response.Meta {
	return &response.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, p.PerPage),
	}
}
