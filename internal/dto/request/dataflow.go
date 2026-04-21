package request

type FlowNodeInput struct {
	NodeID       string   `json:"node_id" validate:"required"`
	Label        string   `json:"label" validate:"required"`
	ResourceType string   `json:"resource_type" validate:"required"`
	Cluster      *string  `json:"cluster"`
	Namespace    *string  `json:"namespace"`
	PositionX    float64  `json:"position_x"`
	PositionY    float64  `json:"position_y"`
	Config       []byte   `json:"config"`
}

type FlowEdgeInput struct {
	EdgeID       string  `json:"edge_id" validate:"required"`
	SourceNodeID string  `json:"source_node_id" validate:"required"`
	TargetNodeID string  `json:"target_node_id" validate:"required"`
	FlowType     string  `json:"flow_type" validate:"required"`
	Label        *string `json:"label"`
	Animated     bool    `json:"animated"`
	Config       []byte  `json:"config"`
}

type CreateDataFlowRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description *string         `json:"description"`
	Nodes       []FlowNodeInput `json:"nodes"`
	Edges       []FlowEdgeInput `json:"edges"`
}

type UpdateDataFlowRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status" validate:"omitempty,oneof=draft deployed running stopped failed"`
}

type UpdateFlowNodesRequest struct {
	Nodes []FlowNodeInput `json:"nodes" validate:"required"`
}

type UpdateFlowEdgesRequest struct {
	Edges []FlowEdgeInput `json:"edges" validate:"required"`
}

type CreateFromTemplateRequest struct {
	TemplateID  string  `json:"template_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}
