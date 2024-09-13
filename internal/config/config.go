package config

type headerKey string

const (
	AUTHORIZATIONHEADER headerKey = "authorization"
	TOKENPREFIX         headerKey = "Bearer "
	USERIDCONTEXTKEY    headerKey = "userID"
	SESSIONIDCONTEXTKEY headerKey = "sessionID"
)

// ServerConfig - shortener server configurations
type ServerConfig struct {
	Secret        string `json:"-"`               // Secret for encryption algorithms
	LogsLevel     string `json:"-"`               // LogsLevel - level of logger
	GRPCAddress   string `json:"grpc_address"`    // GRPCAddress - address which system use to run gRPC server
	DatabaseDSN   string `json:"database_dsn"`    // DatabaseDSN - DSN path for DB connection
	MinIOAddress  string `json:"min_io_address"`  // MinIOAddress - address on which system use to connect to MinIO
	MinIOLogin    string `json:"min_io_login"`    // MinIOLogin - login which system use to connect to MinIO
	MinIOPassword string `json:"min_io_password"` // MinIOPassword - password which system use to connect to MinIO
}
