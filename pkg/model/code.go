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

// CodeGitHubSetting entity
type CodeGitHubSetting struct {
	CodeGitHubSettingID uint32 `gorm:"primary_key;column:code_github_setting_id"`
	Name                string
	ProjectID           uint32
	Type                string
	BaseURL             string
	TargetResource      string
	GitHubUser          string `gorm:"column:github_user"`
	PersonalAccessToken string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (CodeGitHubSetting) TableName() string {
	return "code_github_setting"
}

// CodeGitleaksSetting entity
type CodeGitleaksSetting struct {
	CodeGitHubSettingID uint32 `gorm:"primary_key;column:code_github_setting_id"`
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

// CodeGitleaksCashe entity
type CodeGitleaksCache struct {
	CodeGitHubSettingID uint32 `gorm:"column:code_github_setting_id"`
	RepositoryFullName  string
	ScanAt              time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CodeDependencySetting entity
type CodeDependencySetting struct {
	CodeGitHubSettingID uint32 `gorm:"primary_key;column:code_github_setting_id"`
	CodeDataSourceID    uint32
	ProjectID           uint32
	Status              string
	StatusDetail        string
	ScanAt              time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
