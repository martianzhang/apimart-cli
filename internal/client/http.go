package client

import (
	"net/http"
	"net/url"
)

// ConfigureDefaultClient sets the global http.DefaultClient's transport
// to use the given proxy URL. When proxyURL is empty, it falls back to
// HTTP_PROXY / HTTPS_PROXY / NO_PROXY environment variables.
//
// Call this once at startup so that ALL HTTP requests in the application
// (including those using http.Get(), http.DefaultClient, image downloads,
// ideas search, etc.) respect the proxy configuration.
func ConfigureDefaultClient(proxyURL string) {
	transport := &http.Transport{}
	if proxyURL != "" {
		if parsed, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(parsed)
		}
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}
	http.DefaultClient = &http.Client{Transport: transport}
}
