## **Metrics Server**
![Version](https://img.shields.io/badge/version-v0.8.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiSEFNYVlKSURxN25YRGpuWURwWmZOS05vbkl6YTdHTzNHTFJpdzdHZGJUL001ZlNqS1JhblM0QTl2VytuUzNRQ09WazJwRHVUZnp0dVRCb3dLTUVxb2w4PSIsIml2UGFyYW1ldGVyU3BlYyI6IkJIOGVvTFk2bWVVcnhUTkoiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

([Metrics Server](https://github.com/kubernetes-sigs/metrics-server)) is an AWS supported version of the upstream Metrics Server and is distributed by Amazon EKS add-ons.

### Periodic Reviews
Review [helm chart releases](https://github.com/kubernetes-sigs/metrics-server/releases) periodically to identify new releases and decide on an update plan and an update schedule.

### Updating

1.  To update the GIT_TAG, check the latest available version from the EKS-Addon team using the following AWS CLI command:
    ```bash
    aws eks describe-addon-versions \
    --kubernetes-version 1.32 \
    --addon-name metrics-server \
    --query 'addons[0].addonVersions[0].addonVersion' \
    --output text
    ```
2. For updating HELM_GIT_TAG, monitor [upstream releases](https://github.com/kubernetes-sigs/metrics-server/releases) and changelogs and when to bump the tag.
