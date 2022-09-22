# Mirror images for EKS-A Packages

Most images used in EKS-A Packages are directly built in eks-anywhere-build-tooling.
When there are existing Amazon-built images available, we may choose to deviate from this approach and mirror upstream images instead.

## Rationales for mirroring images
### Why don't we rebuild every image?

Let's start with why we are rebuilding almost all images. We have a promise to customers that all container images of packages are built from source code by Amazon to ensure quality. As most of the upstream of our packages are not Amazon-owned, we build almost all images in our eks-anywhere-build-tooling to keep this promise.

If the upstream repositories are Amazon-owned, the Amazon-built promise has been fulfilled out of box. There are a few benefits for not rebuilding the images in this case:
- We can create a single source of truth for images, which simplifies and standardizes version controls, feature management, and CVE findings management.
- We can reduce duplicate work with another Amazon team, and exploits the synergy between two teams to bring the best to our customers.

One thing to keep in mind is that we cannot default to use upstream images just because they are Amazon-built. Different Amazon teams have different approaches to software development and releases, and other teams' approach may not be compatible with EKS-A packages. Therefore, using existing Amazon-built images require a careful review of how the images
are built and how the Amazon team that owns the images is maintaining 
them.

### If there are Amazon-built images, why don't we use them directly?

Though we don't rebuild images in eks-anywhere-build-tooling, we still
mirror the upstream images in our package image registry and use them as default for packages. This helps ensure 
all the images we used go through standard image
scanning procedures we designed specifically for packages. If there are any CVE findings, we can be notified right away and take proper actions.

## Standards for mirroring images

### Mirroring with AWS ECR pull-through cache (PTC) rules

[PTC](https://docs.aws.amazon.com/AmazonECR/latest/userguide/pull-through-cache.html) is one of the options to mirror images and is currently adopted for packages (i.e. ADOT). It requires mainly two steps to set up:
- Create a PTC rule in the AWS account that hosts all package images.
  - Use the *\<package-project-name\>* as the ECR repository prefix.
  - Ensure the created ECR repository has the following policy actions:
    ```   
    "ecr:GetDownloadUrlForLayer",
    "ecr:BatchGetImage",
    "ecr:BatchImportUpstreamImage"
    ```
- Modify package helm chart to point to the cached images from the package image registry.

Note in order to make eks-anywhere-build-tooling generate helm charts and associated configs properly, image repository names and helm chart repository names need to follow a strict name convention.
- All images repository created by PTC rules follow the name convention of *\<prefix\>/\<ecr-upstream-repo\>* by default. The prefix needs to be *\<package-project-name\>*.
- All helm chart repository need to follow the name convention *\<package-project-name\>/charts/\<image-component-name*\>.

Refer to [helm_require.sh](https://github.com/aws/eks-anywhere-build-tooling/blob/main/build/lib/helm_require.sh) for how name conventions influence the helm charts generation process.

### Mirroring with Skopeo copy

[Skopeo Copy](https://github.com/containers/skopeo/blob/main/docs/skopeo-copy.1.md) is another option to mirror images, and provide more flexibilities for image naming and tagging than PTC. Note eks-anywhere-build-tooling repo hasn't been set up to work with Skopeo copy out of box, so some modifications to the build process will be needed.
