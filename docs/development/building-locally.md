# Building locally

All projects in this repo follow a common building (and updating) process.  Make targets are standardized across the repo.
Projects can be built on a Mac or Linux host. By default when running `go build`, and some of the other key targerts
such as `gather_licenses` and `attribution`, the builder-base docker image used in prow will be used.
This is to ensure checksums and attributions match what is generated in our CI and release systems.
To force running targets on the host: `export FORCE_GO_BUILD_ON_HOST=true`.


## Pre-requisites for building on MacOS host
* Docker Desktop - https://docs.docker.com/desktop/install/mac-install/
* Common utils required - `brew install jq yq`
* Bash 4.2+ is required - `breq install bash`
* A number of targets which run on the host require the gnu variants of common tools
	such as `date` `find` `sed` `stat` `tar` - `brew install coreutils findutils gnu-sed gnu-tar`
* Building container images requires `buildctl` - `brew buildkit`
* Building helm charts requires `helm` and `skopeo`
	* https://helm.sh/docs/intro/install/ - `curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash`
	* https://github.com/containers/skopeo - `brew install skopeo`
* To make pushing to ECR easier: `brew install docker-credential-helper-ecr`

## Pre-requisites for building on Linux (AL2/AL23) host
* Docker
	* `sudo yum install -y docker && sudo usermod -a -G docker ec2-user && id ec2-user && newgrp docker`
	* `sudo systemctl enable docker.service && sudo systemctl start docker.service`
* Common utils required 
	* `sudo yum -y install jq git-core make lz4`
* `yq` is required by a number of targets
	* `sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_$([ "x86_64" = "$(uname -m)" ] && echo amd64 || echo arm64) && sudo chmod +x /usr/local/bin/yq`
* Building helm charts requires `helm` and `skopeo`
	* https://helm.sh/docs/intro/install/ - `curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash`
	* https://github.com/containers/skopeo - There is no binary release of skopeo for linux see [installing](https://github.com/containers/skopeo/blob/main/install.md) or
		copy from builder-base image to host: `CONTAINER_ID=$(docker run -d -i public.ecr.aws/eks-distro-build-tooling/builder-base:standard-latest.al2 bash) && sudo docker cp $CONTAINER_ID:/usr/bin/skopeo /usr/local/bin && docker rm -vf $CONTAINER_ID`
* Building container images for different architectures will require `qemu`
	* `curl -Ls https://github.com/tonistiigi/binfmt/releases/download/deploy%2Fv6.2.0-26/qemu_v6.2.0_linux-amd64.tar.gz | sudo tar xzvf - -C /usr/bin`
	* `docker run --privileged --rm public.ecr.aws/eks-distro-build-tooling/binfmt-misc:qemu-v6.1.0 --install aarch64,amd64`
* To ensure string sorting matches our build pipelines - `export LANG=C.UTF-8`
* To make pushing to ECR easier: `sudo yum install amazon-ecr-credential-helper -y`
* (Skip for AL23) Upgrade to aws cli v2: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
* Building image-builder
	* `tuftool` to download Bottlerocket images - https://github.com/awslabs/tough - There is no binary release of tuftool for linux see [Readme](https://github.com/awslabs/tough/blob/develop/tuftool/README.md) or
		copy from builder-base image to host: `CONTAINER_ID=$(docker run -d -i public.ecr.aws/eks-distro-build-tooling/builder-base:standard-latest.al2 bash) && sudo docker cp $CONTAINER_ID:/usr/bin/tuftool /usr/local/bin && docker rm -vf $CONTAINER_ID`
	* (Skip for AL23) python 3.9 is required for the version of ansible required. AL2 does not ship a python 3.9 package so the easiest way to install it is building it with pyenv:
		* https://github.com/pyenv/pyenv
		* `curl https://pyenv.run | bash`
		* `sudo yum install gcc make patch zlib-devel bzip2 bzip2-devel readline-devel sqlite sqlite-devel openssl11-devel tk-devel libffi-devel xz-devel -y`
		* `pyenv install 3.9 && pyenv global 3.9`


## Typical build targets

Each project folder contains a Makefile with the following build targets.

* `binaries` - clones repo, checks out git tag and builds binaries.
* `local-images` - builds container images and writes to `/tmp/{project}.tar`.
* `images` - builds container images apnd pushes to configured registry. For a few projects (coredns, etcd, kubernetes) local oci tars are also built
	to be uploaded to s3.
* `gather-licenses` - copies licenses from all project dependencies and puts in `_output/LICENSES` to be included in tarballs and container images.
* `attribution` - regenerate the ATTRIBUTION.txt file, golang versions are important here so make sure you have the same versions as the builder-base.
	There is a periodic job which runs every day to keep these up to date, builds will not fail due to these not being 100% accurate.
* `checksums` - update CHECKSUMS file based on currently built binaries.  Builds should be reproducible across Mac and Linux hosts.
	This should be run when bumping GIT_TAG, changing patches or changing build flags, otherwise should not be changed. 
	Builds will fail if these do not match, but the correct checksums will be outputted to make updating easier.
* `build` - run via prow presubmit, clones repo, build binary(s), gather licenses, build container images, build tarballs and generate attribution file.
	Uses `local-images` to build images which does not push to a registry.
* `release` - run via prow postsubmit, same as `build` except runs `images` to push to a remote registry and to push tars/binaries
	to s3.


## Pre-requisites for building container images

* There are two options for building container images
	* `buildctl` can be used instead of docker to match our prow builds.  Running `local-images` targets 
	will export images to a tar, but if running `images` targets which push images, a registry is required.  By default,
	an ECR repo in the currently logged in AWS account will be used.  A common alternative is to run docker registry locally and override
	this default behavior. This can be done with `IMAGE_REPO=localhost:5000 make images`. To run buildkit and the docker registry run from the repo root (or a specific project folder):
		* `curl -L https://github.com/moby/buildkit/releases/download/v0.12.2/buildkit-v0.12.2.linux-$([ "x86_64" = "$(uname -m)" ] && echo amd64 || echo arm64).tar.gz -o /tmp/buildkit.tgz && sudo tar -xf /tmp/buildkit.tgz -C /usr/local bin/buildctl`
		* `make run-buildkit-and-registry`
		* `export BUILDKIT_HOST=docker-container://buildkitd`
		* `make stop-buildkit-and-registry` to cleanup
	* (Experimental) `docker buildx` can be used which may be easier to setup since docker now ships with the buildx plugin.
	docker buildx also uses buildkit behind the scenes so it should result in creating an image which is the same as using `buildctl` directly as prow does.
		* You need to create a builder using either the `docker-container` or `remote` driver. For example: `docker buildx create --name multiarch --driver docker-container --use`
* See [docker-auth.md](./docker-auth.md) for helpful tips for configuring, securing, and using your AWS credentials with the targets in this project.
* By default, `IMAGE_REPO` will be set to your AWS account's private ECR registry based on your AWS_REGION. If you want to create the neccessary registeries
for the current project:
	* `make create-ecr-repos`
	* If you would like to create these repos in your ECR public registry - `IMAGE_REPO=public.ecr.aws/<id> create-ecr-repos`
* Some projects download built artifacts from our dev pipelines. To pull these from the EKS-Anywere dev bucket
	* `export ARTIFACTS_BUCKET=s3://projectbuildpipeline-857-pipelineoutputartifactsb-10ajmk30khe3f`

## (Not recommended) Running go builds on the host

* Multiple versions of golang are required to build the various projects.  There is a helper 
[script](../../build/lib/install_go_versions.sh) modeled from the builder-base to install all golang versions locally.  
**Note** Our build and release pipelines use [EKS Go](https://github.com/aws/eks-distro-build-tooling/blob/main/projects/golang/go/README.md) which
the helper script does not install from. If building on the host, checksums will almost certainly mismatch with what is checked into this repo.
* If testing gathering licenses and generating attribution files on the host, please follow these [instructions](attribution-files.md) to setup.


## Create a long running in builder-base container

By default targets are run in a short-lived container using the builder-base image. If you would like to create a long lived container for testing:

* `make start-docker-builder`
* The container will be named `eks-a-builder` and can be accessed to manually run targets or cleanup:
	* `docker exec -it eks-a-builder bash`
* To stop the container:
	* `make stop-docker-builder`

## Other considerations

* If you use a private `GOPROXY`, add `export GOPROXY=<custom proxy>` to your ~/.bashrc or ~/.zshrc
	* If your private proxy requires authentication, consider creating `~/.netrc`, ex:
		``` 
			machine <machine>
			login <basic auth name>
			password <basic auth password>
		```
		* `chmod og-rw  ~/.netrc`
