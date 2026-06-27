// Package provider detects which AI provider the user has configured
// based on the API base URL. Centralizes provider detection so that
// adding a new provider only requires changes here and in strategy tables.
package provider

import "strings"

// Type represents the identified API provider.
type Type int

const (
	Unknown Type = iota
	APIMart
	OpenAI
	OpenRouter
)

var names = map[Type]string{
	Unknown:    "unknown",
	APIMart:    "APIMart",
	OpenAI:     "OpenAI",
	OpenRouter: "OpenRouter",
}

func (t Type) String() string {
	if s, ok := names[t]; ok {
		return s
	}
	return "unknown"
}

// IsAsync returns true if this provider uses an async task-based model
// (submit → poll → download) for generation.
func (t Type) IsAsync() bool {
	return t == APIMart
}

// apimartDomains lists known APIMart-provided API domains.
var apimartDomains = []string{
	"apimart.ai",
	"apib.ai",
	"aiuxu.com",
	"aishuch.com",
}

// openrouterDomains lists domains where OpenRouter APIs are served.
var openrouterDomains = []string{
	"openrouter.ai",
}

// Detect identifies the provider from an API base URL.
// Returns Unknown if the URL doesn't match any known provider.
func Detect(baseURL string) Type {
	if baseURL == "" {
		return Unknown
	}
	for _, d := range apimartDomains {
		if strings.Contains(baseURL, d) {
			return APIMart
		}
	}
	for _, d := range openrouterDomains {
		if strings.Contains(baseURL, d) {
			return OpenRouter
		}
	}
	// Default to OpenAI-compatible for everything else
	return OpenAI
}

// IsAPIMart is a convenience wrapper around Detect.
func IsAPIMart(baseURL string) bool { return Detect(baseURL) == APIMart }

// IsOpenRouter is a convenience wrapper around Detect.
func IsOpenRouter(baseURL string) bool { return Detect(baseURL) == OpenRouter }
