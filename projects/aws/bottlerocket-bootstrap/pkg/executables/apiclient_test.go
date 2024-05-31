package executables_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/executables"
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/executables/mocks"
)

type apiClientTest struct {
	*WithT
	a *executables.APIClient
	e *mocks.MockExecutable
}

func newAPIClientTest(t *testing.T) *apiClientTest {
	e := mocks.NewMockExecutable(gomock.NewController(t))
	return &apiClientTest{
		WithT: NewWithT(t),
		a:     &executables.APIClient{Executable: e},
		e:     e,
	}
}

func TestSetKubernetesCloudProvider(t *testing.T) {
	tests := []struct {
		testName      string
		cloudProvider string
		args          []string
		wantError     bool
		returnError   bool
	}{
		{
			testName:      "cloud provider empty string",
			cloudProvider: "",
			args:          []string{"set", `kubernetes.cloud-provider=""`},
			returnError:   false,
			wantError:     false,
		},
		{
			testName:      "cloud provider aws",
			cloudProvider: "aws",
			args:          []string{"set", `kubernetes.cloud-provider="aws"`},
			returnError:   false,
			wantError:     false,
		},
		{
			testName:      "cloud provider empty, error",
			cloudProvider: "",
			args:          []string{"set", `kubernetes.cloud-provider=""`},
			returnError:   true,
			wantError:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newAPIClientTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.a.SetKubernetesCloudProvider(tc.cloudProvider)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestSetKubernetesNodeIP(t *testing.T) {
	tests := []struct {
		testName    string
		ip          string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "ip empty string",
			ip:          "",
			args:        []string{"set", "kubernetes.node-ip="},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "valid",
			ip:          "1.2.3.4",
			args:        []string{"set", "kubernetes.node-ip=1.2.3.4"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			ip:          "1.2.3.4",
			args:        []string{"set", "kubernetes.node-ip=1.2.3.4"},
			returnError: true,
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newAPIClientTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.a.SetKubernetesNodeIP(tc.ip)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestSetKubernetesProviderID(t *testing.T) {
	tests := []struct {
		testName    string
		id          string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			id:          "id-0",
			args:        []string{"set", "kubernetes.provider-id=id-0"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			id:          "id-0",
			args:        []string{"set", "kubernetes.provider-id=id-0"},
			returnError: true,
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newAPIClientTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.a.SetKubernetesProviderID(tc.id)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestAPISet(t *testing.T) {
	tests := []struct {
		testName    string
		setting     *executables.APISetting
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName: "kernel config",
			setting: &executables.APISetting{
				Kernel: &executables.Kernel{
					Sysctl: map[string]string{
						"net.core.rmem_max":   "0123",
						"net.core.wmem_max":   "0123",
						"kernel.core_pattern": "/var/corefile/core.%e.%p.%h.%t",
					},
				},
				Kubernetes: &executables.Kubernetes{
					AllowedUnsafeSysctls: []string{
						"net.ipv4.tcp_mtu_probing",
					},
				},
			},
			args:        []string{"set", "--json", `{"kernel":{"sysctl":{"kernel.core_pattern":"/var/corefile/core.%e.%p.%h.%t","net.core.rmem_max":"0123","net.core.wmem_max":"0123"}},"kubernetes":{"allowed-unsafe-sysctls":["net.ipv4.tcp_mtu_probing"]}}`},
			returnError: false,
			wantError:   false,
		},
		{
			testName: "exec error",
			setting: &executables.APISetting{
				Kernel: &executables.Kernel{
					Sysctl: map[string]string{
						"kernel.core_pattern": "pattern-1",
					},
				},
			},
			args:        []string{"set", "--json", `{"kernel":{"sysctl":{"kernel.core_pattern":"pattern-1"}}}`},
			returnError: true,
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newAPIClientTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.a.Set(tc.setting)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestReboot(t *testing.T) {
	tests := []struct {
		testName    string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "success",
			args:        []string{"reboot"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			args:        []string{"reboot"},
			returnError: true,
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newAPIClientTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.a.Reboot()
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}
