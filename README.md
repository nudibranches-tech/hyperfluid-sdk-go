# Bifrost SDK

Go SDK for Hyperfluid data access via GraphQL, REST, and PostgreSQL.

## Quick Start

```bash
# Install
go get bifrost-for-developers/sdk

# globalConfigurationure
cp .env.template .env
# Edit .env with your credentials

# Verify setup
./run_all_tests.sh setup
```

## Usage

```go
import bifrost "bifrost-for-developers/sdk"

func main() {
    bifrost.InitFromEnv()
    
    // GraphQL
    response, _ := bifrost.Get(bifrost.Request{
        Type: bifrost.RequestGraphQL,
        GraphQL: &bifrost.GraphQLRequest{
            Query: `{ catalog { schema { table { field } } } }`,
        },
    })
    
    // REST API
    response, _ := bifrost.Get(bifrost.Request{
        Type: bifrost.RequestOpenAPI,
        OpenAPI: &bifrost.OpenAPIRequest{
            Catalog: "my_catalog",
            Schema:  "my_schema",
            Table:   "my_table",
            Method:  "GET",
            Params:  map[string]string{"limit": "10"},
        },
    })
    
    // PostgreSQL
    response, _ := bifrost.Get(bifrost.Request{
        Type: bifrost.RequestPostgres,
        Postgres: &bifrost.PostgresRequest{
            SQL: "SELECT * FROM catalog.schema.table LIMIT 10",
        },
    })
}
```

## globalConfigurationuration

### Required
- `HYPERFLUID_ORG_ID` - Your organization ID
- `HYPERFLUID_TOKEN` - API token (or use Keycloak)

### Optional
- `HYPERFLUID_BASE_URL` - API endpoint (default: `https://bifrost.hyperfluid.cloud`)
- `HYPERFLUID_POSTGRES_HOST` - PostgreSQL host (default: `bifrost.hyperfluid.cloud`)
- `HYPERFLUID_POSTGRES_PORT` - PostgreSQL port (default: `5432`)

### Keycloak (alternative to token)
- `KEYCLOAK_USERNAME` - Your username
- `KEYCLOAK_PASSWORD` - Your password
- `KEYCLOAK_BASE_URL` - Keycloak server
- `KEYCLOAK_CLIENT_ID` - Client ID
- `KEYCLOAK_REALM` - Realm name

### Testing
- `BIFROST_TEST_CATALOG` - Test catalog name
- `BIFROST_TEST_SCHEMA` - Test schema name
- `BIFROST_TEST_TABLE` - Test table name
- `BIFROST_TEST_COLUMNS` - Columns to select (comma-separated, e.g., `sell,list,acres`)

## Testing

```bash
./run_tests.sh setup       # Verify globalConfigurationuration
./run_tests.sh unit        # Fast tests
./run_tests.sh integration # Full tests
./run_tests.sh all         # Everything
```

## Examples

```bash
cd run_examples && go run .
```

## Local Development Setup

This section describes how to run the complete Hyperfluid stack locally for development and testing.

### 1. Start Hyperfluid Services

```bash
# Start the backend server
just dance

# In another terminal, start the frontend
npm run dev

# Initialize Keycloak with default users
just init-keycloak
```

### 2. globalConfigurationure Hyperfluid Environment

Access the Hyperfluid UI (typically `http://localhost:3000`) and perform the following steps:

#### 2.1 Create an Environment
1. Navigate to **Environments**
2. Click **Create New Environment**
3. Give it a name (e.g., `local-dev`)

#### 2.2 Deploy Required Services
Deploy the following services in your environment:
- **Trino** - Query engine
- **Minio** - Object storage

Wait for both services to be in "Running" state.

#### 2.3 Create Storage Infrastructure
1. **Create a Bucket**:
   - Go to Minio storage
   - Create a new bucket (e.g., `test-data`)

2. **Create Iceberg Container**:
   - Navigate to Data Catalogs
   - Create new Iceberg container
   - Link it to your Minio bucket

3. **Create Schema**:
   - Open the Iceberg container
   - Create a new schema (e.g., `test_schema`)

#### 2.4 Import Sample Data
1. Navigate to your schema
2. Click **Import Data** or **Create Table**
3. Upload the sample data file (CSV/Parquet)
4. globalConfigurationure the table schema and name (e.g., `sample_table`)
5. Convert the file to Iceberg table format

### 3. globalConfigurationure Keycloak Authentication

#### 3.1 Access Keycloak Console
1. Open Keycloak admin console (typically `http://localhost:8443`)
2. Login with the super admin user:
   - Username: `admin`
   - Password: (created during init-keycloak)

#### 3.2 globalConfigurationure Client
1. Navigate to **Clients**
2. Find and edit **`authentication-client`**
3. Under **Settings**:
   - Enable **Direct Access Grants**
   - Set **Access Type** to `confidential` or `public` as needed
4. Click **Save**

### 4. globalConfigurationure Environment Variables

Update your `.env` file with the local globalConfigurationuration:

```bash
# Hyperfluid globalConfigurationuration
HYPERFLUID_BASE_URL=http://localhost:8080
HYPERFLUID_ORG_ID=your-org-id-from-ui
HYPERFLUID_TOKEN=  # Leave empty to use Keycloak

# PostgreSQL (Trino)
HYPERFLUID_POSTGRES_HOST=localhost
HYPERFLUID_POSTGRES_PORT=5432
HYPERFLUID_POSTGRES_USER=your-username
HYPERFLUID_POSTGRES_DB=your-database

# Keycloak Authentication
KEYCLOAK_USERNAME=your-username
KEYCLOAK_PASSWORD=your-password
KEYCLOAK_BASE_URL=http://localhost:8443
KEYCLOAK_CLIENT_ID=authentication-client
KEYCLOAK_REALM=nudibranches-tech

# Test Data globalConfigurationuration
BIFROST_TEST_CATALOG=your-catalog
BIFROST_TEST_SCHEMA=test_schema
BIFROST_TEST_TABLE=sample_table
BIFROST_TEST_COLUMNS=column1,column2,column3  # Columns from your imported data
```

**Note**: Get your `HYPERFLUID_ORG_ID` from the Hyperfluid UI under your organization settings.

### 5. Run Tests

```bash
# Verify setup
./run_tests.sh setup

# Run all tests against local stack
./run_tests.sh all
```

### Troubleshooting

**Keycloak Authentication Fails**:
- Verify `authentication-client` has Direct Access Grants enabled
- Check username/password are correct
- Ensure Keycloak is running on the specified port

**Cannot Connect to Services**:
- Verify all services are running: `docker ps` or check service status
- Check port conflicts
- Ensure firewall allows local connections

**Table Not Found**:
- Verify the table was successfully created in Hyperfluid UI
- Check catalog/schema/table names match your `.env` globalConfigurationuration
- Ensure Trino service is properly connected to Minio

**Sample Data Issues**:
- Verify the data file format is supported (CSV, Parquet, etc.)
- Check column names match the schema definition
- Ensure no special characters in column names

## Project Structure

```
sdk/                  # Core SDK
  bifrost.go         # Public API
  setup_check_test.go # globalConfiguration tests
  core/              # Internal logic
run_examples/        # Usage examples
developpement-tests/ # Test suite
```

## Error Handling

```go
response, error := bifrost.Get(request)
if error != nil {
    // Handle request error
    log.Printf("Request failed: %v", error)
}
if response.Status != bifrost.StatusOK {
    // Handle API error
    log.Printf("API error: %s", response.Error)
}

// Use helper methods
if !response.IsOK() {
    // Alternative way to check status
}
if response.HasError() {
    // Check if error message is present
}
```

## License

Private SDK for internal use.
