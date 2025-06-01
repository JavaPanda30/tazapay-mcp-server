package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/log"

	tools "github.com/tazapay/tazapay-mcp-server/tools/register"
)

func initConfig(logger *slog.Logger) error {
	viper.AutomaticEnv()

	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
		viper.SetConfigName(".tazapay-mcp-server")
		viper.SetConfigType("yaml")

		readErr := viper.ReadInConfig()
		if readErr != nil {
			var notFoundErr viper.ConfigFileNotFoundError
			if !errors.As(readErr, &notFoundErr) {
				logger.ErrorContext(context.Background(), "Config read error", "error", readErr)
				return readErr
			}
		}
	}

	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

	if accessKey == "" || secretKey == "" {
		logger.ErrorContext(context.Background(), "Missing API credentials")
		return constants.ErrMissingAuthKeys
	}

	authString := accessKey + ":" + secretKey
	authToken := base64.StdEncoding.EncodeToString([]byte(authString))
	viper.Set("TAZAPAY_AUTH_TOKEN", authToken)

	logger.InfoContext(context.Background(), "Configuration initialized")

	return nil
}

func main() {
	// Create a logger configuration
	logConfig := log.Config{
		Level:    "info",                           // Example log level
		Format:   "json",                           // Example log format
		FilePath: viper.GetString("LOG_FILE_PATH"), // Can be set via env LOG_FILE_PATH or config
	}

	// Create the logger
	logger, _, logErr := log.New(logConfig) // Empty path = default path near binary
	if logErr != nil {
		os.Exit(1)
	}

	// Initialize config to obtain the env variables
	if err := initConfig(logger); err != nil {
		logger.ErrorContext(context.Background(), "failed to initialize config", "error", err)
		os.Exit(1)
	}

	// Create a new MCP server
	s := server.NewMCPServer("tazapay", "0.0.1")

	// Register tools with the server
	tools.RegisterTools(s, logger)

	logger.InfoContext(context.Background(), "Started Tazapay MCP Server.")

	// server := server.New(s)

	// Gracefully shutdown at completion of execution
	// defer server.Shutdown(context.Background())

	// Run the server
	serveErr := server.ServeStdio(s)
	if serveErr != nil {
		logger.ErrorContext(context.Background(), "server exited with error", "error", serveErr)
		os.Exit(1)
	}
}
