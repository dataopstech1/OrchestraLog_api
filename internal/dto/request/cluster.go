package request

type CreateClusterRequest struct {
	Name         string  `json:"name" validate:"required"`
	Region       string  `json:"region" validate:"required"`
	K8sVersion   string  `json:"k8s_version" validate:"required"`
	APIServerURL *string `json:"api_server_url"`
	Kubeconfig   *string `json:"kubeconfig"`
}

type UpdateClusterRequest struct {
	Name         *string `json:"name"`
	Status       *string `json:"status" validate:"omitempty,oneof=healthy warning critical offline"`
	K8sVersion   *string `json:"k8s_version"`
	APIServerURL *string `json:"api_server_url"`
	Kubeconfig   *string `json:"kubeconfig"`
	NodesTotal   *int    `json:"nodes_total"`
	NodesReady   *int    `json:"nodes_ready"`
}
