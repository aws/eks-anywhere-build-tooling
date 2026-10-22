## **AWS Distro for OpenTelemetry Collector**
![Version](https://img.shields.io/badge/version-v0.50.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiMkdVcDNnUnZnd3NUNE4xeEtERUdyNnpRclN6aXdsbWZhaXdtL1dJYkVRNlJlWVZlMUlGSFlKbHVxTXZIMWgzTUdNWW1kU3FiSHI3ZFY0eDIrcXFNQnNZPSIsIml2UGFyYW1ldGVyU3BlYyI6IlFzaGhYckdMcVI3eTlYTksiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

AWS Distro for OpenTelemetry Collector ([ADOT Collector](https://github.com/aws-observability/aws-otel-collector)) is an AWS supported version of the upstream OpenTelemetry Collector and is distributed by Amazon. It supports the selected components from the OpenTelemetry community.

### Periodic Reviews
Review [image releases](https://github.com/aws-observability/aws-otel-collector/releases) and [helm chart releases](https://github.com/open-telemetry/opentelemetry-helm-charts/releases) periodically to identify new releases and decide on an update plan and an update schedule.

### Updates

#### Version changes
1. Update the `GIT_TAG` and `HELM_GIT_TAG` files to have the new desired version based on the upstream release tags.
   For `HELM_GIT_TAG`, use the latest `opentelemetry-collector` chart whose `appVersion` matches the OpenTelemetry Collector version that the ADOT release is built on (listed in the ADOT release notes).
1. Run `make update-digests` from the project root to update all the `IMAGE_DIGEST` files under `images` directory.
1. Sync `helm/schema.json` with the chart's `values.schema.json`, but **only add properties and presets whose collector components are included in ADOT**.
   ADOT ships a subset of OpenTelemetry components, and the upstream chart's presets assume the full build.
   Presets that render missing components (for example `hostMetrics`, `kubeletMetrics`, `clusterMetrics`, `kubernetesEvents`, `kubernetesObjects`, `annotationDiscovery` and `profiling`) pass schema validation but make the collector exit at startup with `unknown type`.
   To list the supported receivers, run the image with a config that uses an unknown receiver; the error prints the valid values:
   ```bash
   printf 'receivers:\n  bogus: {}\nexporters:\n  debug: {}\nservice:\n  pipelines:\n    metrics:\n      receivers: [bogus]\n      exporters: [debug]\n' > /tmp/adot-probe.yaml
   docker run --rm -v /tmp/adot-probe.yaml:/probe.yaml:ro public.ecr.aws/aws-observability/aws-otel-collector:$(cat GIT_TAG) --config=/probe.yaml
   ```
   The full component list is in [`pkg/defaultcomponents/defaults.go`](https://github.com/aws-observability/aws-otel-collector/blob/main/pkg/defaultcomponents/defaults.go) for the release tag.
1. Keep the root `required` list in `helm/schema.json` to fields that have no chart default (today only `mode`).
   The packages webhook validates the user's config alone, without the chart's `values.yaml` defaults, so copying upstream's `required` (for example `apiVersion`) rejects every config that doesn't set the field.
1. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.
1. Update the version at the top of this `README`.

To make changes to the patches folder, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#generate-patch-files)


To test the upgrade, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#Testing).

#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.
