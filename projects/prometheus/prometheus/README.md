![Version](https://img.shields.io/badge/version-v2.43.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiTldscmdZWkd6NzlhUHJBbFJDRzlMc3NmaGxBOFJlYWE1a3BsVG9KcXhldDRCK05PL0lxNmVVUi9odlMzdXZCYXFxWTBCOUZDbS91R21KL1c5VkdQQ004PSIsIml2UGFyYW1ldGVyU3BlYyI6Im94dGM3UFc0MGRDN0pyREIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

## **Prometheus**

[Prometheus](https://github.com/prometheus/prometheus) is an open-source monitoring solution for collecting and aggregating metrics as time series data.

### Periodic Reviews
Review [image releases](https://github.com/prometheus/prometheus/tags) periodically to identify new releases and decide on an update plan and an update schedule.

### Updates

#### Version changes
1. Update the `GIT_TAG` and `HELM_GIT_TAG` files to have the new desired version based on the upstream release tags.
1. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.

#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.
