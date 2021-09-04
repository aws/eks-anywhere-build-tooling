## **EKS-A CLI tools image**

![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiSFNPUVBIWnk4YnplQllqYzZQUGRMKzE4c0sxTEVqY2wrM3ZrYjZickJBbDcwTzJTSmp3d0ZIZDV4Y0Z0QnpaMmFqL1FuS1BNbGdieVp2NGdVeE1VTnowPSIsIml2UGFyYW1ldGVyU3BlYyI6ImowTGVqR3dIeDYwY251TVIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The EKS-A CLI tools image is packaged with the executables that are invoked by the `eks-a` command-line tool, such as `clusterctl`, `kubectl`, `kind`, `flux`, `govc`, etc. This image serves as the runtime environment when using the CLI, but customers can choose to use their own local binaries by setting the flag `MR_TOOLS_DISABLED` to `true`.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/cli-tools).
