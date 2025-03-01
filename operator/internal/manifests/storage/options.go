package storage

import (
	lokiv1 "github.com/grafana/loki/operator/apis/loki/v1"
)

// Options is used to configure Loki to integrate with
// supported object storages.
type Options struct {
	Schemas     []lokiv1.ObjectStorageSchema
	SharedStore lokiv1.ObjectStorageSecretType

	Azure        *AzureStorageConfig
	GCS          *GCSStorageConfig
	S3           *S3StorageConfig
	Swift        *SwiftStorageConfig
	AlibabaCloud *AlibabaCloudStorageConfig

	SecretName string
	SecretSHA1 string
	TLS        *TLSConfig

	OpenShift OpenShiftOptions
}

// CredentialMode returns which mode is used by the current storage configuration.
// This defaults to CredentialModeStatic, but can be CredentialModeToken
// or CredentialModeManaged depending on the object storage provide, the provided
// secret and whether the operator is running in a managed-auth cluster.
func (o Options) CredentialMode() lokiv1.CredentialMode {
	if o.Azure != nil {
		if o.OpenShift.ManagedAuthEnabled() {
			return lokiv1.CredentialModeManaged
		}

		if o.Azure.WorkloadIdentity {
			return lokiv1.CredentialModeToken
		}
	}

	if o.GCS != nil {
		if o.GCS.WorkloadIdentity {
			return lokiv1.CredentialModeToken
		}
	}

	if o.S3 != nil {
		if o.OpenShift.ManagedAuthEnabled() {
			return lokiv1.CredentialModeManaged
		}

		if o.S3.STS {
			return lokiv1.CredentialModeToken
		}
	}

	return lokiv1.CredentialModeStatic
}

// AzureStorageConfig for Azure storage config
type AzureStorageConfig struct {
	Env              string
	Container        string
	EndpointSuffix   string
	Audience         string
	Region           string
	WorkloadIdentity bool
}

// GCSStorageConfig for GCS storage config
type GCSStorageConfig struct {
	Bucket           string
	Audience         string
	WorkloadIdentity bool
}

// S3StorageConfig for S3 storage config
type S3StorageConfig struct {
	Endpoint string
	Region   string
	Buckets  string
	Audience string
	STS      bool
	SSE      S3SSEConfig
}

type S3SSEType string

const (
	SSEKMSType S3SSEType = "SSE-KMS"
	SSES3Type  S3SSEType = "SSE-S3"
)

type S3SSEConfig struct {
	Type                 S3SSEType
	KMSKeyID             string
	KMSEncryptionContext string
}

// SwiftStorageConfig for Swift storage config
type SwiftStorageConfig struct {
	AuthURL           string
	UserDomainName    string
	UserDomainID      string
	UserID            string
	DomainID          string
	DomainName        string
	ProjectID         string
	ProjectName       string
	ProjectDomainID   string
	ProjectDomainName string
	Region            string
	Container         string
}

// AlibabaCloudStorageConfig for AlibabaCloud storage config
type AlibabaCloudStorageConfig struct {
	Endpoint string
	Bucket   string
}

// TLSConfig for object storage endpoints. Currently supported only by:
// - S3
type TLSConfig struct {
	CA  string
	Key string
}

type OpenShiftOptions struct {
	Enabled          bool
	CloudCredentials CloudCredentials
}

type CloudCredentials struct {
	SecretName string
	SHA1       string
}

func (o OpenShiftOptions) ManagedAuthEnabled() bool {
	return o.CloudCredentials.SecretName != "" && o.CloudCredentials.SHA1 != ""
}
