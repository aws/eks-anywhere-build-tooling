package projectpatchfixer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransformAutoscalerRouterKeepsOnlyClusterAPI(t *testing.T) {
	source := t.TempDir()
	router := filepath.Join(source, "cluster-autoscaler", "cloudprovider", "router")
	if err := os.MkdirAll(router, 0o755); err != nil {
		t.Fatal(err)
	}
	writeAutoscalerTestFile(t, router, "router_aws.go", "package router\n")
	writeAutoscalerTestFile(t, router, "router_clusterapi.go", "package router\n")
	writeAutoscalerTestFile(t, router, "router_all.go", `package router
import (
	_ "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	_ "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/clusterapi"
)
func init() {
	builder.SetDefaultCloudProvider(cloudprovider.GceProviderName)
}
`)

	changed, err := transformAutoscalerRouter(source)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("transformAutoscalerRouter() changed = false")
	}
	content := readAutoscalerTestFile(t, router, "router_all.go")
	if strings.Contains(content, "cloudprovider/aws") {
		t.Fatal("router still imports AWS provider")
	}
	if !strings.Contains(content, "cloudprovider/clusterapi") ||
		!strings.Contains(content, "cloudprovider.ClusterAPIProviderName") {
		t.Fatal("router does not retain Cluster API provider")
	}
	if _, err := os.Stat(filepath.Join(router, "router_aws.go")); !os.IsNotExist(err) {
		t.Fatalf("router_aws.go still exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(router, "router_clusterapi.go")); err != nil {
		t.Fatalf("router_clusterapi.go missing: %v", err)
	}
}

func TestTransformAutoscalerBuilderKeepsOnlyClusterAPI(t *testing.T) {
	source := t.TempDir()
	builder := filepath.Join(source, "cluster-autoscaler", "cloudprovider", "builder")
	if err := os.MkdirAll(builder, 0o755); err != nil {
		t.Fatal(err)
	}
	writeAutoscalerTestFile(t, builder, "builder_aws.go", "package builder\n")
	writeAutoscalerTestFile(t, builder, "builder_clusterapi.go", "package builder\n")
	writeAutoscalerTestFile(t, builder, "builder_all.go", `package builder
import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/clusterapi"
)
var AvailableCloudProviders = []string{
	cloudprovider.AwsProviderName,
	cloudprovider.ClusterAPIProviderName,
}
const DefaultCloudProvider = cloudprovider.GceProviderName
func buildCloudProvider(opts Options) Provider {
	switch opts.CloudProviderName {
	case cloudprovider.AwsProviderName:
		return aws.Build(opts)
	case cloudprovider.ClusterAPIProviderName:
		return clusterapi.BuildClusterAPI(opts)
	}
	return nil
}
`)

	changed, err := transformAutoscalerBuilder(source)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("transformAutoscalerBuilder() changed = false")
	}
	content := readAutoscalerTestFile(t, builder, "builder_all.go")
	if strings.Contains(content, "cloudprovider/aws") || strings.Contains(content, "AwsProviderName") {
		t.Fatal("builder still contains AWS provider")
	}
	if !strings.Contains(content, "ClusterAPIProviderName") {
		t.Fatal("builder does not retain Cluster API provider")
	}
	if _, err := os.Stat(filepath.Join(builder, "builder_aws.go")); !os.IsNotExist(err) {
		t.Fatalf("builder_aws.go still exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(builder, "builder_clusterapi.go")); err != nil {
		t.Fatalf("builder_clusterapi.go missing: %v", err)
	}
}

func TestRemoveAutoscalerGCEDependencies(t *testing.T) {
	source := t.TempDir()
	config := filepath.Join(source, "cluster-autoscaler", "config")
	if err := os.MkdirAll(config, 0o755); err != nil {
		t.Fatal(err)
	}
	writeAutoscalerTestFile(t, config, "autoscaling_options.go", `package config
import gce_localssdsize "example/localssdsize"
type GCEOptions struct {
	// LocalSSDDiskSizeProvider provides sizes.
	LocalSSDDiskSizeProvider gce_localssdsize.Provider
}
`)

	changed, err := removeAutoscalerGCEDependencies(source)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("removeAutoscalerGCEDependencies() changed = false")
	}
	content := readAutoscalerTestFile(t, config, "autoscaling_options.go")
	if strings.Contains(content, "localssdsize") || strings.Contains(content, "LocalSSDDiskSizeProvider") {
		t.Fatal("GCE dependencies remain")
	}
}

func TestRemoveAutoscalerCloudProviderDirectoriesKeepsRuntimeSupport(t *testing.T) {
	source := t.TempDir()
	cloudProvider := filepath.Join(source, "cluster-autoscaler", "cloudprovider")
	for _, name := range []string{"builder", "router", "clusterapi", "mocks", "test", "aws", "gce"} {
		if err := os.MkdirAll(filepath.Join(cloudProvider, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	if err := removeAutoscalerCloudProviderDirectories(source); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"builder", "router", "clusterapi", "mocks", "test"} {
		if _, err := os.Stat(filepath.Join(cloudProvider, name)); err != nil {
			t.Fatalf("%s missing: %v", name, err)
		}
	}
	for _, name := range []string{"aws", "gce"} {
		if _, err := os.Stat(filepath.Join(cloudProvider, name)); !os.IsNotExist(err) {
			t.Fatalf("%s still exists: %v", name, err)
		}
	}
}

func writeAutoscalerTestFile(t *testing.T, directory, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readAutoscalerTestFile(t *testing.T, directory, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(directory, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
