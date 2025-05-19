## **CRI Tools**
![Version](https://img.shields.io/badge/version-v1.33.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiUUlRZXJEVUxWcjI1OE8weVdXQnY4alBSU1lxVm1FOGVoZE83VldDbjJiaFBtY25XT3NIK1RhckZkQXZGclZDSkVLUG5PMmd5K2J2RVlSYk9pclUybC9zPSIsIml2UGFyYW1ldGVyU3BlYyI6IkF3RGUzVDFhVlB0eUlGMWwiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [CRI tools project](https://github.com/kubernetes-sigs/cri-tools) provides a CLI and validation tools for the `kubelet`'s Container Runtime Interface (CRI). This allows CRI runtime developers to debug their runtimes (like `containerd`, `CRI-O`, etc.) without needing to set up Kubernetes components. The `crictl` CLI can perform numerous functions such as running containers, fetching logs, listing conatiner stats, removing images, etc.
