package model

import "time"

// AzureDataSource entity
type AzureDataSource struct {
	AzureDataSourceID uint32 `gorm:"primary_key"`
	Name              string
	Description       string
	MaxScore          float32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Azure entity
type Azure struct {
	AzureID          uint32 `gorm:"primary_key column:azure_id"`
	Name             string
	ProjectID        uint32
	SubscriptionID   string `gorm:"column:subscription_id"`
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// RelAzureDataSource entity
type RelAzureDataSource struct {
	AzureID           uint32 `gorm:"primary_key column:azure_id"`
	AzureDataSourceID uint32 `gorm:"primary_key"`
	ProjectID         uint32
	Status            string
	StatusDetail      string
	ScanAt            time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
