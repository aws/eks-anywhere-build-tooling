package executables_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/executables"
	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/executables/mocks"
)

type fileSystemTest struct {
	*WithT
	f *executables.FileSystem
	e *mocks.MockExecutable
}

func newFileSystemTest(t *testing.T) *fileSystemTest {
	e := mocks.NewMockExecutable(gomock.NewController(t))
	return &fileSystemTest{
		WithT: NewWithT(t),
		f:     &executables.FileSystem{Mount: e, Mkfs: e},
		e:     e,
	}
}

func TestMountVolume(t *testing.T) {
	tests := []struct {
		testName    string
		device      string
		dir         string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			device:      "device-1",
			dir:         "dir-1",
			args:        []string{"device-1", "dir-1"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			device:      "device-1",
			dir:         "dir-1",
			args:        []string{"device-1", "dir-1"},
			returnError: true,
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newFileSystemTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.f.MountVolume(tc.device, tc.dir)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		testName    string
		partition   string
		args        []string
		wantError   bool
		returnError bool
	}{
		{
			testName:    "valid",
			partition:   "device-1",
			args:        []string{"device-1"},
			returnError: false,
			wantError:   false,
		},
		{
			testName:    "exec error",
			partition:   "device-1",
			args:        []string{"device-1"},
			returnError: true,
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newFileSystemTest(t)
			if tc.returnError {
				tt.e.EXPECT().Execute(tc.args).Return(nil, errors.New("exec error"))
			} else {
				tt.e.EXPECT().Execute(tc.args).Return([]byte{}, nil)
			}

			err := tt.f.Partition(tc.partition)
			if tc.wantError {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}
