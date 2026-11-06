## FRRouting

This project compiles [FRRouting](https://github.com/FRRouting/frr) 10.4.5 for
the MetalLB frr-k8s runtime on the EKS-D Amazon Linux 2023 minimal Python base.
It also compiles libyang 6.4.13 and tini 0.19.0 from source. Checked-in libyang
patches include the upstream NETCONF correction and subsequent memory-safety
fixes. A FRR compatibility patch adapts the 10.4.5 source to the newer libyang
API.

The image provides `/sbin/tini`, `/usr/lib/frr/docker-start`, `sh`, `python3`,
`ip`, `flock`, FRR daemons, `vtysh`, and `frr-reload.py`. The `frr` user is UID
100 and GID 101; `frrvty` is GID 102. These IDs match frr-k8s's file-copy init
container. RPKI, Lua scripting, protobuf and documentation are not built.

### Building and checking

`GIT_TAG`, `LIBYANG_GIT_TAG`, and `TINI_GIT_TAG` pin the upstream source tags.
`CHECKSUMS` pins their GitHub source archives. Run `make fetch-sources` to
download and verify them before extraction, then apply the patches under
`patches/frr` and `patches/libyang` without fuzzy matching. Patch changes also
invalidate the prepared sources. `make validate-checksums` verifies the source
archives again. The patch READMEs list upstream commits and any adaptations.

Run `make local-images` on an amd64 or arm64 host to produce an OCI archive for
that host architecture. Release builds use separate native amd64 and arm64
CodeBuild instances and the standard `combine-images` step. The project rejects
compilation for a different host architecture.

After loading the resulting image locally, verify its tools, daemon versions,
user IDs, file access and configuration reload. Integration validation with
frr-k8s must also cover BGP/BFD peers.

### Licenses

FRR's `COPYING` describes a composite binary distributed under GPLv2 or later.
Its full upstream license directory, libyang's BSD license and tini's MIT
license are collected by `make gather-non-golang-licenses` and included under
`/LICENSES/github.com`, alongside `ATTRIBUTION.txt`.

The source URLs, fixed tags and local modifications are listed in
`ATTRIBUTION.txt`. The checked-in patches and build recipes are part of the
modified source distribution.

This project is excluded from independent automatic upgrades. Review and test
the FRR, libyang and frr-k8s pins together when updating MetalLB. The libyang
pin includes security updates beyond FRR's original Alpine recipe; the FRR
compatibility patch must be reviewed when either upstream version changes.

### Updating

1. Update the source tag files, review the patches against the new versions and
   download the corresponding upstream archives.
2. Independently verify those archives and update `CHECKSUMS`; checksums are
   inputs, not generated binary checksums.
3. Regenerate the collected notices and review `ATTRIBUTION.txt`.
4. Build and validate both architectures natively, then validate frr-k8s startup,
   configuration reload, BGP and BFD.
5. Complete any required Amazon Open Source Team review before distribution.
