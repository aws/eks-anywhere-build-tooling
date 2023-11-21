# **Tinkerbell Stack Helm Chart**
![Version](https://img.shields.io/badge/version-0.2.4-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoic0w3TWw2ZDdFblpMMDZtamh2S3RiNmVwbTdRaDVlbmgxWGZkVi9WdGZjMDgvL2J2a1ZGSXJoMVV2dWJlNWZpbjV5Z3k4THRjZ0VyWUlBM0RLTUNWaE4wPSIsIml2UGFyYW1ldGVyU3BlYyI6IllIWGJ1SDFRZm1HM0dnK1giLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

## Updating

1. Make your changes to the chart.
1. When copying CRDs from upstream, make sure CRDs that will/need to be moved around by `clusterctl move` have the appropriate labels.

    ```yaml
    labels:
      clusterctl.cluster.x-k8s.io: ""
      clusterctl.cluster.x-k8s.io/move: ""
    ```

1. Bump the version, appropriately, in the following files: `GIT_TAG`, `helm/sedfile.template`, and `chart/Chart.yaml`.
1. Run `make verify` to verify the chart.
1. Commit and push your changes.
