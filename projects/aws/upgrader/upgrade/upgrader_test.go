package upgrade_test

import (
	"errors"
	"fmt"
	"io/fs"
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
	u *upgrade.Upgrader
	s *mock_upgrade.MockSysCalls
}

func newUpgraderTest(t *testing.T, options ...upgrade.Option) *upgraderTest {
	s := mock_upgrade.NewMockSysCalls(gomock.NewController(t))
	upg := upgrade.NewUpgrader(options...)
	upg.SysCalls = s
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
			tt := newUpgraderTest(t)
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
			tt := newUpgraderTest(t)
			tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s", backUpFolder, newFile)).Return(nil, errors.New(""))
			if tc.ifErr {
				tt.s.EXPECT().ReadFile(oldFile).Return(nil, errors.New(tc.returnErr))
			} else {
				tt.s.EXPECT().ReadFile(oldFile).Return([]byte{}, nil)
				tt.s.EXPECT().WriteFile(fmt.Sprintf("%s/%s", backUpFolder, newFile), []byte{}, fs.FileMode(0o640)).Return(nil)
				tt.s.EXPECT().ReadFile(newFile).Return([]byte{}, nil)
				tt.s.EXPECT().WriteFile(oldFile, []byte{}, fs.FileMode(0o640)).Return(nil)
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
