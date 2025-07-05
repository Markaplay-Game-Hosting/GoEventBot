package main

import (
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func (app *application) ServeFrontendProxy(w http.ResponseWriter, r *http.Request) {

	if !IsViteServerRunning() {
		app.logger.Error("Vite dev server is not running. Please start it with 'pnpm dev' from the ui/ directory")
		return
	}

	viteProxy := app.getViteProxy()

	viteProxy.ServeHTTP(w, r)
}

func (app *application) getViteProxy() *httputil.ReverseProxy {
	// Define your Vite dev server URL (typically runs on port 5173)
	viteURL, err := url.Parse("http://localhost:5173")
	if err != nil {
		app.logger.Error("Failed to parse Vite URL:", "error", err)
	}

	// Create reverse proxy for Vite dev server
	viteProxy := httputil.NewSingleHostReverseProxy(viteURL)

	// Enhanced error handling for Vite proxy
	viteProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		app.logger.Error("Vite proxy error", "error", err)
		http.Error(w, "Vite dev server unavailable", http.StatusBadGateway)
	}

	return viteProxy
}

// IsViteServerRunning
// Check if Vite dev server is running
func IsViteServerRunning() bool {
	conn, err := net.Dial("tcp", "localhost:5173")
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func IsProduction() bool {
	return strings.Contains(strings.ToLower(os.Getenv("GO_ENV")), "prod")
}

func doesFileExist(f fs.FS, path string) bool {
	_, err := fs.Stat(f, path)
	return err == nil
}
