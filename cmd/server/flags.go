package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/PaBah/GophKeeper/internal/config"
	"github.com/PaBah/GophKeeper/internal/logger"
	"go.uber.org/zap"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// ParseFlags - initializer system configuration
func ParseFlags(options *config.ServerConfig) {
	var specified bool
	var logsLevel, databaseDSN, gRPCAddress, configFilePath, minIOAdress, minIOLogin, minIOPassword string

	flag.StringVar(&configFilePath, "c", "", "path to config file")
	flag.StringVar(&options.GRPCAddress, "g", ":3200", "host:port on which gRPC run")
	flag.StringVar(&options.DatabaseDSN, "d", "host=localhost user=paulbahush dbname=gophkeeper password=", "database DSN address")
	flag.StringVar(&options.LogsLevel, "l", "debug", "logs level")
	flag.StringVar(&options.MinIOAddress, "m", "127.0.0.1:9000", "address of minio")
	flag.StringVar(&options.MinIOLogin, "k", "admin", "login for minio")
	flag.StringVar(&options.MinIOPassword, "p", "password123", "password for minio")
	flag.Parse()

	var fileConfig config.ServerConfig
	if configFilePath != "" {
		file, err := os.Open(configFilePath)
		if err == nil {
			err = json.NewDecoder(file).Decode(&fileConfig)
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					logger.Log().Error("can not close file", zap.Error(err))
				}
			}(file)
			if err == nil {
				if !isFlagPassed("d") {
					options.DatabaseDSN = fileConfig.DatabaseDSN
				}
				if !isFlagPassed("g") {
					options.GRPCAddress = fileConfig.GRPCAddress
				}
				if !isFlagPassed("m") {
					options.MinIOAddress = fileConfig.MinIOAddress
				}
				if !isFlagPassed("k") {
					options.MinIOLogin = fileConfig.MinIOLogin
				}
				if !isFlagPassed("p") {
					options.MinIOPassword = fileConfig.MinIOPassword
				}
			}
		}
	}

	logsLevel, specified = os.LookupEnv("LOG_LEVEL")
	if specified {
		options.LogsLevel = logsLevel
	}

	databaseDSN, specified = os.LookupEnv("DATABASE_DSN")
	if specified {
		options.DatabaseDSN = databaseDSN
	}

	gRPCAddress, specified = os.LookupEnv("GRPC_ADDRESS")
	if specified {
		options.GRPCAddress = gRPCAddress
	}

	minIOAdress, specified = os.LookupEnv("MINIO_ADDRESS")
	if specified {
		options.MinIOAddress = minIOAdress
	}

	minIOLogin, specified = os.LookupEnv("MINIO_LOGIN")
	if specified {
		options.MinIOLogin = minIOLogin
	}

	minIOPassword, specified = os.LookupEnv("MINIO_PASSWORD")
	if specified {
		options.MinIOPassword = minIOPassword
	}
}
