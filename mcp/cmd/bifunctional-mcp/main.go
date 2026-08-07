package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/config"
	orgserver "github.com/BifunctionalLabsCo/bifunctional-org-mcp/internal/server"
)

func main() {
	configPath := flag.String("config", envOrDefault("BIFUNCTIONAL_MCP_CONFIG", "config/repos.json"), "path to the repository registry")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(orgserver.Describe())
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("configuration failed", "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := orgserver.Build(cfg).Run(ctx, &mcp.StdioTransport{}); err != nil && ctx.Err() == nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
