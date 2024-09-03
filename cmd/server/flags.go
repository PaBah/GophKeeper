package main

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

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
	var logsLevel, databaseDSN, enableHTTPS, gRPCAddress, configFilePath string

	flag.StringVar(&configFilePath, "c", "", "path to config file")
	flag.StringVar(&options.GRPCAddress, "g", ":3200", "host:port on which gRPC run")
	flag.StringVar(&options.DatabaseDSN, "d", "host=localhost user=paulbahush dbname=gophkeeper password=", "database DSN address")
	flag.StringVar(&options.LogsLevel, "l", "debug", "logs level")
	flag.BoolVar(&options.EnableHTTPS, "s", true, "enable-https")
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
				if !isFlagPassed("s") {
					options.EnableHTTPS = fileConfig.EnableHTTPS
				}
				if !isFlagPassed("g") {
					options.GRPCAddress = fileConfig.GRPCAddress
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

	enableHTTPS, specified = os.LookupEnv("ENABLE_HTTPS")
	if specified {
		options.EnableHTTPS, _ = strconv.ParseBool(enableHTTPS)
	}
}
