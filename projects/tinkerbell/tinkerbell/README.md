# Tinkerbell Mono-repo Build

This project builds the Tinkerbell stack services from the upstream mono-repo (github.com/tinkerbell/tinkerbell).

## Binaries

| Binary | Description | Source |
|--------|-------------|--------|
| `tinkerbell` | Unified binary containing smee, tink-server, tink-controller, rufio, tootles | `./cmd/tinkerbell` |
| `tink-agent` | Worker agent for executing workflow actions on provisioned machines | `./cmd/agent` |

## Services

The `tinkerbell` binary includes the following services used by EKS-A:

| Service | Description | Port |
|---------|-------------|------|
| smee | DHCP and iPXE service (formerly boots) | 67/68 (DHCP), 7171/7272 (HTTP/HTTPS) |
| tink-server | Workflow gRPC service | 42113 |
| tink-controller | Workflow controller | 8080 (metrics), 8081 (probe) |
| rufio | BMC controller | 8082 (metrics), 8083 (probe) |
| tootles | Metadata service (formerly hegel) | 50061 |

Note: `secondstar` (SSH over serial) is excluded via patch as EKS-A does not use it.

## Images

| Image | Component | Base Image |
|-------|-----------|------------|
| `tinkerbell/tinkerbell/tinkerbell` | Main tinkerbell binary | eks-distro-minimal-base |
| `tinkerbell/tinkerbell/tink-agent` | Worker agent | eks-distro-minimal-base-glibc |

## Patches

| Patch | Description |
|-------|-------------|
| `0001-Remove-secondstar-service.patch` | Removes secondstar service (SSH over serial) not used by EKS-A |

## Upstream

- Repository: https://github.com/tinkerbell/tinkerbell
- The mono-repo consolidates: boots→smee, tink, rufio, hegel→tootles

## GIT_TAG

The `GIT_TAG` file uses a commit SHA from the upstream tinkerbell mono-repo. Since helm
requires valid semver for chart versions, the Makefile overrides `HELM_TAG` with a `0.0.1-`
prefix and the sedfile template uses `${HELM_TAG}` for the chart version. The release CLI
in eks-anywhere uses a matching `"0.0.1-<gitTag>"` format for source image lookup.
