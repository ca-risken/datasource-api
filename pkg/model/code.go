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
	ErrorNotifiedAt     time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// CodeGitleaksRepository entity
type CodeGitleaksRepository struct {
	CodeGitleaksRepositoryID uint32    `gorm:"primary_key;column:code_gitleaks_repository_id"`
	CodeGitHubSettingID      uint32    `gorm:"column:code_github_setting_id"`
	RepositoryFullName       string    `gorm:"column:repository_full_name"`
	Status                   string    `gorm:"column:status"`
	StatusDetail             string    `gorm:"column:status_detail"`
	ScanAt                   time.Time `gorm:"column:scan_at"`
	CreatedAt                time.Time `gorm:"column:created_at"`
	UpdatedAt                time.Time `gorm:"column:updated_at"`
}

func (CodeGitleaksRepository) TableName() string {
	return "code_gitleaks_repository"
}

// CodeDependencyRepository entity
type CodeDependencyRepository struct {
	CodeDependencyRepositoryID uint32    `gorm:"primary_key;column:code_dependency_repository_id"`
	CodeGitHubSettingID        uint32    `gorm:"column:code_github_setting_id"`
	RepositoryFullName         string    `gorm:"column:repository_full_name"`
	Status                     string    `gorm:"column:status"`
	StatusDetail               string    `gorm:"column:status_detail"`
	ScanAt                     time.Time `gorm:"column:scan_at"`
	CreatedAt                  time.Time `gorm:"column:created_at"`
	UpdatedAt                  time.Time `gorm:"column:updated_at"`
}

func (CodeDependencyRepository) TableName() string {
	return "code_dependency_repository"
}

// CodeCodeScanRepository entity
type CodeCodeScanRepository struct {
	CodeCodeScanRepositoryID uint32    `gorm:"primary_key;column:code_codescan_repository_id"`
	CodeGitHubSettingID      uint32    `gorm:"column:code_github_setting_id"`
	RepositoryFullName       string    `gorm:"column:repository_full_name"`
	Status                   string    `gorm:"column:status"`
	StatusDetail             string    `gorm:"column:status_detail"`
	ScanAt                   time.Time `gorm:"column:scan_at"`
	CreatedAt                time.Time `gorm:"column:created_at"`
	UpdatedAt                time.Time `gorm:"column:updated_at"`
}

func (CodeCodeScanRepository) TableName() string {
	return "code_codescan_repository"
}

// CodeDependencySetting entity
type CodeDependencySetting struct {
	CodeGitHubSettingID uint32 `gorm:"primary_key;column:code_github_setting_id"`
	CodeDataSourceID    uint32
	ProjectID           uint32
	RepositoryPattern   string
	Status              string
	StatusDetail        string
	ScanAt              time.Time
	ErrorNotifiedAt     time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CodeCodeScanSetting struct {
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
	ErrorNotifiedAt     time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (c *CodeCodeScanSetting) TableName() string {
	return "code_codescan_setting"
}

// CodeGitleaksRepositoryStatus entity
type CodeGitleaksRepositoryStatus struct {
	CodeGitHubSettingID uint32    `gorm:"primary_key;column:code_github_setting_id"`
	RepositoryFullName  string    `gorm:"primary_key;column:repository_full_name"`
	Status              string    `gorm:"column:status"`
	StatusDetail        string    `gorm:"column:status_detail"`
	ScanAt              time.Time `gorm:"column:scan_at"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func (CodeGitleaksRepositoryStatus) TableName() string {
	return "code_gitleaks_repository"
}

// CodeDependencyRepositoryStatus entity
type CodeDependencyRepositoryStatus struct {
	CodeGitHubSettingID uint32    `gorm:"primary_key;column:code_github_setting_id"`
	RepositoryFullName  string    `gorm:"primary_key;column:repository_full_name"`
	Status              string    `gorm:"column:status"`
	StatusDetail        string    `gorm:"column:status_detail"`
	ScanAt              time.Time `gorm:"column:scan_at"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func (CodeDependencyRepositoryStatus) TableName() string {
	return "code_dependency_repository"
}

// CodeCodeScanRepositoryStatus entity
type CodeCodeScanRepositoryStatus struct {
	CodeGitHubSettingID uint32    `gorm:"primary_key;column:code_github_setting_id"`
	RepositoryFullName  string    `gorm:"primary_key;column:repository_full_name"`
	Status              string    `gorm:"column:status"`
	StatusDetail        string    `gorm:"column:status_detail"`
	ScanAt              time.Time `gorm:"column:scan_at"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func (CodeCodeScanRepositoryStatus) TableName() string {
	return "code_codescan_repository"
}
