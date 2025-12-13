## **Prometheus Node Exporter**
![Version](https://img.shields.io/badge/version-v1.10.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVk9qbzVQdVlyQmVSNE44amtYY0U0YVJDM25yWnJjQlExd25ycDZQWnU1czlVMGt5M2hWMDBSaWlSL1JVU0cwMXBQeUIzczlkWkRZWVhleUpBWFdkOUY4PSIsIml2UGFyYW1ldGVyU3BlYyI6Im1nbzJUbTE1ZUN5SmowN2EiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Prometheus Node Exporter](https://github.com/prometheus/node_exporter) is a Prometheus exporter that exposes a wide variety of hardware and OS metrics. It directly monitors and scrapes metrics from the host machines.

### Updates

#### Version changes
1. Update the `GIT_TAG` and `GOLANG_VERSION` file to have the new desired version based on the upstream release tags.
2. Run `make build` or `make release` to build package, if `apply patch` step fails during build follow the  steps below to update the patch and rerun build/release again.
3. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.
4. Update the version at the top of this `README`.


To make changes to the patches folder, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#generate-patch-files)

To test the upgrade, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#Testing).

#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.
