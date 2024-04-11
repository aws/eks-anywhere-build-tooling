package upgrade_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	cmpDir = "/foo/binaries"
)

func TestCNIPluginsUpgrade(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/cni-plugins/.", cmpDir), "/").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, nil).Times(1)

	err := tt.u.CNIPluginsUpgrade(ctx)
	tt.Expect(err).To(BeNil())
}

func TestCNIPluginsUpgradeComponentsDirError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("", errors.New("")).Times(1)

	err := tt.u.CNIPluginsUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestCNIPluginsUpgradeFirstVersionCmdError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, errors.New("")).Times(1)

	err := tt.u.CNIPluginsUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestCNIPluginsUpgradeRecursiveCopyError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/cni-plugins/.", cmpDir), "/").Return(nil, errors.New("")).Times(1)

	err := tt.u.CNIPluginsUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestCNIPluginsUpgradeSecondVersionCmdError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/cni-plugins/.", cmpDir), "/").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, errors.New("")).Times(1)

	err := tt.u.CNIPluginsUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}
