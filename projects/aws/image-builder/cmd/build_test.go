package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/eks-anywhere-build-tooling/image-builder/builder"
)

func TestValidateSupportedHypervisor(t *testing.T) {
	testCases := []struct {
		testName     string
		buildOptions builder.BuildOptions
		wantErr      string
	}{
		{
			testName: "vSphere hypervisor",
			buildOptions: builder.BuildOptions{
				Hypervisor: "vsphere",
			},
			wantErr: "",
		},
		{
			testName: "AMI hypervisor",
			buildOptions: builder.BuildOptions{
				Hypervisor: "ami",
			},
			wantErr: "",
		},
		{
			testName: "Unknown hypervisor",
			buildOptions: builder.BuildOptions{
				Hypervisor: "unknown-hypervisor",
			},
			wantErr: "unknown-hypervisor is not supported yet. Please select one of vsphere,baremetal,nutanix,cloudstack,ami",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := validateSupportedHypervisors(tt.buildOptions.Hypervisor)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidateOSHypervisorCombinations(t *testing.T) {
	testCases := []struct {
		testName     string
		buildOptions builder.BuildOptions
		wantErr      string
	}{
		{
			testName: "Cloudstack hypervisor with Redhat OS",
			buildOptions: builder.BuildOptions{
				Hypervisor: "cloudstack",
				Os:         "redhat",
			},
			wantErr: "",
		},
		{
			testName: "AMI hypervisor with Ubuntu OS",
			buildOptions: builder.BuildOptions{
				Hypervisor: "ami",
				Os:         "ubuntu",
			},
			wantErr: "",
		},
		{
			testName: "Nutanix hypervisor with Redhat OS",
			buildOptions: builder.BuildOptions{
				Hypervisor: "nutanix",
				Os:         "redhat",
			},
			wantErr: "Invalid OS type. Only ubuntu OS is supported for Nutanix",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := validateOSHypervisorCombinations(tt.buildOptions.Os, tt.buildOptions.Hypervisor)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}
