package model

import (
	"time"

	"github.com/google/uuid"
)

type DataFlow struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	Status      string     `db:"status" json:"status"`
	TemplateID  *string    `db:"template_id" json:"template_id"`
	Layout      string     `db:"layout" json:"layout"`
	CreatedBy   *uuid.UUID `db:"created_by" json:"created_by"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type FlowNode struct {
	ID           uuid.UUID `db:"id" json:"id"`
	DataFlowID   uuid.UUID `db:"data_flow_id" json:"data_flow_id"`
	NodeID       string    `db:"node_id" json:"node_id"`
	Label        string    `db:"label" json:"label"`
	ResourceType string    `db:"resource_type" json:"resource_type"`
	Cluster      *string   `db:"cluster" json:"cluster"`
	Namespace    *string   `db:"namespace" json:"namespace"`
	Status       string    `db:"status" json:"status"`
	PositionX    float64   `db:"position_x" json:"position_x"`
	PositionY    float64   `db:"position_y" json:"position_y"`
	Config       []byte    `db:"config" json:"config,omitempty"`
	Metrics      []byte    `db:"metrics" json:"metrics,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type FlowEdge struct {
	ID           uuid.UUID `db:"id" json:"id"`
	DataFlowID   uuid.UUID `db:"data_flow_id" json:"data_flow_id"`
	EdgeID       string    `db:"edge_id" json:"edge_id"`
	SourceNodeID string    `db:"source_node_id" json:"source_node_id"`
	TargetNodeID string    `db:"target_node_id" json:"target_node_id"`
	FlowType     string    `db:"flow_type" json:"flow_type"`
	Label        *string   `db:"label" json:"label"`
	Animated     bool      `db:"animated" json:"animated"`
	Config       []byte    `db:"config" json:"config,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type FlowTemplate struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  *string   `db:"description" json:"description"`
	Category     *string   `db:"category" json:"category"`
	Nodes        []byte    `db:"nodes" json:"nodes"`
	Edges        []byte    `db:"edges" json:"edges"`
	ThumbnailURL *string   `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type DataFlowWithGraph struct {
	DataFlow
	Nodes []FlowNode `json:"nodes"`
	Edges []FlowEdge `json:"edges"`
}
