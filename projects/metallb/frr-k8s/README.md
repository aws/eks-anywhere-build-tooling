## **FRR-K8s**
![Version](https://img.shields.io/badge/version-v0.0.25-blue)

[FRR-K8s](https://github.com/metallb/frr-k8s) configures FRR on Kubernetes nodes.
The image provides `/frr-k8s`, `/frr-metrics`, `/frr-status`, `/statuscleaner`,
and `/frr-reloader.sh` for the FRR-K8s chart used by MetalLB.

The Go programs are built statically for `linux/amd64` and `linux/arm64`.
The image uses EKS Distro minimal base with Bash and `cp`, which the chart's
init containers need when copying programs and the reloader into shared volumes.
The reloader runs in the separate FRR container after the init step.
Go licenses and attribution are included in `/LICENSES` and `/ATTRIBUTION.txt`.

This project is excluded from independent automatic upgrades. Its version must
match the frr-k8s subchart shipped by MetalLB; update the frr-k8s and FRR pins
together when updating MetalLB.

### Updating

1. Review the upstream release notes and Dockerfile for changes to programs,
   runtime requirements, and chart init commands.
2. Update the MetalLB, frr-k8s and FRR version combination, `GOLANG_VERSION`,
   and the version above.
3. Run `make attribution-checksums` to regenerate attribution and checksums for
   both architectures.
4. Run `make build` and verify that the chart's init commands can copy the
   reloader, metrics exporter, and status exporter with executable permissions.
5. Run `make generate` from the repository root to update `UPSTREAM_PROJECTS.yaml`.

This is a new upstream project. Its initial `ATTRIBUTION.txt` requires Amazon
Open Source Team review before publication; generating the file does not grant
that approval.

### Make target changes

Run `make add-generated-help-block` from this project directory.
