package utils

import (
	"time"
)

type Configuration struct {
	BaseURL    string
	OrgID      string
	DataDockID string
	Token      string

	SkipTLSVerify  bool
	RequestTimeout time.Duration
	MaxRetries     int

	KeycloakBaseURL      string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
	KeycloakUsername     string
	KeycloakPassword     string

	MinIORegion    string
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL    string
	MinIOUseOIDC   string
}

type Response struct {
	Status   string
	Data     any
	Error    string
	HTTPCode int
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)
