package logengineevent

import "time"

// StorageConfigEvent represents a per-tenant storage backend configuration
type StorageConfigEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Provider  string    `json:"provider"`  // S3, GCS, AZURE_BLOB, SFTP, LOCAL, GOOGLE_DRIVE
	Enabled   bool      `json:"enabled"`
	IsDefault bool      `json:"isDefault"` // If true, this is the tenant's default storage

	// Common settings
	PathPrefix string `json:"pathPrefix"` // Base path/prefix for this tenant

	// S3/MinIO settings
	S3Config *S3StorageConfig `json:"s3Config,omitempty"`

	// Google Cloud Storage settings
	GCSConfig *GCSStorageConfig `json:"gcsConfig,omitempty"`

	// Azure Blob Storage settings
	AzureConfig *AzureBlobConfig `json:"azureConfig,omitempty"`

	// SFTP settings
	SFTPConfig *SFTPStorageConfig `json:"sftpConfig,omitempty"`

	// Local filesystem settings
	LocalConfig *LocalStorageConfig `json:"localConfig,omitempty"`

	// Google Drive settings
	GoogleDriveConfig *GoogleDriveConfig `json:"googleDriveConfig,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StorageProvider constants
const (
	StorageProviderS3          = "S3"
	StorageProviderMinIO       = "MINIO"
	StorageProviderGCS         = "GCS"
	StorageProviderAzureBlob   = "AZURE_BLOB"
	StorageProviderSFTP        = "SFTP"
	StorageProviderLocal       = "LOCAL"
	StorageProviderGoogleDrive = "GOOGLE_DRIVE"
)

// S3StorageConfig holds S3/MinIO specific settings
type S3StorageConfig struct {
	Bucket       string `json:"bucket"`
	Region       string `json:"region"`
	Endpoint     string `json:"endpoint,omitempty"`     // For MinIO or S3-compatible
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`              // Encrypted in transit
	UsePathStyle bool   `json:"usePathStyle,omitempty"` // For MinIO
	UseSSL       bool   `json:"useSsl"`
}

// GCSStorageConfig holds Google Cloud Storage specific settings
type GCSStorageConfig struct {
	Bucket          string `json:"bucket"`
	ProjectID       string `json:"projectId"`
	CredentialsJSON string `json:"credentialsJson"` // Service account JSON (encrypted)
}

// AzureBlobConfig holds Azure Blob Storage specific settings
type AzureBlobConfig struct {
	AccountName   string `json:"accountName"`
	AccountKey    string `json:"accountKey"` // Encrypted
	ContainerName string `json:"containerName"`
	Endpoint      string `json:"endpoint,omitempty"` // For Azure Stack or emulator
}

// SFTPStorageConfig holds SFTP specific settings
type SFTPStorageConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password,omitempty"`       // Encrypted
	PrivateKey     string `json:"privateKey,omitempty"`     // Encrypted
	PrivateKeyPass string `json:"privateKeyPass,omitempty"` // Encrypted
	BasePath       string `json:"basePath"`
	KnownHostsFile string `json:"knownHostsFile,omitempty"`
}

// LocalStorageConfig holds local filesystem specific settings
type LocalStorageConfig struct {
	BasePath    string `json:"basePath"`
	Permissions string `json:"permissions,omitempty"` // e.g., "0755"
}

// GoogleDriveConfig holds Google Drive specific settings
type GoogleDriveConfig struct {
	CredentialsJSON string `json:"credentialsJson"` // OAuth or service account JSON (encrypted)
	FolderID        string `json:"folderId"`        // Root folder ID
	SharedDrive     bool   `json:"sharedDrive"`     // If using shared drive
	SharedDriveID   string `json:"sharedDriveId,omitempty"`
}

