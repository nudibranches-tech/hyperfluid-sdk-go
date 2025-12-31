package fluent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// S3Builder provides a fluent interface for S3/MinIO operations using OIDC STS
type S3Builder struct {
	client interface {
		GetConfig() utils.Configuration
	}

	errors []error

	bucket    string
	key       string
	s3Client  *s3.Client
	stsClient *sts.Client

	idToken     string
	sessionName string
	roleArn     string

	stsMethod   string // "oidc" or ""
	oidcEnabled bool
}

// NewS3Builder creates a new S3Builder instance configured for MinIO
func NewS3Builder(client interface {
	GetConfig() utils.Configuration
}) (*S3Builder, error) {
	cfg := client.GetConfig()
	err := verifyBasicConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Check if we should use OIDC or static credentials
	useOIDC := getEnvOrConfig(cfg, "MINIO_USE_OIDC", "false") == "true"

	if useOIDC {
		return newS3BuilderWithOIDC(client)
	}

	return newS3BuilderWithStaticCreds(client)
}

func verifyBasicConfig(cfg utils.Configuration) error {
	if cfg.MinIOEndpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT is required")
	}
	if cfg.MinIORegion == "" {
		return fmt.Errorf("MINIO_REGION is required")
	}
	return nil
}

// newS3BuilderWithStaticCreds creates S3Builder with static MinIO credentials
func newS3BuilderWithStaticCreds(client interface {
	GetConfig() utils.Configuration
}) (*S3Builder, error) {
	cfg := client.GetConfig()

	if cfg.MinIOAccessKey == "" {
		return nil, fmt.Errorf("MINIO_ACCESS_KEY is required")
	}
	if cfg.MinIOSecretKey == "" {
		return nil, fmt.Errorf("MINIO_SECRET_KEY is required")
	}

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.MinIORegion),
		config.WithBaseEndpoint(cfg.MinIOEndpoint),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.MinIOAccessKey,
			cfg.MinIOSecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load MinIO config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Builder{
		client:      client,
		s3Client:    s3Client,
		errors:      []error{},
		oidcEnabled: false,
	}, nil
}

// newS3BuilderWithOIDC creates S3Builder configured for OIDC STS
func newS3BuilderWithOIDC(client interface {
	GetConfig() utils.Configuration
}) (*S3Builder, error) {
	cfg := client.GetConfig()
	ctx := context.Background()

	// Create base config with anonymous credentials for STS
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.MinIORegion),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	isHttps, err := isHTTPS(cfg.MinIOEndpoint)
	if err != nil {
		return nil, fmt.Errorf("MinIO endpoint incorecctly formatted")
	}

	// Create STS client pointing to MinIO's STS endpoint
	stsClient := sts.NewFromConfig(awsCfg, func(o *sts.Options) {
		// MinIO STS endpoint is typically at the base endpoint
		o.BaseEndpoint = aws.String(cfg.MinIOEndpoint)
		o.EndpointOptions.DisableHTTPS = !isHttps
	})

	// Create S3 client (will be updated after STS)
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.MinIOEndpoint)
		o.UsePathStyle = true
		o.EndpointOptions.DisableHTTPS = !isHttps
	})

	return &S3Builder{
		client:      client,
		s3Client:    s3Client,
		stsClient:   stsClient,
		errors:      []error{},
		oidcEnabled: true,
	}, nil
}

// isHTTPS checks if endpoint uses HTTPS
func isHTTPS(endpoint string) (bool, error) {
	URL, err := url.Parse(endpoint)
	if err != nil {
		return false, err
	}
	return URL.Scheme == "https", nil
}

// OIDC sets OIDC JWT token for AssumeRoleWithWebIdentity
func (s *S3Builder) OIDC(idToken string) *S3Builder {
	if !s.oidcEnabled {
		s.errors = append(
			s.errors,
			fmt.Errorf("OIDC cannot be used with static credentials; enable MINIO_USE_OIDC=true"),
		)
		return s
	}

	if idToken == "" {
		s.errors = append(s.errors, fmt.Errorf("OIDC token cannot be empty"))
	}

	s.idToken = idToken
	s.stsMethod = "oidc"
	return s
}

// RoleArn sets the role ARN for AssumeRoleWithWebIdentity
func (s *S3Builder) RoleArn(roleArn string) *S3Builder {
	if roleArn == "" {
		s.errors = append(s.errors, fmt.Errorf("role ARN cannot be empty"))
	}
	s.roleArn = roleArn
	return s
}

// SessionName sets the session name
func (s *S3Builder) SessionName(sessionName string) *S3Builder {
	s.sessionName = sessionName
	return s
}

// assumeRoleWithWebIdentity calls MinIO STS and updates the S3 client
func (s *S3Builder) assumeRoleWithWebIdentity(ctx context.Context) error {
	if s.idToken == "" {
		return fmt.Errorf("OIDC token is required for STS")
	}

	// Build session name if not provided
	sessionName := s.sessionName
	if sessionName == "" {
		sessionName = fmt.Sprintf("minio-session-%d", time.Now().Unix())
	}

	// Build input for AssumeRoleWithWebIdentity
	// Note: RoleArn is optional for MinIO. MinIO determines permissions from JWT claims
	// when RoleArn is not provided or uses RolePolicy when it is provided
	input := &sts.AssumeRoleWithWebIdentityInput{
		WebIdentityToken: aws.String(s.idToken),
		RoleSessionName:  aws.String(sessionName),
		DurationSeconds:  aws.Int32(3600), // 1 hour
	}

	// RoleArn is optional for MinIO but required by AWS SDK
	if s.roleArn != "" {
		// User explicitly provided a RoleArn
		input.RoleArn = aws.String(s.roleArn)
	} else {
		// Use placeholder - MinIO ignores this and uses JWT claims for authorization
		input.RoleArn = aws.String("arn:xxx:xxx:xxx:xxxx")
	}

	// Call STS
	output, err := s.stsClient.AssumeRoleWithWebIdentity(ctx, input)
	if err != nil {
		return fmt.Errorf("AssumeRoleWithWebIdentity failed: %w", err)
	}

	if output.Credentials == nil {
		return fmt.Errorf("STS returned no credentials")
	}

	// Extract temporary credentials
	creds := output.Credentials
	accessKey := aws.ToString(creds.AccessKeyId)
	secretKey := aws.ToString(creds.SecretAccessKey)
	sessionToken := aws.ToString(creds.SessionToken)

	// Create new credentials provider with STS credentials
	staticCreds := credentials.NewStaticCredentialsProvider(
		accessKey,
		secretKey,
		sessionToken,
	)

	// Get MinIO config
	cfg := s.client.GetConfig()
	// Recreate AWS config with new credentials
	ctx2 := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx2,
		config.WithRegion(cfg.MinIORegion),
		config.WithBaseEndpoint(cfg.MinIOEndpoint),
		config.WithCredentialsProvider(staticCreds),
	)
	if err != nil {
		return fmt.Errorf("failed to create config with STS credentials: %w", err)
	}

	isHttps, err := isHTTPS(cfg.MinIOEndpoint)
	if err != nil {
		return fmt.Errorf("MinIO endpoint incorecctly formatted")
	}

	// Recreate S3 client with STS credentials
	s.s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.EndpointOptions.DisableHTTPS = !isHttps
	})

	return nil
}

// Helper function to get config from environment or Configuration struct
func getEnvOrConfig(cfg utils.Configuration, key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	switch key {
	case "MINIO_ENDPOINT":
		if cfg.MinIOEndpoint != "" {
			return cfg.MinIOEndpoint
		}
	case "MINIO_ACCESS_KEY":
		if cfg.MinIOAccessKey != "" {
			return cfg.MinIOAccessKey
		}
	case "MINIO_SECRET_KEY":
		if cfg.MinIOSecretKey != "" {
			return cfg.MinIOSecretKey
		}
	case "MINIO_REGION":
		if cfg.MinIORegion != "" {
			return cfg.MinIORegion
		}
	case "MINIO_USE_OIDC":
		if cfg.MinIOUseOIDC != "" {
			return cfg.MinIOUseOIDC
		}
	}

	return fallback
}

// Bucket sets the S3 bucket name
func (s *S3Builder) Bucket(bucket string) *S3Builder {
	if bucket == "" {
		s.errors = append(s.errors, fmt.Errorf("bucket name cannot be empty"))
	}
	s.bucket = bucket
	return s
}

// Key sets the S3 object key (file path)
func (s *S3Builder) Key(key string) *S3Builder {
	if key == "" {
		s.errors = append(s.errors, fmt.Errorf("object key cannot be empty"))
	}
	s.key = key
	return s
}

// validate checks that all required fields are set and runs STS if needed
func (s *S3Builder) validate(ctx context.Context) error {
	if len(s.errors) > 0 {
		return fmt.Errorf("validation failed: %s", s.errors[0].Error())
	}
	if s.bucket == "" {
		return fmt.Errorf("%w: bucket required", utils.ErrInvalidRequest)
	}
	if s.key == "" {
		return fmt.Errorf("%w: key required", utils.ErrInvalidRequest)
	}

	// If OIDC method is set, assume role before proceeding
	if s.stsMethod == "oidc" {
		return s.assumeRoleWithWebIdentity(ctx)
	}

	return nil
}

// S3Object represents a downloaded object from MinIO/S3
type S3Object struct {
	Bucket       string
	Key          string
	Size         *int64
	ContentType  string
	LastModified *time.Time
	Metadata     map[string]string
	Body         io.ReadCloser // stream the content
}

// Get retrieves the object from MinIO and returns a stream
func (s *S3Builder) Get(ctx context.Context) (*S3Object, error) {
	if err := s.validate(ctx); err != nil {
		return nil, err
	}

	result, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}

	// Return a struct with Body as io.ReadCloser for streaming
	obj := &S3Object{
		Bucket:       s.bucket,
		Key:          s.key,
		Size:         result.ContentLength,
		ContentType:  aws.ToString(result.ContentType),
		LastModified: result.LastModified,
		Metadata:     result.Metadata,
		Body:         result.Body, // caller is responsible for closing
	}

	return obj, nil
}

// validateList checks validation errors and runs STS if needed (no key required)
func (s *S3Builder) validateList(ctx context.Context) error {
	if len(s.errors) > 0 {
		return fmt.Errorf("validation failed: %s", s.errors[0].Error())
	}
	if s.bucket == "" {
		return fmt.Errorf("%w: bucket required", utils.ErrInvalidRequest)
	}

	if s.stsMethod == "oidc" {
		return s.assumeRoleWithWebIdentity(ctx)
	}

	return nil
}

// List lists objects in the bucket with optional prefix
func (s *S3Builder) List(ctx context.Context, prefix string) (*utils.Response, error) {
	if err := s.validateList(ctx); err != nil {
		return nil, err
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	result, err := s.s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return &utils.Response{
			Status:   utils.StatusError,
			Error:    fmt.Sprintf("failed to list objects from MinIO: %v", err),
			HTTPCode: http.StatusInternalServerError,
		}, err
	}

	objects := make([]map[string]interface{}, 0, len(result.Contents))
	for _, obj := range result.Contents {
		var lastModified *string
		if obj.LastModified != nil {
			s := obj.LastModified.Format(time.RFC3339)
			lastModified = &s
		}

		objects = append(objects, map[string]interface{}{
			"key":           aws.ToString(obj.Key),
			"size":          obj.Size,
			"last_modified": lastModified, // nil-safe
		})
	}

	return &utils.Response{
		Status: utils.StatusOK,
		Data: map[string]interface{}{
			"bucket":  s.bucket,
			"objects": objects,
			"count":   len(objects),
		},
		HTTPCode: http.StatusOK,
	}, nil
}
