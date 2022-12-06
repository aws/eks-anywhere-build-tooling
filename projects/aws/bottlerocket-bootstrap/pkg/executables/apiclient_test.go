package executables_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/executables"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/executables/mocks"
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
			args:          []string{"set", `kubernetes.cloud-providers=""`},
			returnError:   false,
			wantError:     false,
		},
		{
			testName:      "cloud provider aws",
			cloudProvider: "aws",
			args:          []string{"set", `kubernetes.cloud-providers="aws"`},
			returnError:   false,
			wantError:     false,
		},
		{
			testName:      "cloud provider empty, error",
			cloudProvider: "",
			args:          []string{"set", `kubernetes.cloud-providers=""`},
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

func TestSetKubernetesAllowUnsafeSysctls(t *testing.T) {
	tests := []struct {
		testName    string
		sysctls     []string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			sysctls:     []string{"sysctl-1", "sysctl-2", "sysctl-3"},
			args:        []string{"set", `allowed-unsafe-sysctls=["sysctl-1","sysctl-2","sysctl-3"]`},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			sysctls:     []string{"sysctl-1"},
			args:        []string{"set", `allowed-unsafe-sysctls=["sysctl-1"]`},
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

			err := tt.a.SetKubernetesAllowUnsafeSysctls(tc.sysctls)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestSetKernelRmemMax(t *testing.T) {
	tests := []struct {
		testName    string
		size        string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			size:        "0123",
			args:        []string{"set", "kernel.sysctl.net.core.rmem_max=0123"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			size:        "0123",
			args:        []string{"set", "kernel.sysctl.net.core.rmem_max=0123"},
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

			err := tt.a.SetKernelRmemMax(tc.size)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestSetKernelWmemMax(t *testing.T) {
	tests := []struct {
		testName    string
		size        string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			size:        "0123",
			args:        []string{"set", "kernel.sysctl.net.core.wmem_max=0123"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			size:        "0123",
			args:        []string{"set", "kernel.sysctl.net.core.wmem_max=0123"},
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

			err := tt.a.SetKernelWmemMax(tc.size)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestSetKernelCorePattern(t *testing.T) {
	tests := []struct {
		testName    string
		pattern     string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			pattern:     "pattern-1",
			args:        []string{"set", `kernel.core_pattern="pattern-1"`},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			pattern:     "pattern-1",
			args:        []string{"set", `kernel.core_pattern="pattern-1"`},
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

			err := tt.a.SetKernelCorePattern(tc.pattern)
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
