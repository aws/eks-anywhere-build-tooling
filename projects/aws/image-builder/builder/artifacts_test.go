package builder

import (
	"testing"
)

func TestGetArtifactsFilePathFromUrl(t *testing.T) {
	testcases := []struct {
		name    string
		url     string
		path    string
		wantErr bool
	}{
		{
			name:    "Successful path parsing",
			url:     "https://distro.eks.amazonaws.com/kubernetes-1-24/releases/24/artifacts/kubernetes/v1.24.16/bin/linux/amd64/kubeadm",
			path:    "kubernetes-1-24/releases/24/artifacts/kubernetes/v1.24.16/bin/linux/amd64/kubeadm",
			wantErr: false,
		},
		{
			name:    "Unsuccessful path parsing",
			url:     "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getArtifactFilePathFromUrl(tc.url)
			if got != tc.path {
				t.Fatalf("Unexpected path. Expected: %s, Got: %s", tc.path, got)
			}
			if (err != nil) != tc.wantErr {
				t.Fatalf("Unexpected error. Got: %v", err)
			}
		})
	}
}
