## **Image Builder**
![Version](https://img.shields.io/badge/version-105b63a6b281f55026627711a4ff32651a944c55-blue)
| Artifact | Build Status |
| --- | --- |
| 1-20 OVA | ![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiazF2R2J2ell0Y0tHT1RnYmF6WXdnRjMwTHMyMTlSVXZnMVoyRytWZ0FDaE5HOU5WejA2VjFzSVNObWlXTjM0eHh2akpBbjgwV0xaTjl5cjFOZlFrZlNNPSIsIml2UGFyYW1ldGVyU3BlYyI6IjhxZTMzVVhZZnR6V0JBOU4iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |
| 1-21 OVA | ![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibHJVYmMvSUF0ZlkrVEJsMVBwZU9xLy9ndUZ0U3dGZStpelk2RDRpRTBLQnBrQWNqVkU2TW9qWWI1aFBJM1hpQ1B6TzhaeVduTWdxcE5JeS9XWGhDME5RPSIsIml2UGFyYW1ldGVyU3BlYyI6InVGUS9yandMWmd1cWRsOWciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |
| 1-20 Raw | ![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiSGszQi9CcGZTNDgzQjZzUkpQMTNjaEc3UmdvTk45YlpqLzc1OVo2dVdFc2xNSVN1bjV3dGtzK2dyTW5wSzBPWU41ZjZZZ0Ztb0E5eGVVaDRPOG1tdW5ZPSIsIml2UGFyYW1ldGVyU3BlYyI6IlQxbHdMQ0Q0QkNmWURkL3QiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |
| 1-21 Raw | ![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiNVVkU201eEZoUUQ4RVQ5eE5Hang3cW84Vzl0OHdRYVBpbnNFK0dLVXBqeGZjMnRtZERtUTBBNHlVZFlyVEZPbTRLZ3VkRTRQT3J1WDBsUmV0Q0RQWHp3PSIsIml2UGFyYW1ldGVyU3BlYyI6Imcvb0VFWmNVdDYxK01oZEIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main) |

The [Image Builder project](https://github.com/kubernetes-sigs/image-builder) offers a collection of cross-provider Kubernetes virtual machine image building utilities. It can be used to build images intended for use with Kubernetes Cluster API providers. Each provider has its own format of images that it can work with, for example, AMIs for AWS instances, and OVAs for vSphere. The Image Builder project relies on Packer configuration files and Ansible playbooks to build the images and store them in appropriate locations and accounts.

EKS-A CLI project uses these images as the node image when constructing workload clusters for different infrastructure providers like AWS and vSphere.

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/kubernetes-sigs/image-builder) and decide on new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @vignesh-goutham or @jaxesn.
1. Follow these steps for changes to the patches/ folder:
    1. Checkout the desired tag on upstream [repo](https://github.com/kubernetes-sigs/image-builder) and create a new branch on your local workspace.
    1. Review the patches under patches/ folder in this repo. Apply the required patches to the new branch in the local clone of upstream repo created in the above step.
        1. Run `git am <path to patches>` on the upstream clone.
        1. For patches that need some manual changes, you will see a similar error: `Patch failed at *`
        1. For that patch, run `git apply --reject --whitespace=fix *.patch`. This will apply hunks of the patch that do apply correctly, leaving
           the failing parts in a new file ending in `.rej`. This file shows what changes weren't applied and you need to manually apply.
        1. Once the changes are done, delete the `.rej` file and run `git add .` and `git am --continue`
    1. Remove any patches that are either merged upstream or no longer needed. Please reach out to @vignesh-goutham or @jaxesn if there are any questions regarding keeping/removing patches.
    1. Run `git format-patch <commit>`, where `<commit>` is the last upstream commit on that tag. Move the generated patches from under the upstream fork to the projects/kubernetes-sigs/image-builder/patches/ folder in this repo.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. If any of the make targets used in projects/kubernetes-sigs/image-builder/Makefile to call upstream make changed, make those appropriate changes.
1. Update the version at the top of this Readme.
1. Monitor node image builds and e2e tests as updates to this project can potentially break cluster create and update process.

### Building your own image

Users can build their own node images from a custom base image using this image-building process. The generated image can be used 
with EKS-A CLI to create clusters.

Currently this process only supports building RHEL images for KVM hypervisor.

#### Prerequisites
- Builder machine running Ubuntu 20.04 on metal

#### Steps
- Clone this repository on the builder machine `git clone https://github.com/aws/eks-anywhere-build-tooling.git`
- Change directory into image-builder project at `eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder`
- Download the base image to the current working path
- Setup KVM on the builder machine
  ```
  sudo apt update -y
  sudo apt install qemu-kvm libvirt-daemon-system libvirt-clients virtinst cpu-checker libguestfs-tools libosinfo-bin unzip ansible -y
  sudo usermod -a -G kvm ubuntu
  sudo chown root:kvm /dev/kvm
  sudo snap install yq
  sudo apt install jq
  ```
- Setup exports
  * export BASE_IMAGE=<file name of base image>
  * export RELEASE_BRANCH=<1-21 or 1-20>
  * export RHSM_USER=<RedHat username>
  * export RHSM_PASS=<RedHat password>
  * export ARTIFACTS_BUCKET=s3://projectbuildpipeline-857-pipelineoutputartifactsb-10ajmk30khe3f
- Run the local build make command. This command will verify required build dependencies and run the image building process.
    ```
    make build-qemu-rhel-local
    ```