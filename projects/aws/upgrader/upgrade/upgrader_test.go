package upgrade_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	upgrade "github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade"
	mock_upgrade "github.com/aws/eks-anywhere-build-tooling/projects/aws/upgrader/upgrade/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const (
	kubeConfigPath = "/etc/kubernetes/admin.conf"
	testDataDir    = "testdata"
)

type upgraderTest struct {
	*WithT
	u *upgrade.InPlaceUpgrader
	s *mock_upgrade.MockSysCalls
}

func newSysCallsMock(s *mock_upgrade.MockSysCalls) upgrade.SysCalls {
	return upgrade.SysCalls{
		WriteFile: func(name string, data []byte, perm fs.FileMode) error {
			return s.WriteFile(name, data, perm)
		},
		ReadFile: func(name string) ([]byte, error) {
			return s.ReadFile(name)
		},
		OpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
			return s.OpenFile(name, flag, perm)
		},
		Stat: func(name string) (fs.FileInfo, error) {
			return s.Stat(name)
		},
		Executable: func() (string, error) {
			return s.Executable()
		},
		ExecCommand: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			return s.ExecCommand(ctx, name, arg...)
		},
		MkdirAll: func(name string, perm fs.FileMode) error {
			return s.MkdirAll(name, perm)
		},
	}
}

func newInPlaceUpgraderTest(t *testing.T, options ...upgrade.Option) *upgraderTest {
	s := mock_upgrade.NewMockSysCalls(gomock.NewController(t))
	upg := upgrade.NewInPlaceUpgrader(options...)
	upg.SysCalls = newSysCallsMock(s)
	return &upgraderTest{
		WithT: NewWithT(t),
		u:     &upg,
		s:     s,
	}
}

func TestCurrDir(t *testing.T) {
	tests := []struct {
		testName  string
		ifErr     bool
		returnErr string
	}{
		{
			testName:  "exec path found success",
			ifErr:     false,
			returnErr: "",
		}, {
			testName:  "path error",
			ifErr:     true,
			returnErr: "path error",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			tt := newInPlaceUpgraderTest(t)
			if tc.ifErr {
				tt.s.EXPECT().Executable().Return("", errors.New(tc.returnErr))
			} else {
				tt.s.EXPECT().Executable().Return("", nil)
			}
			_, err := tt.u.CurrDir()
			if tc.ifErr {
				tt.Expect(err).NotTo(BeNil())
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}

func TestBackUpAndReplace(t *testing.T) {
	tests := []struct {
		testName  string
		ifErr     bool
		returnErr string
	}{
		{
			testName:  "backup success",
			ifErr:     false,
			returnErr: "",
		},
		{
			testName:  "backup error",
			ifErr:     true,
			returnErr: "backup file path error",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			oldFile := "/home/dummy/path"
			backUpFolder := "home/dummy/backup"
			newFile := "path.bk"
			tt := newInPlaceUpgraderTest(t)
			tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s", backUpFolder, newFile)).Return(nil, errors.New(""))
			if tc.ifErr {
				tt.s.EXPECT().ReadFile(oldFile).Return(nil, errors.New(tc.returnErr))
			} else {
				tt.s.EXPECT().ReadFile(oldFile).Return([]byte{}, nil)
				tt.s.EXPECT().WriteFile(fmt.Sprintf("%s/%s", backUpFolder, newFile), []byte{}, fileMode416).Return(nil)
				tt.s.EXPECT().ReadFile(newFile).Return([]byte{}, nil)
				tt.s.EXPECT().WriteFile(oldFile, []byte{}, fileMode416).Return(nil)
			}
			err := tt.u.BackUpAndReplace(oldFile, backUpFolder, newFile)
			if tc.ifErr {
				tt.Expect(err).To(Not(BeNil()))
			} else {
				tt.Expect(err).To(BeNil())
			}
		})
	}
}
