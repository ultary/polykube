package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	APIGroup   string          `gorm:"column:api_group"`
	APIVersion string          `gorm:"column:api_version"`
	Kind       string          `gorm:"column:kind"`
	Name       string          `gorm:"column:name"`
	Namespace  string          `gorm:"column:namespace"`
	Manifest   json.RawMessage `gorm:"column:manifest;type:jsonb"`
	UID        uuid.UUID       `gorm:"column:uid;type:uuid"`
	CreatedAt  time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  *time.Time      `gorm:"column:deleted_at;default:null"`
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
