package provider

// StorageProviderType defines supported storage provider types.
type StorageProviderType string

const (
	StorageProviderS3        StorageProviderType = "S3"
	StorageProviderR2        StorageProviderType = "R2"
	StorageProviderGCS       StorageProviderType = "GCS"
	StorageProviderAzureBlob StorageProviderType = "AZURE_BLOB"
	StorageProviderMinIO     StorageProviderType = "MINIO"
	StorageProviderSFTP      StorageProviderType = "SFTP"
	StorageProviderFTP       StorageProviderType = "FTP"
	StorageProviderGenericS3 StorageProviderType = "GENERIC_S3"
)

// StoragePurpose defines the intended use for a storage binding.
type StoragePurpose string

const (
	StoragePurposeLogs      StoragePurpose = "LOGS"
	StoragePurposeDocuments StoragePurpose = "DOCUMENTS"
	StoragePurposeBackups   StoragePurpose = "BACKUPS"
	StoragePurposeMedia     StoragePurpose = "MEDIA"
	StoragePurposeReports   StoragePurpose = "REPORTS"
	StoragePurposeGeneral   StoragePurpose = "GENERAL"
)

// AllStoragePurposes returns every supported purpose.
func AllStoragePurposes() []StoragePurpose {
	return []StoragePurpose{
		StoragePurposeLogs, StoragePurposeDocuments, StoragePurposeBackups,
		StoragePurposeMedia, StoragePurposeReports, StoragePurposeGeneral,
	}
}

// StorageProviderCategory distinguishes cloud vs local/on-premise providers.
type StorageProviderCategory string

const (
	StorageCategoryCloud StorageProviderCategory = "cloud"
	StorageCategoryLocal StorageProviderCategory = "local"
)

// StorageProviderInfo holds display metadata and field definitions for one storage provider.
type StorageProviderInfo struct {
	Value       string                  `json:"value"`
	Label       string                  `json:"label"`
	Description string                  `json:"description"`
	Category    StorageProviderCategory `json:"category"`
	Fields      []FieldSchema           `json:"fields"`
}

// StoragePurposeInfo holds display metadata for a storage purpose.
type StoragePurposeInfo struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// StorageScopeInfo holds display metadata for a scope option.
type StorageScopeInfo struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// StorageFormConfig is the complete form metadata the frontend needs to render the storage binding form.
type StorageFormConfig struct {
	Providers []StorageProviderInfo `json:"providers"`
	Scopes    []StorageScopeInfo   `json:"scopes"`
	Purposes  []StoragePurposeInfo `json:"purposes"`
}

// StorageProviderRegistry is the single source of truth for all storage providers.
// Adding a new provider = add one entry here; frontend auto-renders fields.
var StorageProviderRegistry = map[string]StorageProviderInfo{
	"S3": {
		Value:       "S3",
		Label:       "Amazon S3",
		Description: "Amazon Simple Storage Service — scalable cloud object storage",
		Category:    StorageCategoryCloud,
		Fields: []FieldSchema{
			{Key: "endpoint", Label: "Endpoint", Placeholder: "https://s3.amazonaws.com", IsURL: true, MaxLength: 500},
			{Key: "region", Label: "Region", Placeholder: "us-east-1", Required: true, MinLength: 2, MaxLength: 64, Options: []FieldOption{
				{Value: "us-east-1", Label: "US East (N. Virginia)"},
				{Value: "us-west-2", Label: "US West (Oregon)"},
				{Value: "eu-west-1", Label: "EU (Ireland)"},
				{Value: "ap-south-1", Label: "Asia Pacific (Mumbai)"},
				{Value: "ap-southeast-1", Label: "Asia Pacific (Singapore)"},
			}},
			{Key: "bucketName", Label: "Bucket Name", Placeholder: "my-storage-bucket", Required: true, MinLength: 3, MaxLength: 63},
			{Key: "accessKeyId", Label: "Access Key ID", Placeholder: "AKIAIOSFODNN7EXAMPLE", Required: true, Sensitive: true, MinLength: 16, MaxLength: 128},
			{Key: "secretAccessKey", Label: "Secret Access Key", Placeholder: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", Required: true, Sensitive: true, MinLength: 16, MaxLength: 128},
		},
	},
	"R2": {
		Value:       "R2",
		Label:       "Cloudflare R2",
		Description: "Cloudflare R2 — S3-compatible object storage with zero egress fees",
		Category:    StorageCategoryCloud,
		Fields: []FieldSchema{
			{Key: "accountId", Label: "Account ID", Placeholder: "Cloudflare Account ID", Required: true, MinLength: 16, MaxLength: 64},
			{Key: "endpoint", Label: "Endpoint Override", Placeholder: "https://{accountId}.r2.cloudflarestorage.com", IsURL: true, MaxLength: 500},
			{Key: "bucketName", Label: "Bucket Name", Placeholder: "my-r2-bucket", Required: true, MinLength: 3, MaxLength: 63},
			{Key: "accessKeyId", Label: "Access Key ID", Placeholder: "R2 Access Key ID", Required: true, Sensitive: true, MinLength: 16, MaxLength: 128},
			{Key: "secretAccessKey", Label: "Secret Access Key", Placeholder: "R2 Secret Access Key", Required: true, Sensitive: true, MinLength: 16, MaxLength: 128},
		},
	},
	"GCS": {
		Value:       "GCS",
		Label:       "Google Cloud Storage",
		Description: "Google Cloud Storage — unified object storage with global edge caching",
		Category:    StorageCategoryCloud,
		Fields: []FieldSchema{
			{Key: "endpoint", Label: "API Endpoint", Placeholder: "https://storage.googleapis.com", IsURL: true, MaxLength: 500},
			{Key: "projectId", Label: "Project ID", Placeholder: "my-gcp-project", Required: true, MinLength: 4, MaxLength: 64},
			{Key: "bucketName", Label: "Bucket Name", Placeholder: "my-gcs-bucket", Required: true, MinLength: 3, MaxLength: 63},
			{Key: "serviceAccountKey", Label: "Service Account Key (JSON)", Placeholder: "Base64-encoded service account JSON key", Required: true, Sensitive: true, MinLength: 10, MaxLength: 10000},
		},
	},
	"AZURE_BLOB": {
		Value:       "AZURE_BLOB",
		Label:       "Azure Blob Storage",
		Description: "Microsoft Azure Blob Storage — massively scalable cloud object storage",
		Category:    StorageCategoryCloud,
		Fields: []FieldSchema{
			{Key: "endpoint", Label: "Endpoint", Placeholder: "https://blob.core.windows.net", IsURL: true, MaxLength: 500},
			{Key: "storageAccountName", Label: "Storage Account Name", Placeholder: "mystorageaccount", Required: true, MinLength: 3, MaxLength: 24},
			{Key: "containerName", Label: "Container Name", Placeholder: "my-container", Required: true, MinLength: 3, MaxLength: 63},
			{Key: "accessKey", Label: "Access Key", Placeholder: "Azure storage access key", Required: true, Sensitive: true, MinLength: 10, MaxLength: 200},
		},
	},
	"MINIO": {
		Value:       "MINIO",
		Label:       "MinIO",
		Description: "MinIO — high-performance S3-compatible object storage (self-hosted or cloud)",
		Category:    StorageCategoryCloud,
		Fields: []FieldSchema{
			{Key: "endpoint", Label: "Endpoint", Placeholder: "https://minio.example.com:9000", Required: true, IsURL: true, MaxLength: 500},
			{Key: "bucketName", Label: "Bucket Name", Placeholder: "my-minio-bucket", Required: true, MinLength: 3, MaxLength: 63},
			{Key: "accessKeyId", Label: "Access Key ID", Placeholder: "MinIO access key", Required: true, Sensitive: true, MinLength: 4, MaxLength: 128},
			{Key: "secretAccessKey", Label: "Secret Access Key", Placeholder: "MinIO secret key", Required: true, Sensitive: true, MinLength: 8, MaxLength: 128},
			{Key: "useSSL", Label: "Use SSL", Placeholder: "true", MaxLength: 5},
		},
	},
	"SFTP": {
		Value:       "SFTP",
		Label:       "SFTP",
		Description: "SSH File Transfer Protocol — secure file transfer over SSH",
		Category:    StorageCategoryLocal,
		Fields: []FieldSchema{
			{Key: "host", Label: "Host", Placeholder: "sftp.example.com", Required: true, MinLength: 1, MaxLength: 255},
			{Key: "port", Label: "Port", Placeholder: "22", Required: true, MinLength: 1, MaxLength: 5},
			{Key: "username", Label: "Username", Placeholder: "sftp-user", Required: true, MinLength: 1, MaxLength: 128},
			{Key: "password", Label: "Password", Placeholder: "Password (optional if using private key)", Sensitive: true, MaxLength: 256},
			{Key: "privateKey", Label: "Private Key (PEM)", Placeholder: "-----BEGIN RSA PRIVATE KEY-----", Sensitive: true, MaxLength: 10000},
			{Key: "basePath", Label: "Base Path", Placeholder: "/data/storage", MaxLength: 500},
		},
	},
	"FTP": {
		Value:       "FTP",
		Label:       "FTP",
		Description: "File Transfer Protocol — traditional file transfer (TLS supported)",
		Category:    StorageCategoryLocal,
		Fields: []FieldSchema{
			{Key: "host", Label: "Host", Placeholder: "ftp.example.com", Required: true, MinLength: 1, MaxLength: 255},
			{Key: "port", Label: "Port", Placeholder: "21", Required: true, MinLength: 1, MaxLength: 5},
			{Key: "username", Label: "Username", Placeholder: "ftp-user", Required: true, MinLength: 1, MaxLength: 128},
			{Key: "password", Label: "Password", Placeholder: "FTP password", Sensitive: true, MaxLength: 256},
			{Key: "basePath", Label: "Base Path", Placeholder: "/data/storage", MaxLength: 500},
			{Key: "passive", Label: "Passive Mode", Placeholder: "true", MaxLength: 5},
			{Key: "useTLS", Label: "Use TLS", Placeholder: "true", MaxLength: 5},
		},
	},
	"GENERIC_S3": {
		Value:       "GENERIC_S3",
		Label:       "Generic S3-Compatible",
		Description: "Any S3-compatible storage provider (Ceph, Wasabi, Linode, etc.)",
		Category:    StorageCategoryCloud,
		Fields: []FieldSchema{
			{Key: "endpoint", Label: "Endpoint", Placeholder: "https://s3.us-east-1.wasabisys.com", Required: true, IsURL: true, MaxLength: 500},
			{Key: "region", Label: "Region", Placeholder: "us-east-1", MaxLength: 64},
			{Key: "bucketName", Label: "Bucket Name", Placeholder: "my-bucket", Required: true, MinLength: 3, MaxLength: 63},
			{Key: "accessKeyId", Label: "Access Key ID", Placeholder: "Access key", Required: true, Sensitive: true, MinLength: 4, MaxLength: 128},
			{Key: "secretAccessKey", Label: "Secret Access Key", Placeholder: "Secret key", Required: true, Sensitive: true, MinLength: 8, MaxLength: 128},
			{Key: "forcePathStyle", Label: "Force Path Style", Placeholder: "true", MaxLength: 5},
		},
	},
}

// StoragePurposeRegistry provides display metadata for each storage purpose.
var StoragePurposeRegistry = map[StoragePurpose]StoragePurposeInfo{
	StoragePurposeLogs:      {Value: "LOGS", Label: "Logs", Description: "Application and system logs"},
	StoragePurposeDocuments: {Value: "DOCUMENTS", Label: "Documents", Description: "User documents and files"},
	StoragePurposeBackups:   {Value: "BACKUPS", Label: "Backups", Description: "Database and system backups"},
	StoragePurposeMedia:     {Value: "MEDIA", Label: "Media", Description: "Images, videos, and media assets"},
	StoragePurposeReports:   {Value: "REPORTS", Label: "Reports", Description: "Generated reports and exports"},
	StoragePurposeGeneral:   {Value: "GENERAL", Label: "General", Description: "General purpose storage"},
}

// GetStorageFormConfig returns the complete form configuration for the frontend.
func GetStorageFormConfig() StorageFormConfig {
	providers := make([]StorageProviderInfo, 0, len(StorageProviderRegistry))
	for _, info := range StorageProviderRegistry {
		providers = append(providers, info)
	}

	purposes := make([]StoragePurposeInfo, 0, len(AllStoragePurposes()))
	for _, p := range AllStoragePurposes() {
		if info, ok := StoragePurposeRegistry[p]; ok {
			purposes = append(purposes, info)
		}
	}

	scopes := []StorageScopeInfo{
		{Value: "OwnerOnly", Label: "Owner Only", Description: "Only this tenant can use this binding"},
		{Value: "OwnerAndDescendants", Label: "Owner & Descendants", Description: "This tenant and all child tenants inherit this binding"},
		{Value: "ExplicitTenants", Label: "Explicit Tenants", Description: "Only explicitly specified tenants can use this binding"},
	}

	return StorageFormConfig{
		Providers: providers,
		Scopes:    scopes,
		Purposes:  purposes,
	}
}
