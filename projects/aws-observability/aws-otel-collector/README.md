## **AWS Distro for OpenTelemetry Collector**

AWS Distro for OpenTelemetry Collector ([ADOT Collector](https://github.com/aws-observability/aws-otel-collector)) is an AWS supported version of the upstream OpenTelemetry Collector and is distributed by Amazon. It supports the selected components from the OpenTelemetry community.

### Updating

1. Review [releases](https://github.com/aws-observability/aws-otel-collector/releases)
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Verify the golang version has not changed. 
1. Verify no changes have been made to the dockerfiles
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.