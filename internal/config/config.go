package config

const (
	AUTHORIZATIONHEADER = "authorization"
	TOKENPREFIX         = "Bearer "
	TOKENCONTEXTKEY     = "token"
	USEREMAILCONTEXTKEY = "email"
	USERIDCONTEXTKEY    = "userID"
)

// ServerConfig - shortener server configurations
type ServerConfig struct {
	Secret      string `json:"-"`            // Secret for encryption algorithms
	LogsLevel   string `json:"-"`            // LogsLevel - level of logger
	GRPCAddress string `json:"grpc_address"` // GRPCAddress - address which system use to run gRPC server
	EnableHTTPS bool   `json:"enable_https"` // EnableHTTPS - flag to enable HTTPS server mode
	DatabaseDSN string `json:"database_dsn"` // DatabaseDSN - DSN path for DB connection
}
