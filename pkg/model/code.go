package model

import "time"

// CodeDataSource entity
type CodeDataSource struct {
	CodeDataSourceID uint32 `gorm:"primary_key"`
	Name             string
	Description      string
	MaxScore         float32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CodeGithubSetting entity
type CodeGithubSetting struct {
	CodeGithubSettingID uint32 `gorm:"primary_key"`
	Name                string
	ProjectID           uint32
	Type                string
	BaseURL             string
	TargetResource      string
	GithubUser          string
	PersonalAccessToken string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CodeGitleaksSetting entity
type CodeGitleaksSetting struct {
	CodeGithubSettingID uint32 `gorm:"primary_key"`
	CodeDataSourceID    uint32
	ProjectID           uint32
	RepositoryPattern   string
	ScanPublic          bool
	ScanInternal        bool
	ScanPrivate         bool
	Status              string
	StatusDetail        string
	ScanAt              time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CodeGithubEnterpriseOrg entity
type CodeGithubEnterpriseOrg struct {
	CodeGithubSettingID uint32 `gorm:"primary_key"`
	Organization        string `gorm:"primary_key"`
	ProjectID           uint32
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
