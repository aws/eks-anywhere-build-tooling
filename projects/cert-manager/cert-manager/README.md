## **cert-manager**
![Version](https://img.shields.io/badge/version-v1.19.3-blue)

([cert-manager](https://github.com/cert-manager/cert-manager)) is an AWS supported version of the upstream cert-manager and is distributed by Amazon EKS add-ons.

### Periodic Reviews
Review [helm chart releases](https://github.com/cert-manager/cert-manager/releases) periodically to identify new releases and decide on an update plan and an update schedule.

### Updating

1.  To update the EKS_ADDON_IMAGE_TAG, check the latest available version from the EKS-Addon team using the following AWS CLI command:
    ```bash
    aws eks describe-addon-versions \
    --kubernetes-version 1.35 \
    --addon-name cert-manager \
    --query 'addons[0].addonVersions[0].addonVersion' \
    --output text
    ```
2. For updating HELM_GIT_TAG, monitor [upstream releases](https://github.com/cert-manager/cert-manager/releases) and changelogs and when to bump the tag.

### Notes
- The startupapicheck component is disabled because the EKS Add-on team does not publish a startupapicheck image.
- Images are sourced from the EKS Add-on ECR repo (602401143452.dkr.ecr.us-west-2.amazonaws.com) under the `eks/` prefix.
