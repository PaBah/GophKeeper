package main

import (
	"os"
	"testing"

	"github.com/PaBah/GophKeeper/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name          string
		expectedValue []string
		envValues     []string
	}{
		{
			name:          "got from ENV",
			expectedValue: []string{":8888", "test", "info", "minio:9000", "test", "test", "test"},
			envValues:     []string{":8888", "test", "info", "minio:9000", "test", "test", "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &config.ServerConfig{}
			if tt.envValues != nil {
				os.Setenv("GRPC_ADDRESS", tt.envValues[0])
				os.Setenv("DATABASE_DSN", tt.envValues[1])
				os.Setenv("LOG_LEVEL", tt.envValues[2])
				os.Setenv("MINIO_ADDRESS", tt.envValues[3])
				os.Setenv("MINIO_LOGIN", tt.envValues[4])
				os.Setenv("MINIO_PASSWORD", tt.envValues[5])
			}
			ParseFlags(options)
			assert.Equal(t, options.GRPCAddress, tt.expectedValue[0], "Правльно распаршеный GRPC_ADDRESS")
			assert.Equal(t, options.DatabaseDSN, tt.expectedValue[1], "Правльно распаршеный DATABASE_DSN")
			assert.Equal(t, options.LogsLevel, tt.expectedValue[2], "Правльно распаршеный LOG_LEVEL")
			assert.Equal(t, options.MinIOAddress, tt.expectedValue[3], "Правльно распаршеный MINIO_ADDRESS")
			assert.Equal(t, options.MinIOLogin, tt.expectedValue[4], "Правльно распаршеный MINIO_LOGIN")
			assert.Equal(t, options.MinIOPassword, tt.expectedValue[5], "Правльно распаршеный MINIO_PASSWORD")
		})
	}
}
