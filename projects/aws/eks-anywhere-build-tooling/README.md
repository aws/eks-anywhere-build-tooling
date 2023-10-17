## **EKS-A CLI tools image**
![Version](https://img.shields.io/badge/version-v0.17.4-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiNVpyOFVBOHNqZkE0OEVta1Q1Z2xlSytId0l2NTNYYUNXRzdoL2xsV2N5cWlzUDErZjRvQm42ZGRLeWQ2TzQ2eGtEM3l0Z21pZksxbGczTG90YzFuR3J3PSIsIml2UGFyYW1ldGVyU3BlYyI6IkRDeENUYkFXQk53MUNTYVYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The EKS-A CLI tools image is packaged with the executables that are invoked by the `eks-a` command-line tool, such as `clusterctl`, `kubectl`, `kind`, `flux`, `govc`, etc. This image serves as the runtime environment when using the CLI, but customers can choose to use their own local binaries by setting the flag `MR_TOOLS_DISABLED` to `true`.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/cli-tools).

## Building
This project depends on other artifacts from this repo.  To build image locally, `ARTIFACTS_BUCKET` must be supplied. For ex
the following is the presubmit bucket:

`ARTIFACTS_BUCKET=s3://projectbuildpipeline-857-pipelineoutputartifactsb-10ajmk30khe3f make build`
