## **Linux Boot Config**
![Version](https://img.shields.io/badge/version-v5.17-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQzVkSWVDSHVIYkNHSWdxMTVwMkI3RnFXVEMyeWdJbWplOUdUbFcxbkpEd3Zld05LaDhUYlRIK09xaFhTQ3doMVM2U3doWkpkOCt4NmVpYi85VzV1SnNNPSIsIml2UGFyYW1ldGVyU3BlYyI6Ik9ySnBzTE9LMnRwRWNhekgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Linux Boot Config](https://github.com/torvalds/linux/tree/master/tools/bootconfig) project is a tool on the linux repository that helps expand 
the kernel command line to support additional key-value data when booting the kernel in an efficient way. This allows administrators to pass a structured key config file

This tool is used by Bottlerocket metal variant to help parse the user input of kernel parameters into an initrd binary file that the kernel understands.

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/torvalds/linux/tree/master) and decide on the new version
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @vignesh-goutham, @zmrow or @jaxesn
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. If any of the make targets used in projects/torvalds/linux/Makefile to call upstream make changed, make those appropriate changes.
1. Update the version at the top of this Readme.
1. Monitor node image builds and e2e tests as updates to this project can potentially break cluster create and update process.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.