package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/martianzhang/apimart-cli/internal/types"
)

const marketplaceURL = "https://api.apimart.ai/api/marketplace/models"

var modelType string

// modelsCmd represents the `apimart-cli models` command.
var modelsCmd = &cobra.Command{
	Use:   "models [image|video|chat]",
	Short: "List available AI models from APIMart marketplace",
	Long: `Query and display models from the APIMart marketplace.

Supports filtering by type: image, video, or chat. No API key required.

Examples:
  apimart-cli models
  apimart-cli models image
  apimart-cli models video
  apimart-cli models chat`,
	Args: cobra.MaximumNArgs(1),
	RunE: runModels,
}

func runModels(cmd *cobra.Command, args []string) error {
	// Determine type filter
	mediaType := ""
	if len(args) > 0 {
		mediaType = args[0]
	}
	if modelType != "" {
		mediaType = modelType
	}

	url := marketplaceURL + "?sort=newest&page=1&page_size=50"
	if mediaType != "" {
		url += "&type=" + mediaType
	}

	client := httpProxyClient()
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result types.MarketplaceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("API returned error")
	}

	data := result.Data
	if len(data.Models) == 0 {
		fmt.Println("No models found.")
		return nil
	}

	// Build vendor → models map
	type group struct {
		vendor string
		models []types.MarketplaceModel
	}
	var groups []group
	vendorMap := make(map[string][]types.MarketplaceModel)

	for _, m := range data.Models {
		vName := "Other"
		if m.Vendor != nil && m.Vendor.Name != "" {
			vName = m.Vendor.Name
		}
		vendorMap[vName] = append(vendorMap[vName], m)
	}

	for vName, mods := range vendorMap {
		groups = append(groups, group{vName, mods})
	}

	// Print header
	title := "All Models"
	if mediaType != "" {
		title = strings.ToUpper(mediaType[:1]) + mediaType[1:] + " Models"
	}
	fmt.Printf("\n%s (%d total)\n\n", title, data.Total)

	for _, g := range groups {
		fmt.Printf("  %s:\n", g.vendor)
		for _, m := range g.models {
			price := formatPrice(m)
			tags := strings.Join(m.Tags, ", ")
			line := fmt.Sprintf("    %-30s", m.ModelName)
			line += fmt.Sprintf("  %-12s", price)
			if tags != "" {
				line += "  " + tags
			}
			fmt.Println(line)
		}
		fmt.Println()
	}
	return nil
}

func formatPrice(m types.MarketplaceModel) string {
	if !m.Pricing.HasPrice {
		return "—"
	}
	unit := m.Pricing.PriceUnit
	if unit == "" {
		unit = "/次"
	}
	switch m.Pricing.BillingType {
	case "per_token":
		return fmt.Sprintf("$%.4f%s", m.Pricing.StartingPrice, unit)
	default:
		return fmt.Sprintf("$%.4f%s", m.Pricing.StartingPrice, unit)
	}
}

// httpProxyClient returns an HTTP client that respects the configured proxy.
func httpProxyClient() *http.Client {
	transport := &http.Transport{}
	if httpProxy != "" {
		if proxyURL, err := url.Parse(httpProxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}
	return &http.Client{Transport: transport}
}

func init() {
	modelsCmd.Flags().StringVarP(&modelType, "type", "t", "", "Filter by type: image, video, chat")
	rootCmd.AddCommand(modelsCmd)
}
