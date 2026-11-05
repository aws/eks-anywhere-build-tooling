## FRRouting

This project compiles [FRRouting](https://github.com/FRRouting/frr) 10.4.3 for
the MetalLB frr-k8s runtime on the EKS-D Amazon Linux 2023 minimal Python base.
It also compiles libyang 2.1.128 and tini 0.19.0 from source.

The image provides `/sbin/tini`, `/usr/lib/frr/docker-start`, `sh`, `python3`,
`ip`, `flock`, FRR daemons, `vtysh`, and `frr-reload.py`. The `frr` user is UID
100 and GID 101; `frrvty` is GID 102. These IDs match frr-k8s's file-copy init
container. RPKI, Lua scripting, protobuf and documentation are not built.

### Building and checking

`GIT_TAG`, `LIBYANG_GIT_TAG`, and `TINI_GIT_TAG` pin the upstream source tags.
`CHECKSUMS` pins their GitHub source archives. Run `make fetch-sources` to
download and verify them before extraction. `make validate-checksums` verifies
the source archives again.

Run `make local-images` on an amd64 or arm64 host to produce an OCI archive for
that host architecture. Release builds use separate native amd64 and arm64
CodeBuild instances and the standard `combine-images` step. The project rejects
compilation for a different host architecture.

After loading the resulting image locally, verify its tools, daemon versions,
user IDs, file access and configuration reload. Integration validation with
frr-k8s must also cover BGP/BFD peers.

### Licenses and first release

FRR's `COPYING` describes a composite binary distributed under GPLv2 or later.
Its full upstream license directory, libyang's BSD license and tini's MIT
license are collected by `make gather-non-golang-licenses` and included under
`/LICENSES/github.com`, alongside `ATTRIBUTION.txt`.

**The first Amazon Open Source Team review is pending.** This project does not
record or imply approval. Obtain and record the team's actual review before the
first distribution, including FRR's GPL obligations and the runtime components.
The source URLs and fixed tags are listed in `ATTRIBUTION.txt`.

This project is excluded from independent automatic upgrades. The FRR and
libyang versions must remain aligned with the frr-k8s version shipped by
MetalLB; review all pins together when updating MetalLB.

### Updating

1. Update the source tag files and download the corresponding upstream archives.
2. Independently verify those archives and update `CHECKSUMS`; checksums are
   inputs, not generated binary checksums.
3. Regenerate the collected notices and review `ATTRIBUTION.txt`.
4. Build and validate both architectures natively, then validate frr-k8s startup,
   configuration reload, BGP and BFD.
5. Complete any required Amazon Open Source Team review before distribution.
