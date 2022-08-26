## **AWS Distro for OpenTelemetry Collector**

AWS Distro for OpenTelemetry Collector ([ADOT Collector](https://github.com/aws-observability/aws-otel-collector)) is an AWS supported version of the upstream OpenTelemetry Collector and is distributed by Amazon. It supports the selected components from the OpenTelemetry community.

### Periodic Reviews
Review [image releases](https://github.com/aws-observability/aws-otel-collector/releases) and [helm chart releases](https://github.com/open-telemetry/opentelemetry-helm-charts/releases) periodically to identify new releases and decide on an update plan and an update schedule.

### Updates

#### Version changes
1. Update the `GIT_TAG` and `HELM_GIT_TAG` files to have the new desired version based on the upstream release tags.
1. Run `make update-digests` from the project root to update all the `IMAGE_DIGEST` files under `images` directory.
1. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.

#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.