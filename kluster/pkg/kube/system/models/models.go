package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AppName string

const (
	AppCilium        = "cilium"
	AppOtelAgent     = "otel-agent"
	AppOtelCollector = "otel-collector"
	AppOtelGateway   = "otel-gateway"
)

type AppStatus string

const (
	AppCreate = "create"
	AppUpdate = "update"
	AppDelete = "delete"
	AppDone   = "done"
)

type Application struct {
	ID        uint64     `gorm:"column:id;primaryKey"`
	Name      AppName    `gorm:"column:name"`
	Status    AppStatus  `gorm:"column:status"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;default:null"`
}

func (Application) TableName() string {
	return "kube_system_applications"
}

type Resource struct {
	ID        uint64          `gorm:"column:id;primaryKey"`
	Name      string          `gorm:"column:name"`
	Namespace string          `gorm:"column:namespace"`
	Manifest  json.RawMessage `gorm:"column:manifest;type:jsonb"`
	UID       uuid.UUID       `gorm:"column:uid;type:uuid"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time      `gorm:"column:deleted_at;default:null"`
}

func (Resource) TableName() string {
	return "kluster_resources"
}

type LatestRsourceKindVersion struct {
	ResourceVersion uint64    `gorm:"column:resource_version;primaryKey"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (LatestRsourceKindVersion) TableName() string {
	return "kluster_latest_event_resource_version"
}
