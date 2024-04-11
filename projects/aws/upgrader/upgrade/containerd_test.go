package upgrade_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestContainerdUpgrade(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "containerd").Return(nil, nil).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).To(BeNil())
}

func TestContainerdUpgradeComponentsDirError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("", errors.New("")).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestContainerdUpgradeFirstVersionCmdError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, errors.New("")).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestContainerdUpgradeRecursiveCopyError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/").Return(nil, errors.New("")).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestContainerdUpgradeSecondVersionCmdError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, errors.New("")).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestContainerdUpgradeDaemonReloadError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, errors.New("")).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}

func TestContainerdUpgradeContainerdRestartError(t *testing.T) {
	ctx := context.TODO()
	tt := newInPlaceUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "containerd").Return(nil, errors.New("")).Times(1)

	err := tt.u.ContainerdUpgrade(ctx)
	tt.Expect(err).ToNot(BeNil())
}
