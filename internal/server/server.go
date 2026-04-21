package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/orchestralog/api/internal/config"
	"github.com/orchestralog/api/internal/handler"
	"github.com/orchestralog/api/internal/middleware"
	"github.com/orchestralog/api/internal/repository"
	"github.com/orchestralog/api/internal/service"
)

type Server struct {
	cfg    *config.Config
	db     *sqlx.DB
	router *chi.Mux
}

func New(cfg *config.Config, db *sqlx.DB) *Server {
	s := &Server{cfg: cfg, db: db}
	s.router = s.buildRouter()
	return s
}

func (s *Server) buildRouter() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.RateLimit(200, time.Minute))

	// Repositories
	userRepo := repository.NewUserRepository(s.db)
	clusterRepo := repository.NewClusterRepository(s.db)
	kafkaRepo := repository.NewKafkaRepository(s.db)
	bigDataRepo := repository.NewBigDataRepository(s.db)
	mlaiRepo := repository.NewMLAIRepository(s.db)
	biRepo := repository.NewBIRepository(s.db)
	dataFlowRepo := repository.NewDataFlowRepository(s.db)

	// Services
	authService := service.NewAuthService(userRepo, s.cfg)
	clusterService := service.NewClusterService(clusterRepo)
	kafkaService := service.NewKafkaService(kafkaRepo)
	bigDataService := service.NewBigDataService(bigDataRepo)
	mlaiService := service.NewMLAIService(mlaiRepo)
	biService := service.NewBIService(biRepo)
	dataFlowService := service.NewDataFlowService(dataFlowRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	clusterHandler := handler.NewClusterHandler(clusterService)
	kafkaHandler := handler.NewKafkaHandler(kafkaService)
	bigDataHandler := handler.NewBigDataHandler(bigDataService)
	mlaiHandler := handler.NewMLAIHandler(mlaiService)
	biHandler := handler.NewBIHandler(biService)
	dataFlowHandler := handler.NewDataFlowHandler(dataFlowService)
	monitoringHandler := handler.NewMonitoringHandler()
	dashboardHandler := handler.NewDashboardHandler(s.db)
	userHandler := handler.NewUserHandler(userService)

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)

			// Protected auth routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(s.cfg.JWT.AccessSecret))
				r.Post("/logout", authHandler.Logout)
				r.Get("/me", authHandler.Me)
			})
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.cfg.JWT.AccessSecret))

			// Kafka
			r.Route("/kafka/clusters", func(r chi.Router) {
				r.Get("/", kafkaHandler.ListClusters)
				r.With(middleware.RequireAdminOrOperator()).Post("/", kafkaHandler.CreateCluster)
				r.Get("/{id}", kafkaHandler.GetCluster)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}", kafkaHandler.UpdateCluster)
				r.With(middleware.RequireAdmin()).Delete("/{id}", kafkaHandler.DeleteCluster)
				r.Get("/{id}/metrics", kafkaHandler.GetMetrics)
				r.Get("/{id}/topics", kafkaHandler.ListTopics)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/topics", kafkaHandler.CreateTopic)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}/topics/{topicId}", kafkaHandler.UpdateTopic)
				r.With(middleware.RequireAdmin()).Delete("/{id}/topics/{topicId}", kafkaHandler.DeleteTopic)
				r.Get("/{id}/consumers", kafkaHandler.ListConsumers)
				r.Get("/{id}/consumers/{groupId}", kafkaHandler.GetConsumer)
				r.Get("/{id}/resources", kafkaHandler.ListResources)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/resources", kafkaHandler.CreateResource)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}/resources/{resId}", kafkaHandler.UpdateResource)
				r.With(middleware.RequireAdmin()).Delete("/{id}/resources/{resId}", kafkaHandler.DeleteResource)
			})

			// Spark
			r.Route("/spark/clusters", func(r chi.Router) {
				r.Get("/", bigDataHandler.ListSparkClusters)
				r.With(middleware.RequireAdminOrOperator()).Post("/", bigDataHandler.CreateSparkCluster)
				r.Get("/{id}", bigDataHandler.GetSparkCluster)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}", bigDataHandler.UpdateSparkCluster)
				r.With(middleware.RequireAdmin()).Delete("/{id}", bigDataHandler.DeleteSparkCluster)
				r.Get("/{id}/applications", bigDataHandler.ListSparkApplications)
				r.Get("/{id}/applications/{appId}", bigDataHandler.GetSparkApplication)
				r.Get("/{id}/metrics", bigDataHandler.MetricsStub)
			})

			// Flink
			r.Route("/flink/clusters", func(r chi.Router) {
				r.Get("/", bigDataHandler.ListFlinkClusters)
				r.With(middleware.RequireAdminOrOperator()).Post("/", bigDataHandler.CreateFlinkCluster)
				r.Get("/{id}", bigDataHandler.GetFlinkCluster)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}", bigDataHandler.UpdateFlinkCluster)
				r.With(middleware.RequireAdmin()).Delete("/{id}", bigDataHandler.DeleteFlinkCluster)
				r.Get("/{id}/jobs", bigDataHandler.ListFlinkJobs)
				r.Get("/{id}/jobs/{jobId}", bigDataHandler.GetFlinkJob)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/jobs/{jobId}/cancel", bigDataHandler.CancelFlinkJob)
				r.Get("/{id}/metrics", bigDataHandler.MetricsStub)
			})

			// Hive
			r.Route("/hive/instances", func(r chi.Router) {
				r.Get("/", bigDataHandler.ListHiveInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", bigDataHandler.CreateHiveInstance)
				r.Get("/{id}", bigDataHandler.GetHiveInstance)
				r.Get("/{id}/tables", bigDataHandler.ListHiveTables)
				r.Get("/{id}/tables/{tableId}", bigDataHandler.GetHiveTable)
				r.Get("/{id}/metrics", bigDataHandler.MetricsStub)
			})

			// HDFS
			r.Route("/hdfs/clusters", func(r chi.Router) {
				r.Get("/", bigDataHandler.ListHDFSClusters)
				r.With(middleware.RequireAdminOrOperator()).Post("/", bigDataHandler.CreateHDFSCluster)
				r.Get("/{id}", bigDataHandler.GetHDFSCluster)
				r.Get("/{id}/metrics", bigDataHandler.MetricsStub)
			})

			// NiFi
			r.Route("/nifi/instances", func(r chi.Router) {
				r.Get("/", bigDataHandler.ListNiFiInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", bigDataHandler.CreateNiFiInstance)
				r.Get("/{id}", bigDataHandler.GetNiFiInstance)
				r.Get("/{id}/process-groups", bigDataHandler.ListNiFiProcessGroups)
				r.Get("/{id}/metrics", bigDataHandler.MetricsStub)
			})

			// MLFlow
			r.Route("/mlflow/instances", func(r chi.Router) {
				r.Get("/", mlaiHandler.ListMLFlowInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", mlaiHandler.CreateMLFlowInstance)
				r.Get("/{id}", mlaiHandler.GetMLFlowInstance)
				r.Get("/{id}/experiments", mlaiHandler.ListExperiments)
				r.Get("/{id}/experiments/{expId}", mlaiHandler.GetExperiment)
				r.Get("/{id}/models", mlaiHandler.ListModels)
				r.Get("/{id}/models/{modelId}", mlaiHandler.GetModel)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}/models/{modelId}/stage", mlaiHandler.UpdateModelStage)
				r.Get("/{id}/metrics", mlaiHandler.MetricsStub)
			})

			// Feast
			r.Route("/feast/instances", func(r chi.Router) {
				r.Get("/", mlaiHandler.ListFeastInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", mlaiHandler.CreateFeastInstance)
				r.Get("/{id}", mlaiHandler.GetFeastInstance)
				r.Get("/{id}/entities", mlaiHandler.ListFeastEntities)
				r.Get("/{id}/feature-views", mlaiHandler.ListFeastFeatureViews)
				r.Get("/{id}/metrics", mlaiHandler.MetricsStub)
			})

			// JupyterHub
			r.Route("/jupyterhub/instances", func(r chi.Router) {
				r.Get("/", mlaiHandler.ListJupyterHubInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", mlaiHandler.CreateJupyterHubInstance)
				r.Get("/{id}", mlaiHandler.GetJupyterHubInstance)
				r.Get("/{id}/notebooks", mlaiHandler.ListNotebooks)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/notebooks", mlaiHandler.CreateNotebook)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}/notebooks/{nbId}", mlaiHandler.UpdateNotebook)
				r.With(middleware.RequireAdmin()).Delete("/{id}/notebooks/{nbId}", mlaiHandler.DeleteNotebook)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/notebooks/{nbId}/start", mlaiHandler.StartNotebook)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/notebooks/{nbId}/stop", mlaiHandler.StopNotebook)
				r.Get("/{id}/metrics", mlaiHandler.MetricsStub)
			})

			// LLM
			r.Route("/llm/deployments", func(r chi.Router) {
				r.Get("/", mlaiHandler.ListLLMDeployments)
				r.With(middleware.RequireAdminOrOperator()).Post("/", mlaiHandler.CreateLLMDeployment)
				r.Get("/{id}", mlaiHandler.GetLLMDeployment)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}", mlaiHandler.UpdateLLMDeployment)
				r.With(middleware.RequireAdmin()).Delete("/{id}", mlaiHandler.DeleteLLMDeployment)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/scale", mlaiHandler.ScaleLLMDeployment)
				r.Get("/{id}/metrics", mlaiHandler.MetricsStub)
			})

			// Superset
			r.Route("/superset/instances", func(r chi.Router) {
				r.Get("/", biHandler.ListSupersetInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", biHandler.CreateSupersetInstance)
				r.Get("/{id}", biHandler.GetSupersetInstance)
				r.Get("/{id}/dashboards", biHandler.ListSupersetDashboards)
				r.Get("/{id}/metrics", biHandler.MetricsStub)
			})

			// Metabase
			r.Route("/metabase/instances", func(r chi.Router) {
				r.Get("/", biHandler.ListMetabaseInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", biHandler.CreateMetabaseInstance)
				r.Get("/{id}", biHandler.GetMetabaseInstance)
				r.Get("/{id}/dashboards", biHandler.ListMetabaseDashboards)
				r.Get("/{id}/metrics", biHandler.MetricsStub)
			})

			// N8N
			r.Route("/n8n/instances", func(r chi.Router) {
				r.Get("/", biHandler.ListN8NInstances)
				r.With(middleware.RequireAdminOrOperator()).Post("/", biHandler.CreateN8NInstance)
				r.Get("/{id}", biHandler.GetN8NInstance)
				r.Get("/{id}/workflows", biHandler.ListN8NWorkflows)
				r.Get("/{id}/workflows/{wfId}", biHandler.GetN8NWorkflow)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/workflows/{wfId}/activate", biHandler.ActivateN8NWorkflow)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/workflows/{wfId}/deactivate", biHandler.DeactivateN8NWorkflow)
				r.Get("/{id}/metrics", biHandler.MetricsStub)
			})

			// Users
			r.Route("/users", func(r chi.Router) {
				r.With(middleware.RequireAdmin()).Get("/", userHandler.List)
				r.With(middleware.RequireAdmin()).Post("/", userHandler.Create)
				r.With(middleware.RequireAdmin()).Get("/{id}", userHandler.GetByID)
				r.With(middleware.RequireAdmin()).Put("/{id}", userHandler.Update)
				r.With(middleware.RequireAdmin()).Delete("/{id}", userHandler.Delete)
			})

			// DataFlows
			r.Route("/data-flows", func(r chi.Router) {
				r.Get("/templates", dataFlowHandler.ListTemplates)
				r.Get("/", dataFlowHandler.List)
				r.With(middleware.RequireAdminOrOperator()).Post("/", dataFlowHandler.Create)
				r.Get("/{id}", dataFlowHandler.GetByID)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}", dataFlowHandler.Update)
				r.With(middleware.RequireAdmin()).Delete("/{id}", dataFlowHandler.Delete)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}/nodes", dataFlowHandler.UpdateNodes)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}/edges", dataFlowHandler.UpdateEdges)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/deploy", dataFlowHandler.Deploy)
				r.With(middleware.RequireAdminOrOperator()).Post("/{id}/stop", dataFlowHandler.Stop)
			})

			// Monitoring
			r.Route("/monitoring/clusters/{id}", func(r chi.Router) {
				r.Get("/metrics", monitoringHandler.GetMetrics)
				r.Get("/nodes", monitoringHandler.ListNodes)
				r.Get("/nodes/{nodeId}", monitoringHandler.GetNode)
				r.Get("/resources", monitoringHandler.GetResources)
				r.Get("/pods", monitoringHandler.ListPods)
				r.Get("/events", monitoringHandler.ListEvents)
				r.Get("/alerts", monitoringHandler.ListAlerts)
			})

			// Dashboard
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/summary", dashboardHandler.Summary)
				r.Get("/services-status", dashboardHandler.ServicesStatus)
				r.Get("/recent-activity", dashboardHandler.RecentActivity)
				r.Get("/alerts", dashboardHandler.Alerts)
			})

			// Clusters
			r.Route("/clusters", func(r chi.Router) {
				r.Get("/", clusterHandler.List)
				r.With(middleware.RequireAdminOrOperator()).Post("/", clusterHandler.Create)
				r.Get("/{id}", clusterHandler.GetByID)
				r.With(middleware.RequireAdminOrOperator()).Put("/{id}", clusterHandler.Update)
				r.With(middleware.RequireAdmin()).Delete("/{id}", clusterHandler.Delete)
				r.Get("/{id}/namespaces", clusterHandler.ListNamespaces)
				r.Get("/{id}/namespaces/{nsId}/resources", clusterHandler.ListResources)
			})
		})

		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	return r
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}

	return srv.ListenAndServe()
}
