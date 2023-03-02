module github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.14.0
	github.com/aws/aws-sdk-go-v2/config v1.14.0
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.11.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.14.0
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/godbus/dbus/v5 v5.0.4
	github.com/golang/mock v1.4.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/vishvananda/netlink v1.1.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/yaml v1.2.0
)

replace (
	golang.org/x/net/http => golang.org/x/net/http v0.0.0-20220906165146-f3363e06e74c
	golang.org/x/text => golang.org/x/text v0.3.8
)
