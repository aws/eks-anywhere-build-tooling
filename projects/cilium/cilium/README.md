## **Cilium**
![Version](https://img.shields.io/badge/version-v1.18.2--0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYTh2UnBFVGhjQ1EyeENsWU91ZlJzMktyZHRINlpFWlc0RkZ5amU3Yy96b3p2Z2dxNThZZVQ5ZjRPTEZndGVNQVMwNkMvVmZZR000bGJXWDFqWDFnUlZVPSIsIml2UGFyYW1ldGVyU3BlYyI6ImZRZ2JzZmhRcWZtNFNHZTciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Cilium](https://github.com/cilium/cilium) is an open-source software for providing and transparently securing network connectivity and load-balancing between application workloads such as application containers or processes. Cilium operates at Layer 3/4 to provide traditional networking and security services as well as at Layer 7 to protect and secure use of modern application protocols such as HTTP, gRPC and Kafka. Cilium is integrated into common orchestration frameworks such as Kubernetes. A new Linux kernel technology called [eBPF](https://ebpf.io) is at the foundation of Cilium. It supports dynamic insertion of eBPF bytecode into the Linux kernel at various integration points such as network I/O, application sockets and tracepoints to implement security, networking and visibility logic.

The Cilium runtime image based on `bpftool`, `iproute2` and `llvm`, and provides the runtime environment required for Cilium setup. It also includes `gops` for debugging.

The Cilium image includes `cilium-agent` and other binaries, including `cilium`, `envoy`, `cilium-health` and `hubble` CLI. It is built on top of the Cilium runtime image.

The Cilium operator image includes the generic cilium-operator binary and the `gops` binary for debugging. It also contains CA certificates for SSL cert verification.

### Updating
1. Update GIT_TAG file with new Cilium version.
2. Install `skopeo` if you don't have it ([instructions](https://github.com/containers/skopeo/blob/main/install.md)).
3. Run `make update-digests` in this folder.
4. Update the version at the top of this Readme.
5. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
