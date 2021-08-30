## **Source Controller**
![Version](https://img.shields.io/badge/version-v0.12.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiS1ZJY3BFVGg0a21PUmpDVWM2T0pnc2VxV25uYWt5aGJjQktVSURIVnBsd0VBUmljSlUxTVNyeG5pSzhFbXNaMkdiUGdBRWU5L2plMG9ldVFxcHhrYjd3PSIsIml2UGFyYW1ldGVyU3BlYyI6IjgybDlDK2ZHLzJQVmNZNFoiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [source-controller](https://github.com/fluxcd/source-controller) is a Kubernetes operator specialized in artifacts acquisition from external sources such as Git, Helm repositories and S3 buckets. The controller watches for `Source` objects in a cluster and acts on them. It was designed with the goal of offloading the sources' registration, authentication, verification and resource-fetching to a dedicated controller.

Some of the features of the Source controller are:

* Authenticates to sources (SSH, user/password, API token)
* Validates source authenticity (PGP)
* Detects source changes based on update policies (semver)
* Fetches resources on-demand and on-a-schedule
* Packages the fetched resources into a well-known format (tar.gz, yaml)
* Makes the artifacts addressable by their source identifier (SHA, version, ts)
* Makes the artifacts available in-cluster to interested third-parties
* Notifies interested third-parties of source changes and availability (status conditions, events, hooks)
* Reacts to Git push and Helm chart upload events

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/l0g8r8j6/fluxcd/source-controller).