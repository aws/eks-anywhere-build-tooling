package upgrade_test

import (
	"context"
	"fmt"
	"testing"
)

const (
	cmpDir = "/foo/binaries"
)

func TestCniPluginsUpgrade(t *testing.T) {
	ctx := context.TODO()
	tt := newUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "cp", "-rf", fmt.Sprintf("%s/cni-plugins/.", cmpDir), "/").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "/opt/cni/bin/loopback", "--version").Return(nil, nil)
	tt.u.CniPluginsUpgrade(ctx)
}
