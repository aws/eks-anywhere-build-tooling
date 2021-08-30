## **GoVMOMI**
![Version](https://img.shields.io/badge/version-v0.24.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZ1FxODROWXBIdytIZVBsNUFzODdBcngreGlZdlVwdUliRThoTGNDajBab0YzdDZ3NzVKSnBTVDBTS0lzY25sUG82MzZPMWdteE14VkZrK0F2TlppKzBjPSIsIml2UGFyYW1ldGVyU3BlYyI6IkJHNTRwbGtDV2xYRCtaZ0wiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[GoVMOMI](https://github.com/vmware/govmomi) is a Go library for interacting with VMware vSphere APIs (ESXi and/or vCenter). It primarily provides convenience functions for working with the vSphere API. It provides Go bindings to the default implementation of the VMware Managed Object Management Interface (VMOMI)

In addition to the vSphere API client, this project also includes `govc`, a CLI for vSphere. The `eks-a` tool invokes govc to perform validations on templates, fetching OVAs for building vSphere clusters, cleaning up stale VMs, etc.
