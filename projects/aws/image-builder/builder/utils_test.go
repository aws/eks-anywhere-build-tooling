package builder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSupportedReleaseBranchesSuccess(t *testing.T) {
	b := BuildOptions{
		ReleaseChannel: "1-24",
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if !SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		t.Fatalf("GetSupportedReleaseBranches error: supported branches does not contain the release channel"+
			": %s", b.ReleaseChannel)
	}
}

func TestGetSupportedReleaseBranchesFailure(t *testing.T) {
	b := BuildOptions{
		ReleaseChannel: "1-16",
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		t.Fatalf("GetSupportedReleaseBranches error: supported branches does not contain the release channel"+
			": %s", b.ReleaseChannel)
	}
}

func TestSetRhsmProxy(t *testing.T) {
	testcases := []struct {
		proxy            *ProxyConfig
		rhsm             *RhsmConfig
		outProxyHostname string
		outProxyPort     string
		name             string
		wantErr          bool
	}{
		{
			name: "Successful proxy set on rhsm",
			proxy: &ProxyConfig{
				HttpProxy:  "http://test.com:4578",
				HttpsProxy: "https://test.com:4578",
			},
			rhsm: &RhsmConfig{
				ServerHostname: "redhat-satellite.com",
			},
			outProxyHostname: "test.com",
			outProxyPort:     "4578",
			wantErr:          false,
		},
		{
			name: "Failed proxy set from parsing url",
			proxy: &ProxyConfig{
				HttpProxy:  "http:/test.com:4578",
				HttpsProxy: "http:/test.com:4578",
			},
			rhsm: &RhsmConfig{
				ServerHostname: "redhat-satellite.com",
			},
			wantErr: true,
		},
		{
			name: "Successful proxy set with no proxy",
			proxy: &ProxyConfig{
				HttpProxy:  "http://test.com:4578",
				HttpsProxy: "https://test.com:4578",
				NoProxy:    "no-proxy.com",
			},
			rhsm: &RhsmConfig{
				ServerHostname: "redhat-satellite.com",
			},
			outProxyHostname: "test.com",
			outProxyPort:     "4578",
			wantErr:          false,
		},
		{
			name: "Successful nil proxy set with satellite in no proxy",
			proxy: &ProxyConfig{
				HttpProxy:  "http://test.com:4578",
				HttpsProxy: "https://test.com:4578",
				NoProxy:    "redhat-satellite.com",
			},
			rhsm: &RhsmConfig{
				ServerHostname: "redhat-satellite.com",
			},
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := setRhsmProxy(tc.proxy, tc.rhsm)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Expected nil error. Got: %v", err)
			}
			if err == nil {
				if tc.rhsm.ProxyHostname != tc.outProxyHostname {
					t.Fatalf("Unexpected Proxy Hostname, Expected: %s, Got: %s", tc.outProxyHostname, tc.rhsm.ProxyHostname)
				}
				if tc.rhsm.ProxyPort != tc.outProxyPort {
					t.Fatalf("Unexpected Proxy Port, Expected: %s, Got: %s", tc.outProxyPort, tc.rhsm.ProxyPort)
				}
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working dir: %v", err)
	}
	testDir := filepath.Join(cwd, "testdata")
	os.MkdirAll(testDir, 0755)
	testcases := []struct {
		name    string
		path    string
		url     string
		wantErr bool
	}{
		{
			name:    "Download file to local dir",
			path:    filepath.Join(testDir, "test-A"),
			url:     "https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml",
			wantErr: false,
		},
		{
			name:    "Fail download with bad url",
			path:    filepath.Join(testDir, "test-B"),
			url:     "https:////anywhere-assets.eks.amazonaws.com/releases/eks-a/bad-manifest.yaml",
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := downloadFile(tc.path, tc.url)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Unexpected error status. Got error: %v, Expected error: %t", err, tc.wantErr)
			}
		})
	}
	os.RemoveAll(testDir)
}
