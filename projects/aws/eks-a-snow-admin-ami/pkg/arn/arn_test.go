package arn_test

import (
	"testing"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/arn"
)

func TestForImageBuilderObject(t *testing.T) {
	tests := []struct {
		testName string
		account  string
		region   string
		kind     string
		name     string
		want     string
	}{
		{
			testName: "no spaces, no uppercase",
			name:     "componentname",
			region:   "us-west-2",
			account:  "111",
			kind:     "component",
			want:     "arn:aws:imagebuilder:us-west-2:111:component/componentname",
		},
		{
			testName: "no spaces, with uppercase",
			name:     "componentName",
			region:   "us-west-2",
			account:  "111",
			kind:     "component",
			want:     "arn:aws:imagebuilder:us-west-2:111:component/componentname",
		},
		{
			testName: "spaces and uppercase",
			name:     "Component Name",
			region:   "us-west-2",
			account:  "111",
			kind:     "component",
			want:     "arn:aws:imagebuilder:us-west-2:111:component/component-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arn.ForImageBuilderObject(tt.account, tt.region, tt.kind, tt.name); got != tt.want {
				t.Errorf("ForImageBuilderObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
