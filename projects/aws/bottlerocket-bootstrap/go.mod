module github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.14.0
	github.com/aws/aws-sdk-go-v2/config v1.14.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.14.0
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/godbus/dbus/v5 v5.0.4
	github.com/golang/mock v1.4.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/yaml v1.2.0
)
