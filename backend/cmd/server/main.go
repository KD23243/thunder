/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package main is the entry point for starting the Thunder server.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/asgardeo/thunder/internal/flow"
	"github.com/asgardeo/thunder/internal/system/cert"
	"github.com/asgardeo/thunder/internal/system/managers"

	"github.com/asgardeo/thunder/internal/oauth/jwt"
	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/log"
)

func main() {
	logger := log.GetLogger()

	thunderHome := getThunderHome(logger)

	cfg := initThunderConfigurations(logger, thunderHome)
	if cfg == nil {
		logger.Fatal("Failed to initialize configurations")
	}

	mux := initMultiplexer(logger)
	if mux == nil {
		logger.Fatal("Failed to initialize multiplexer")
	}

	initFlowService(logger)

	startServer(logger, cfg, mux, thunderHome)
}

// getThunderHome retrieves and return the Thunder home directory.
func getThunderHome(logger *log.Logger) string {
	// Parse project directory from command line arguments.
	projectHome := ""
	projectHomeFlag := flag.String("thunderHome", "", "Path to Thunder home directory")
	flag.Parse()

	if *projectHomeFlag != "" {
		logger.Info("Using thunderHome from command line argument", log.String("thunderHome", *projectHomeFlag))
		projectHome = *projectHomeFlag
	} else {
		// If no command line argument is provided, use the current working directory.
		dir, dirErr := os.Getwd()
		if dirErr != nil {
			logger.Fatal("Failed to get current working directory", log.Error(dirErr))
		}
		projectHome = dir
	}

	return projectHome
}

// initThunderConfigurations initializes the Thunder configurations.
func initThunderConfigurations(logger *log.Logger, thunderHome string) *config.Config {
	// Load the configurations.
	configFilePath := path.Join(thunderHome, "repository/conf/deployment.yaml")
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.Fatal("Failed to load configurations", log.Error(err))
	}

	// Load the server's private key for signing JWTs.
	if err := jwt.LoadPrivateKey(cfg, thunderHome); err != nil {
		logger.Fatal("Failed to load private key", log.Error(err))
	}

	// Initialize runtime configurations.
	if err := config.InitializeThunderRuntime(thunderHome, cfg); err != nil {
		logger.Fatal("Failed to initialize thunder runtime", log.Error(err))
	}

	return cfg
}

// initMultiplexer initializes the HTTP multiplexer and registers the services.
func initMultiplexer(logger *log.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	serviceManager := managers.NewServiceManager(mux)

	// Register the services.
	err := serviceManager.RegisterServices()
	if err != nil {
		logger.Fatal("Failed to register the services", log.Error(err))
	}

	return mux
}

// initFlowService initializes the flow service.
func initFlowService(logger *log.Logger) {
	svc := flow.GetFlowService()
	if err := svc.Init(); err != nil {
		logger.Fatal("Failed to initialize flow service", log.Error(err))
	}
}

// startServer starts the HTTP server with the given configurations and multiplexer.
func startServer(logger *log.Logger, cfg *config.Config, mux *http.ServeMux, thunderHome string) {
	// Get TLS configuration from the certificate and key files.
	tlsConfig, err := cert.GetTLSConfig(cfg, thunderHome)
	if err != nil {
		logger.Fatal("Failed to load TLS configuration", log.Error(err))
	}

	// Build the server address using hostname and port from the configurations.
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)

	ln, err := tls.Listen("tcp", serverAddr, tlsConfig) // listener bound to port
	if err != nil {
		logger.Fatal("Failed to start TLS listener", log.Error(err))
	}

	logger.Info("WSO2 Thunder server started...", log.String("address", serverAddr))

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // Mitigate Slowloris attacks
	}

	if err := server.Serve(ln); err != nil {
		logger.Fatal("Failed to serve requests", log.Error(err))
	}
}
