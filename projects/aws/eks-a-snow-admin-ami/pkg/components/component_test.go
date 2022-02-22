package components_test

import (
	"testing"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/components"
)

func TestComponentLastVersionARN(t *testing.T) {
	component := &components.Component{
		Name: "Component Name",
	}
	region := "us-west-2"
	account := "111"
	want := "arn:aws:imagebuilder:us-west-2:111:component/component-name/x.x.x"

	if got := component.LastVersionARN(account, region); got != want {
		t.Errorf("Component.LastVersionARN() = %v, want %v", got, want)
	}
}
