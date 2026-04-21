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

type KafkaMetrics struct {
	Brokers struct {
		Total           int `json:"total"`
		Online          int `json:"online"`
		UnderReplicated int `json:"under_replicated"`
	} `json:"brokers"`
	Topics struct {
		Total             int `json:"total"`
		Partitions        int `json:"partitions"`
		ReplicationFactor int `json:"replication_factor"`
	} `json:"topics"`
	Consumers struct {
		Total  int   `json:"total"`
		Active int   `json:"active"`
		Lag    int64 `json:"lag"`
	} `json:"consumers"`
}

type KafkaService struct {
	kafkaRepo *repository.KafkaRepository
}

func NewKafkaService(kafkaRepo *repository.KafkaRepository) *KafkaService {
	return &KafkaService{kafkaRepo: kafkaRepo}
}

// ── Clusters ───────────────────────────────────────────────────────────────

func (s *KafkaService) ListClusters(ctx context.Context, clusterID *uuid.UUID, p pagination.Params) ([]model.KafkaCluster, *response.Meta, error) {
	clusters, total, err := s.kafkaRepo.ListClusters(ctx, clusterID, p.PerPage, p.Offset)
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

func (s *KafkaService) GetClusterByID(ctx context.Context, id uuid.UUID) (*model.KafkaCluster, error) {
	kc, err := s.kafkaRepo.FindClusterByID(ctx, id)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if kc == nil {
		return nil, apierror.ErrNotFound
	}
	return kc, nil
}

func (s *KafkaService) CreateCluster(ctx context.Context, req *request.CreateKafkaClusterRequest, createdBy string) (*model.KafkaCluster, error) {
	creatorID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, apierror.ErrBadRequest
	}

	retentionHours := req.LogRetentionHours
	if retentionHours == 0 {
		retentionHours = 168
	}
	replFactor := req.DefaultReplicationFactor
	if replFactor == 0 {
		replFactor = 3
	}
	minInsync := req.MinInsyncReplicas
	if minInsync == 0 {
		minInsync = 2
	}

	brokerCPU := req.Brokers.CPU
	if brokerCPU == "" {
		brokerCPU = "2"
	}
	brokerMem := req.Brokers.Memory
	if brokerMem == "" {
		brokerMem = "4Gi"
	}
	brokerStorage := req.Brokers.Storage
	if brokerStorage == "" {
		brokerStorage = "100Gi"
	}
	storageClass := req.Brokers.StorageClass
	if storageClass == "" {
		storageClass = "standard"
	}
	jvmHeap := req.Brokers.JVMHeap
	if jvmHeap == "" {
		jvmHeap = "2g"
	}

	var authMech *string
	if req.Security.AuthMechanism != "" {
		authMech = &req.Security.AuthMechanism
	}
	var authzType *string
	if req.Security.AuthzType != "" {
		authzType = &req.Security.AuthzType
	}
	tlsVersion := req.Security.TLSVersion
	if tlsVersion == "" {
		tlsVersion = "TLSv1.3"
	}

	srReplicas := req.Services.SchemaRegistryReplicas
	if srReplicas == 0 {
		srReplicas = 1
	}
	kcReplicas := req.Services.KafkaConnectReplicas
	if kcReplicas == 0 {
		kcReplicas = 1
	}
	ksqlReplicas := req.Services.KSQLReplicas
	if ksqlReplicas == 0 {
		ksqlReplicas = 1
	}

	zkReplicas := req.Zookeeper.Replicas
	if zkReplicas == 0 {
		zkReplicas = 3
	}
	zkCPU := req.Zookeeper.CPU
	if zkCPU == "" {
		zkCPU = "0.5"
	}
	zkMem := req.Zookeeper.Memory
	if zkMem == "" {
		zkMem = "1Gi"
	}
	zkStorage := req.Zookeeper.Storage
	if zkStorage == "" {
		zkStorage = "10Gi"
	}

	kc := &model.KafkaCluster{
		ID:                       uuid.New(),
		ClusterID:                req.ClusterID,
		Name:                     req.Name,
		Version:                  req.Version,
		Environment:              req.Environment,
		Namespace:                req.Namespace,
		Status:                   "pending",
		BrokerCount:              req.Brokers.Count,
		BrokerCPU:                brokerCPU,
		BrokerMemory:             brokerMem,
		BrokerStorage:            brokerStorage,
		BrokerStorageClass:       storageClass,
		BrokerJVMHeap:            jvmHeap,
		ZKEnabled:                req.Zookeeper.Enabled,
		ZKReplicas:               zkReplicas,
		ZKCPU:                    zkCPU,
		ZKMemory:                 zkMem,
		ZKStorage:                zkStorage,
		AuthEnabled:              req.Security.AuthEnabled,
		AuthMechanism:            authMech,
		AuthzEnabled:             req.Security.AuthzEnabled,
		AuthzType:                authzType,
		TLSEnabled:               req.Security.TLSEnabled,
		TLSVersion:               tlsVersion,
		ServiceType:              "ClusterIP",
		SchemaRegistryEnabled:    req.Services.SchemaRegistryEnabled,
		SchemaRegistryReplicas:   srReplicas,
		KafkaConnectEnabled:      req.Services.KafkaConnectEnabled,
		KafkaConnectReplicas:     kcReplicas,
		KSQLEnabled:              req.Services.KSQLEnabled,
		KSQLReplicas:             ksqlReplicas,
		MonitoringEnabled:        req.Services.MonitoringEnabled,
		PrometheusEnabled:        req.Services.PrometheusEnabled,
		GrafanaEnabled:           req.Services.GrafanaEnabled,
		JMXExporterEnabled:       true,
		LogRetentionHours:        retentionHours,
		DefaultReplicationFactor: replFactor,
		MinInsyncReplicas:        minInsync,
		CreatedBy:                &creatorID,
	}

	if err := s.kafkaRepo.CreateCluster(ctx, kc); err != nil {
		return nil, apierror.ErrInternal
	}
	return kc, nil
}

func (s *KafkaService) UpdateCluster(ctx context.Context, id uuid.UUID, req *request.UpdateKafkaClusterRequest) (*model.KafkaCluster, error) {
	kc, err := s.GetClusterByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		kc.Name = *req.Name
	}
	if req.Status != nil {
		kc.Status = *req.Status
	}
	if req.Version != nil {
		kc.Version = *req.Version
	}
	if err := s.kafkaRepo.UpdateCluster(ctx, kc); err != nil {
		return nil, apierror.ErrInternal
	}
	return kc, nil
}

func (s *KafkaService) DeleteCluster(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetClusterByID(ctx, id); err != nil {
		return err
	}
	if err := s.kafkaRepo.DeleteCluster(ctx, id); err != nil {
		return apierror.ErrInternal
	}
	return nil
}

func (s *KafkaService) GetMetrics(ctx context.Context, id uuid.UUID) (*KafkaMetrics, error) {
	kc, err := s.GetClusterByID(ctx, id)
	if err != nil {
		return nil, err
	}

	topicCount, _ := s.kafkaRepo.CountTopics(ctx, id)
	consumerCount, _ := s.kafkaRepo.CountConsumerGroups(ctx, id)
	totalLag, _ := s.kafkaRepo.TotalConsumerLag(ctx, id)

	m := &KafkaMetrics{}
	m.Brokers.Total = kc.BrokerCount
	m.Brokers.Online = kc.BrokerCount
	m.Topics.Total = topicCount
	m.Topics.ReplicationFactor = kc.DefaultReplicationFactor
	m.Consumers.Total = consumerCount
	m.Consumers.Active = consumerCount
	m.Consumers.Lag = totalLag

	return m, nil
}

// ── Topics ─────────────────────────────────────────────────────────────────

func (s *KafkaService) ListTopics(ctx context.Context, kafkaClusterID uuid.UUID, p pagination.Params) ([]model.KafkaTopic, *response.Meta, error) {
	if _, err := s.GetClusterByID(ctx, kafkaClusterID); err != nil {
		return nil, nil, err
	}
	topics, total, err := s.kafkaRepo.ListTopics(ctx, kafkaClusterID, p.PerPage, p.Offset)
	if err != nil {
		return nil, nil, apierror.ErrInternal
	}
	meta := &response.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: pagination.TotalPages(total, p.PerPage),
	}
	return topics, meta, nil
}

func (s *KafkaService) CreateTopic(ctx context.Context, kafkaClusterID uuid.UUID, req *request.CreateKafkaTopicRequest) (*model.KafkaTopic, error) {
	if _, err := s.GetClusterByID(ctx, kafkaClusterID); err != nil {
		return nil, err
	}
	retentionMs := req.RetentionMs
	if retentionMs == 0 {
		retentionMs = 604800000
	}
	t := &model.KafkaTopic{
		ID:                uuid.New(),
		KafkaClusterID:    kafkaClusterID,
		Name:              req.Name,
		Partitions:        req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
		RetentionMs:       retentionMs,
		Status:            "active",
	}
	if err := s.kafkaRepo.CreateTopic(ctx, t); err != nil {
		return nil, apierror.ErrInternal
	}
	return t, nil
}

func (s *KafkaService) UpdateTopic(ctx context.Context, kafkaClusterID, topicID uuid.UUID, req *request.UpdateKafkaTopicRequest) (*model.KafkaTopic, error) {
	t, err := s.kafkaRepo.FindTopicByID(ctx, topicID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if t == nil || t.KafkaClusterID != kafkaClusterID {
		return nil, apierror.ErrNotFound
	}
	if req.Partitions != nil {
		t.Partitions = *req.Partitions
	}
	if req.RetentionMs != nil {
		t.RetentionMs = *req.RetentionMs
	}
	if err := s.kafkaRepo.UpdateTopic(ctx, t); err != nil {
		return nil, apierror.ErrInternal
	}
	return t, nil
}

func (s *KafkaService) DeleteTopic(ctx context.Context, kafkaClusterID, topicID uuid.UUID) error {
	t, err := s.kafkaRepo.FindTopicByID(ctx, topicID)
	if err != nil {
		return apierror.ErrInternal
	}
	if t == nil || t.KafkaClusterID != kafkaClusterID {
		return apierror.ErrNotFound
	}
	return s.kafkaRepo.DeleteTopic(ctx, topicID)
}

// ── Consumer Groups ────────────────────────────────────────────────────────

func (s *KafkaService) ListConsumerGroups(ctx context.Context, kafkaClusterID uuid.UUID) ([]model.KafkaConsumerGroup, error) {
	if _, err := s.GetClusterByID(ctx, kafkaClusterID); err != nil {
		return nil, err
	}
	groups, err := s.kafkaRepo.ListConsumerGroups(ctx, kafkaClusterID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	return groups, nil
}

func (s *KafkaService) GetConsumerGroup(ctx context.Context, kafkaClusterID, groupID uuid.UUID) (*model.KafkaConsumerGroup, error) {
	g, err := s.kafkaRepo.FindConsumerGroupByID(ctx, groupID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if g == nil || g.KafkaClusterID != kafkaClusterID {
		return nil, apierror.ErrNotFound
	}
	return g, nil
}

// ── Resources ──────────────────────────────────────────────────────────────

func (s *KafkaService) ListResources(ctx context.Context, kafkaClusterID uuid.UUID) ([]model.KafkaResource, error) {
	if _, err := s.GetClusterByID(ctx, kafkaClusterID); err != nil {
		return nil, err
	}
	return s.kafkaRepo.ListResources(ctx, kafkaClusterID)
}

func (s *KafkaService) CreateResource(ctx context.Context, kafkaClusterID uuid.UUID, req *request.CreateKafkaResourceRequest) (*model.KafkaResource, error) {
	if _, err := s.GetClusterByID(ctx, kafkaClusterID); err != nil {
		return nil, err
	}
	var ns *string
	if req.Namespace != "" {
		ns = &req.Namespace
	}
	res := &model.KafkaResource{
		ID:             uuid.New(),
		KafkaClusterID: kafkaClusterID,
		Name:           req.Name,
		Type:           req.Type,
		Namespace:      ns,
		Configuration:  req.Configuration,
		Status:         "active",
	}
	if err := s.kafkaRepo.CreateResource(ctx, res); err != nil {
		return nil, apierror.ErrInternal
	}
	return res, nil
}

func (s *KafkaService) UpdateResource(ctx context.Context, kafkaClusterID, resID uuid.UUID, req *request.UpdateKafkaResourceRequest) (*model.KafkaResource, error) {
	res, err := s.kafkaRepo.FindResourceByID(ctx, resID)
	if err != nil {
		return nil, apierror.ErrInternal
	}
	if res == nil || res.KafkaClusterID != kafkaClusterID {
		return nil, apierror.ErrNotFound
	}
	if req.Status != nil {
		res.Status = *req.Status
	}
	if req.Configuration != nil {
		res.Configuration = req.Configuration
	}
	if err := s.kafkaRepo.UpdateResource(ctx, res); err != nil {
		return nil, apierror.ErrInternal
	}
	return res, nil
}

func (s *KafkaService) DeleteResource(ctx context.Context, kafkaClusterID, resID uuid.UUID) error {
	res, err := s.kafkaRepo.FindResourceByID(ctx, resID)
	if err != nil {
		return apierror.ErrInternal
	}
	if res == nil || res.KafkaClusterID != kafkaClusterID {
		return apierror.ErrNotFound
	}
	return s.kafkaRepo.DeleteResource(ctx, resID)
}
