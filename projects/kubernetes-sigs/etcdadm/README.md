## **etcdadm**
![Version](https://img.shields.io/badge/version-f089d308442c18f487a52d09fd067ae9ac7cd8f2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiK0pzWGNJc01qaEVYTU9JcjY5MzdFTFVlSmV2aE1ESUVlODhKNHErSUNJSlkrV1o2bDlPS1hRU1BsWGJhNTZEVkNEYXVGeGRpRnJ4VkpjdFNiR2ZVQ21nPSIsIml2UGFyYW1ldGVyU3BlYyI6Ikh6dkhlYVh0QnE1TytCaU0iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[etcdadm](https://github.com/kubernetes-sigs/etcdadm) is a command-line tool for operating an etcd cluster. It downloads a specific etcd release, installs the binary, configures a systemd service, generates CA certificates, calls the etcd API to add (or remove) a member, and verifies that the new member is healthy. Its user experience is inspired by kubeadm.

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/kubernetes-sigs/etcdadm) and decide on new version. 
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @g-gaston.
1. Follow these steps for changes to the patches/ folder:
    1. Checkout the upstream repo and create a new branch from the desired upstream release tag.
    1. Review the patches under patches/ folder in this repo. Apply the required patches to the new branch created in the above step. Remove any patches that are either
    merged upstream or no longer needed. Please reach out to @g-gaston if there are any questions regarding keeping/removing patches.
    1. Run `git format-patch <commit>`, where `<commit>` is the last upstream commit on that tag. Move the generated patches under the patches/ folder in this repo.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [0.13.2 compared to 0.23.0](https://github.com/kubernetes-sigs/etcdadm/compare/v0.1.3...v0.1.5). Check the `$(BIN)` target for
any build flag changes, tag changes, dependencies, etc. 
1. Verify the golang version has not changed. The version specified in the variable `GO_IMAGE` in the Makefile seems to be kept up to date.
1. Update checksums and attribution using `make run-attribution-checksums-in-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
1. Update the `ETCDADM_VERSION` to the upstream tag [here](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/kubernetes-sigs/image-builder/build/setup_packer_configs.sh#L76).
