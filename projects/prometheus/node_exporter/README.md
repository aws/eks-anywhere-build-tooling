## **Prometheus Node Exporter**
![Version](https://img.shields.io/badge/version-v1.5.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVk9qbzVQdVlyQmVSNE44amtYY0U0YVJDM25yWnJjQlExd25ycDZQWnU1czlVMGt5M2hWMDBSaWlSL1JVU0cwMXBQeUIzczlkWkRZWVhleUpBWFdkOUY4PSIsIml2UGFyYW1ldGVyU3BlYyI6Im1nbzJUbTE1ZUN5SmowN2EiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Prometheus Node Exporter](https://github.com/prometheus/node_exporter) is a Prometheus exporter that exposes a wide variety of hardware and OS metrics. It directly monitors and scrapes metrics from the host machines.

### Updates

#### Version changes
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Run `make generate` from the root of the repo to update the `UPSTREAM_PROJECTS.yaml` file.

#### Make target changes
1. Run `make add-generated-help-block` from the project root to update available make targets.
