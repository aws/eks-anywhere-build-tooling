package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	kubecmd "k8s.io/client-go/tools/clientcmd"
	kubecmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGetCurrentContextAPIServer(t *testing.T) {
	tests := []struct {
		name       string
		update     func(*kubecmdapi.Config)
		wantServer string
		wantError  string
	}{
		{
			name:       "local API server",
			wantServer: "https://192.168.1.10:6443",
		},
		{
			name: "discovery API server",
			update: func(config *kubecmdapi.Config) {
				config.Clusters["local"].Server = "https://192.168.1.100:6443"
			},
			wantServer: "https://192.168.1.100:6443",
		},
		{
			name: "current context selects cluster",
			update: func(config *kubecmdapi.Config) {
				config.Clusters["discovery"] = &kubecmdapi.Cluster{Server: "https://192.168.1.100:6443"}
				config.Contexts["discovery"] = &kubecmdapi.Context{Cluster: "discovery"}
			},
			wantServer: "https://192.168.1.10:6443",
		},
		{
			name: "IPv6 API server",
			update: func(config *kubecmdapi.Config) {
				config.Clusters["local"].Server = "https://[fd00::10]:6443"
			},
			wantServer: "https://[fd00::10]:6443",
		},
		{
			name: "current context is empty",
			update: func(config *kubecmdapi.Config) {
				config.CurrentContext = ""
			},
			wantError: "current context is empty",
		},
		{
			name: "current context is missing",
			update: func(config *kubecmdapi.Config) {
				config.CurrentContext = "missing"
			},
			wantError: `current context "missing" not found`,
		},
		{
			name: "current context is nil",
			update: func(config *kubecmdapi.Config) {
				config.Contexts["local"] = nil
			},
			wantError: `current context "local" not found`,
		},
		{
			name: "cluster name is empty",
			update: func(config *kubecmdapi.Config) {
				config.Contexts["local"].Cluster = ""
			},
			wantError: `cluster name is empty for current context "local"`,
		},
		{
			name: "cluster is missing",
			update: func(config *kubecmdapi.Config) {
				config.Contexts["local"].Cluster = "missing"
			},
			wantError: `cluster "missing" referenced by current context "local" not found`,
		},
		{
			name: "cluster is nil",
			update: func(config *kubecmdapi.Config) {
				config.Clusters["local"] = nil
			},
			wantError: `cluster "local" referenced by current context "local" not found`,
		},
		{
			name: "API server is empty",
			update: func(config *kubecmdapi.Config) {
				config.Clusters["local"].Server = ""
			},
			wantError: `API server is empty for cluster "local"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := validKubeConfig()
			if tt.update != nil {
				tt.update(&config)
			}

			server, err := getCurrentContextAPIServer(config)
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("getCurrentContextAPIServer() error = %v, want error containing %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("getCurrentContextAPIServer() error = %v, want nil", err)
			}
			if server != tt.wantServer {
				t.Fatalf("getCurrentContextAPIServer() = %q, want %q", server, tt.wantServer)
			}
		})
	}
}

func TestGetApiServerFromKubeConfig(t *testing.T) {
	path := writeKubeConfig(t, validKubeConfig())
	server, err := GetApiServerFromKubeConfig(path)
	if err != nil {
		t.Fatalf("GetApiServerFromKubeConfig() error = %v, want nil", err)
	}
	if server != "https://192.168.1.10:6443" {
		t.Fatalf("GetApiServerFromKubeConfig() = %q, want %q", server, "https://192.168.1.10:6443")
	}
}

func TestGetApiServerFromKubeConfigInvalidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(path, []byte("clusters: ["), 0o600); err != nil {
		t.Fatalf("writing kubeconfig: %v", err)
	}

	if _, err := GetApiServerFromKubeConfig(path); err == nil {
		t.Fatal("GetApiServerFromKubeConfig() error = nil, want an invalid kubeconfig error")
	}
}

func TestGetApiServerFromKubeConfigMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	if _, err := GetApiServerFromKubeConfig(path); err == nil {
		t.Fatal("GetApiServerFromKubeConfig() error = nil, want a file read error")
	}
}

func validKubeConfig() kubecmdapi.Config {
	return kubecmdapi.Config{
		CurrentContext: "local",
		Contexts: map[string]*kubecmdapi.Context{
			"local": {Cluster: "local"},
		},
		Clusters: map[string]*kubecmdapi.Cluster{
			"local": {Server: "https://192.168.1.10:6443"},
		},
	}
}

func writeKubeConfig(t *testing.T, config kubecmdapi.Config) string {
	t.Helper()

	data, err := kubecmd.Write(config)
	if err != nil {
		t.Fatalf("serializing kubeconfig: %v", err)
	}

	path := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("writing kubeconfig: %v", err)
	}
	return path
}
