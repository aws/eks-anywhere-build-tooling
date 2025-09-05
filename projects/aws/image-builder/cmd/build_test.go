package cmd

import (
	"io/ioutil"
	"log"
	"os"
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
			wantErr: "",
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

func TestValidateOSVersionCombinations(t *testing.T) {
	testCases := []struct {
		testName     string
		buildOptions builder.BuildOptions
		wantErr      string
	}{
		{
			testName: "Ubuntu 20.04",
			buildOptions: builder.BuildOptions{
				Os:        "ubuntu",
				OsVersion: "20.04",
			},
			wantErr: "",
		},
		{
			testName: "Ubuntu 22.04",
			buildOptions: builder.BuildOptions{
				Os:        "ubuntu",
				OsVersion: "22.04",
			},
			wantErr: "",
		},
		{
			testName: "Ubuntu 24.04",
			buildOptions: builder.BuildOptions{
				Os:        "ubuntu",
				OsVersion: "24.04",
			},
			wantErr: "",
		},
		{
			testName: "Ubuntu 26.04",
			buildOptions: builder.BuildOptions{
				Os:        "ubuntu",
				OsVersion: "26.04",
			},
			wantErr: "26.04 is not a supported version of Ubuntu. Please select one of 20.04,22.04,24.04",
		},
		{
			testName: "Redhat 8",
			buildOptions: builder.BuildOptions{
				Os:        "redhat",
				OsVersion: "8",
			},
			wantErr: "",
		},
		{
			testName: "Redhat 9",
			buildOptions: builder.BuildOptions{
				Os:        "redhat",
				OsVersion: "9",
			},
			wantErr: "",
		},
		{
			testName: "Rockylinux 1",
			buildOptions: builder.BuildOptions{
				Os:        "rocky",
				OsVersion: "1",
			},
			wantErr: "rocky is not a supported OS",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := validateOSVersion(tt.buildOptions.Os, tt.buildOptions.OsVersion)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidateFirmware(t *testing.T) {
	testCases := []struct {
		testName     string
		buildOptions builder.BuildOptions
		wantErr      string
	}{
		{
			testName: "Ubuntu ova with efi",
			buildOptions: builder.BuildOptions{
				Os:         "ubuntu",
				OsVersion:  "22",
				Hypervisor: "vsphere",
				Firmware:   "efi",
			},
			wantErr: "",
		},
		{
			testName: "Ubuntu raw with efi",
			buildOptions: builder.BuildOptions{
				Os:         "ubuntu",
				OsVersion:  "22",
				Hypervisor: "baremetal",
				Firmware:   "efi",
			},
			wantErr: "",
		},
		{
			testName: "Ubuntu raw with bios",
			buildOptions: builder.BuildOptions{
				Os:         "ubuntu",
				OsVersion:  "20",
				Hypervisor: "baremetal",
				Firmware:   "bios",
			},
			wantErr: "Ubuntu Raw builds only support EFI firmware",
		},
		{
			testName: "Redhat raw with efi",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "9",
				Hypervisor: "baremetal",
				Firmware:   "efi",
			},
			wantErr: "",
		},
		{
			testName: "Redhat raw with no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "8",
				Hypervisor: "baremetal",
				Firmware:   "",
			},
			wantErr: "",
		},
		{
			testName: "Redhat 9 ova with bios",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "9",
				Hypervisor: "vsphere",
				Firmware:   "bios",
			},
			wantErr: "",
		},
		{
			testName: "Redhat 8 ova with efi",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "8",
				Hypervisor: "vsphere",
				Firmware:   "efi",
			},
			wantErr: "Only RedHat version 9 supports EFI firmware",
		},
		{
			testName: "Redhat 9 ova with efi",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "9",
				Hypervisor: "vsphere",
				Firmware:   "efi",
			},
			wantErr: "",
		},
		{
			testName: "Redhat raw with bad firmware",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				Hypervisor: "baremetal",
				Firmware:   "bad",
			},
			wantErr: "bad is not a firmware. Please select one of bios,efi",
		},
		{
			testName: "Redhat raw with unsupported efi version",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "8",
				Hypervisor: "baremetal",
				Firmware:   "efi",
			},
			wantErr: "Only RedHat version 9 supports EFI firmware",
		},
		{
			testName: "Redhat raw with unsupported bios version",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "9",
				Hypervisor: "baremetal",
				Firmware:   "bios",
			},
			wantErr: "RedHat version 9 Raw builds only support EFI firmware",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := validateFirmware(tt.buildOptions.Firmware, tt.buildOptions.Os, tt.buildOptions.OsVersion, tt.buildOptions.Hypervisor)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidateInputsFirmwareDefaulting(t *testing.T) {
	baremetalConfigFile = "test-baremetal-config"
	cloudstackConfigFile = baremetalConfigFile
	vSphereConfigFile = baremetalConfigFile
	err = ioutil.WriteFile(baremetalConfigFile, []byte(`{"rhel_username": "un", "rhel_password": "pw", "iso_url": "fake-iso", "iso_checksum": "fake", "iso_checksum_type": "sha256"}`), 0o644)
	if err != nil {
		log.Fatal(err)
	}
	testCases := []struct {
		testName     string
		buildOptions builder.BuildOptions
		wantFirmware string
	}{
		{
			testName: "Ubuntu 20.04 Raw no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "ubuntu",
				OsVersion:  "20.04",
				Hypervisor: "baremetal",
				Firmware:   "",
			},
			wantFirmware: "efi",
		},
		{
			testName: "Ubuntu 22.04 Raw no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "ubuntu",
				OsVersion:  "22.04",
				Hypervisor: "baremetal",
				Firmware:   "",
			},
			wantFirmware: "efi",
		},
		{
			testName: "Redhat 9 Raw no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "9",
				Hypervisor: "baremetal",
				Firmware:   "",
			},
			wantFirmware: "efi",
		},
		{
			testName: "Redhat 8 Raw no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "8",
				Hypervisor: "baremetal",
				Firmware:   "",
			},
			wantFirmware: "bios",
		},
		{
			testName: "Redhat 9 cloudstack no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "redhat",
				OsVersion:  "9",
				Hypervisor: "cloudstack",
				Firmware:   "",
			},
			wantFirmware: "bios",
		},
		{
			testName: "Ubuntu 22.04 vsphere no firmware",
			buildOptions: builder.BuildOptions{
				Os:         "ubuntu",
				OsVersion:  "22.04",
				Hypervisor: "vsphere",
				Firmware:   "",
			},
			wantFirmware: "bios",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			err := ValidateInputs(&tt.buildOptions)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantFirmware, tt.buildOptions.Firmware)
		})
	}
	os.RemoveAll(baremetalConfigFile)
}
