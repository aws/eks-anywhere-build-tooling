package upgrade_test

import (
	"context"
	"fmt"
	"testing"
)

func TestContainerdUpgrade(t *testing.T) {
	ctx := context.TODO()
	tt := newUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/containerd/.", cmpDir), "/").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "containerd", "--version").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "containerd").Return(nil, nil)
	tt.u.ContainerdUpgrade(ctx)
}
