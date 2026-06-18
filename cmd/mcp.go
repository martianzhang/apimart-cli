package cmd

import (
	"github.com/martianzhang/apimart-cli/internal/config"
	"github.com/martianzhang/apimart-cli/internal/mcp"
	"github.com/spf13/cobra"
)

// mcpCmd represents the `apimart-cli mcp` command.
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server for AI agent integration",
	Long: `Start an MCP (Model Context Protocol) server over stdio.

This allows AI agents (Claude Desktop, Cursor, etc.) to call APIMart
tools directly: generate images, generate videos, query models, etc.

Configuration is read from config.yaml, environment variables, and --config flag.

Example MCP host config:
{
  "mcpServers": {
    "apimart": {
      "command": "apimart-cli",
      "args": ["mcp"]
    }
  }
}
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config (optional) to get defaults
		cfg, _ := config.Load(cfgFile)

		mcpCfg := &mcp.Config{
			APIKey:  apiKey,
			BaseURL: apiBase,
			Proxy:   httpProxy,
			Output:  outputDir,
		}

		// Copy defaults from config if present
		if cfg != nil {
			if mcpCfg.APIKey == "" {
				mcpCfg.APIKey = cfg.APIKey
			}
			if mcpCfg.BaseURL == "" {
				mcpCfg.BaseURL = cfg.BaseURL
			}
			if mcpCfg.Proxy == "" {
				mcpCfg.Proxy = cfg.HTTPProxy
			}
			mcpCfg.Defaults = cfg.Defaults
		}

		return mcp.Run(mcpCfg)
	},
}

func init() {
	// Override PersistentPreRunE to skip the api key check for mcp command.
	// Some MCP tools (list_models, get_model_pricing) work without API key.
	mcpCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Load config to populate apiKey, apiBase, httpProxy, outputDir if not set via flags
		if c, err := config.Load(cfgFile); err == nil {
			if apiKey == "" {
				apiKey = c.APIKey
			}
			if apiBase == "" {
				apiBase = c.BaseURL
			}
			if httpProxy == "" {
				httpProxy = c.HTTPProxy
			}
			if !cmd.Flags().Changed("output") {
				outputDir = c.BaseURL // fallback stays as "."
			}
		}
		// Don't error on missing API key - tools will handle it gracefully
		return nil
	}

	rootCmd.AddCommand(mcpCmd)
}
