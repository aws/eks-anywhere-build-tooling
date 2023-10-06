# **Tinkerbell Stack Helm Chart**

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
