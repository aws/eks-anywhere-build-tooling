# **Tinkerbell CRDs Helm Chart**
![Version](https://img.shields.io/badge/version-v0.22.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoic0w3TWw2ZDdFblpMMDZtamh2S3RiNmVwbTdRaDVlbmgxWGZkVi9WdGZjMDgvL2J2a1ZGSXJoMVV2dWJlNWZpbjV5Z3k4THRjZ0VyWUlBM0RLTUNWaE4wPSIsIml2UGFyYW1ldGVyU3BlYyI6IllIWGJ1SDFRZm1HM0dnK1giLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

This chart installs Tinkerbell CRDs (Custom Resource Definitions) on Kubernetes.

## How it works

CRDs are automatically generated from the [tinkerbell/tinkerbell](https://github.com/tinkerbell/tinkerbell) mono-repo during build:

1. The build clones the mono-repo at the version specified in `GIT_TAG`
2. `scripts/generate-crds-chart.sh` copies CRDs from `crd/bases/` 
3. Required annotations and labels are added automatically:
   - `helm.sh/resource-policy: keep` - prevents Helm from deleting CRDs on uninstall
   - `clusterctl.cluster.x-k8s.io` labels - enables CAPI move operations

## Updating

1. Update `GIT_TAG` to the new mono-repo version (e.g., `v0.22.1`)
2. Bump the chart version in `chart/Chart.yaml` and `helm/sedfile.template`
3. Update `appVersion` in `chart/Chart.yaml` to match `GIT_TAG`
4. Run `make verify` to verify the chart builds correctly
5. Commit and push your changes

## Files

- `chart/` - Chart metadata (Chart.yaml, values.yaml) - these are copied to the generated chart
- `scripts/generate-crds-chart.sh` - Script that generates the chart from mono-repo CRDs
- `scripts/verify.sh` - Verification script
- `helm/sedfile.template` - Sed file for version substitution
