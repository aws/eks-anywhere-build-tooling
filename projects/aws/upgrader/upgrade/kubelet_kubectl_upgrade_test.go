package upgrade_test

import (
	"context"
	"fmt"
	"os"
	"testing"
)

const (
	kubeletConf = "/etc/sysconfig/kubelet"
)

func TestKubeletKubectlUpgrade(t *testing.T) {
	ctx := context.TODO()
	tt := newUpgraderTest(t)
	tt.s.EXPECT().Executable().Return("/foo/eks-upgrades/tools", nil).Times(1)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubectl")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "stop", "kubelet").Return(nil, nil)
	tt.s.EXPECT().Stat(fmt.Sprintf("%s/%s.bk", upgCompBinDir, "kubelet")).Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "kubeadm", "version", "-oshort").Return([]byte("v1.25.x"), nil)
	tt.s.EXPECT().Stat(kubeletConf).Return(nil, os.ErrNotExist)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "daemon-reload").Return(nil, nil)
	tt.s.EXPECT().ExecCommand(ctx, "systemctl", "restart", "kubelet").Return(nil, nil)
	tt.u.KubeletKubectlUpgrade(ctx)
}
