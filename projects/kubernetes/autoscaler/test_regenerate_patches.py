import importlib.util
import tempfile
import unittest
from pathlib import Path


MODULE_PATH = Path(__file__).with_name("regenerate_patches.py")
SPEC = importlib.util.spec_from_file_location("autoscaler_regenerator", MODULE_PATH)
regenerator = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(regenerator)


class AutoscalerRegeneratorTest(unittest.TestCase):
    def test_transform_router_keeps_only_clusterapi(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            source = Path(temp_dir)
            router = source / "cluster-autoscaler" / "cloudprovider" / "router"
            router.mkdir(parents=True)
            (router / "router_aws.go").write_text("package router\n", encoding="utf-8")
            (router / "router_clusterapi.go").write_text("package router\n", encoding="utf-8")
            (router / "router_all.go").write_text(
                """package router
import (
\t_ "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
\t_ "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/clusterapi"
)
func init() {
\tbuilder.SetDefaultCloudProvider(cloudprovider.GceProviderName)
}
""",
                encoding="utf-8",
            )

            self.assertTrue(regenerator.transform_router(source))
            content = (router / "router_all.go").read_text(encoding="utf-8")
            self.assertNotIn("cloudprovider/aws", content)
            self.assertIn("cloudprovider/clusterapi", content)
            self.assertIn("cloudprovider.ClusterAPIProviderName", content)
            self.assertFalse((router / "router_aws.go").exists())
            self.assertTrue((router / "router_clusterapi.go").exists())

    def test_transform_builder_keeps_only_clusterapi(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            source = Path(temp_dir)
            builder = source / "cluster-autoscaler" / "cloudprovider" / "builder"
            builder.mkdir(parents=True)
            (builder / "builder_aws.go").write_text("package builder\n", encoding="utf-8")
            (builder / "builder_clusterapi.go").write_text("package builder\n", encoding="utf-8")
            (builder / "builder_all.go").write_text(
                """package builder
import (
\t"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
\t"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
\t"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/clusterapi"
)
var AvailableCloudProviders = []string{
\tcloudprovider.AwsProviderName,
\tcloudprovider.ClusterAPIProviderName,
}
const DefaultCloudProvider = cloudprovider.GceProviderName
func buildCloudProvider(opts Options) Provider {
\tswitch opts.CloudProviderName {
\tcase cloudprovider.AwsProviderName:
\t\treturn aws.Build(opts)
\tcase cloudprovider.ClusterAPIProviderName:
\t\treturn clusterapi.BuildClusterAPI(opts)
\t}
\treturn nil
}
""",
                encoding="utf-8",
            )

            self.assertTrue(regenerator.transform_builder(source))
            content = (builder / "builder_all.go").read_text(encoding="utf-8")
            self.assertNotIn("cloudprovider/aws", content)
            self.assertNotIn("AwsProviderName", content)
            self.assertIn("ClusterAPIProviderName", content)
            self.assertFalse((builder / "builder_aws.go").exists())
            self.assertTrue((builder / "builder_clusterapi.go").exists())

    def test_remove_gce_dependencies(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            source = Path(temp_dir)
            path = source / "cluster-autoscaler" / "config" / "autoscaling_options.go"
            path.parent.mkdir(parents=True)
            path.write_text(
                """package config
import gce_localssdsize "example/localssdsize"
type GCEOptions struct {
\t// LocalSSDDiskSizeProvider provides sizes.
\tLocalSSDDiskSizeProvider gce_localssdsize.Provider
}
""",
                encoding="utf-8",
            )

            self.assertTrue(regenerator.remove_gce_dependencies(source))
            content = path.read_text(encoding="utf-8")
            self.assertNotIn("localssdsize", content)
            self.assertNotIn("LocalSSDDiskSizeProvider", content)

    def test_remove_cloud_provider_directories_keeps_runtime_support(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            source = Path(temp_dir)
            cloudprovider = source / "cluster-autoscaler" / "cloudprovider"
            for name in ("builder", "router", "clusterapi", "mocks", "test", "aws", "gce"):
                (cloudprovider / name).mkdir(parents=True)

            regenerator.remove_cloud_provider_directories(source)

            for name in ("builder", "router", "clusterapi", "mocks", "test"):
                self.assertTrue((cloudprovider / name).exists())
            for name in ("aws", "gce"):
                self.assertFalse((cloudprovider / name).exists())


if __name__ == "__main__":
    unittest.main()
