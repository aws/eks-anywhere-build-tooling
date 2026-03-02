![Version](https://img.shields.io/badge/version-v0.310.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiTldscmdZWkd6NzlhUHJBbFJDRzlMc3NmaGxBOFJlYWE1a3BsVG9KcXhldDRCK05PL0lxNmVVUi9odlMzdXZCYXFxWTBCOUZDbS91R21KL1c5VkdQQ004PSIsIml2UGFyYW1ldGVyU3BlYyI6Im94dGM3UFc0MGRDN0pyREIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

## **Prometheus**

[Prometheus](https://github.com/prometheus/prometheus) is an open-source monitoring solution for collecting and aggregating metrics as time series data.

### Periodic Reviews
Review [image releases](https://github.com/prometheus/prometheus/tags) and [chart releases](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus) periodically to identify new releases and decide on an update plan and an update schedule.

### Updates

#### Version changes
1. Update the `GIT_TAG`, `HELM_GIT_TAG`, and `GOLANG_VERSION` files to have the new desired version based on the upstream release tags.
1. Review the patches under `patches/` folder and remove any that are either merged upstream or no longer needed.
1. Current patch information:
    * `helm/patches`:
        1. 0001 patch removes prometheus chart dependencies.
		1. 0002 patch adds node exporter component.
		1. 0003 patch adds changes to update image repo
		1. 0004 patch adds changes to update namespace
		1. 0005 patch add changes for pod update strategy due to config map changes
		1. 0006 patch update values.yaml and adds values.schema.json
1. Run `make build` or `make release` to build package, if `apply patch` step fails during build follow the  steps below to update the patch and rerun build/release again.
1. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.
1. Update the version at the top of this `README`.


To make changes to the patches folder, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#generate-patch-files)


To test the upgrade, follow the steps mentioned [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/packages/update-helm-charts.md#Testing).

#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.
