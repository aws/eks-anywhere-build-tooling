# libyang upstream fixes

These production-source patches apply to the checksum-verified libyang
6.4.13 archive in numeric order. They contain upstream fixes published after
that release; the build does not follow an unpinned development branch.

| Patch | Upstream commit | Change |
| --- | --- | --- |
| 0001 | [927fce90](https://github.com/CESNET/libyang/commit/927fce908774289a7d9c07fff5704c33965c9773) | NETCONF envelope ownership on error (CVE-2026-115327) |
| 0002 | [72cf9791](https://github.com/CESNET/libyang/commit/72cf97917a9a35dc66643207cf82dd509b5bc9eb) | Extension parent references after array reallocation |
| 0003 | [fbb673b4](https://github.com/CESNET/libyang/commit/fbb673b4e80588f8487b7ab81dfdb84ca4cf1d22) | Context printer copy direction |
| 0004 | [48103dc3](https://github.com/CESNET/libyang/commit/48103dc31d0dd0903fe51eb8ff6cb41aea86a98b) | Schema diff variable initialization |
| 0005 | [005a8b46](https://github.com/CESNET/libyang/commit/005a8b46c3d379de1c57b22cde84a362c08827bc) | JSON member name ownership |
| 0006 | [9ddf2d7e](https://github.com/CESNET/libyang/commit/9ddf2d7edbb6c3584980545c72558a11b62f5c23) | CBOR input log ownership |
| 0007 | [e2667c92](https://github.com/CESNET/libyang/commit/e2667c920996b3d1ca542c1ac0885c0effebfcdb) | CBOR opaque list ownership |
| 0008 | [362f2e44](https://github.com/CESNET/libyang/commit/362f2e44713863a4d27aed21f95c392f96ed776a) | IPv4 prefix print buffer size |
| 0009 | [e2059d18](https://github.com/CESNET/libyang/commit/e2059d18e4f9094e851ea2f97968acde088af15c) | Sorted tree metadata ownership |

Eight patches use the upstream production-source hunks unchanged. Patch 0003
adjusts one removed/context line to match the release's lack of explicit pointer
casts; its corrected replacement lines are unchanged from upstream.

The patches do not change installed public headers. Keep their commit references
and include the patches with the modified source distribution. When updating
libyang, remove fixes already present in the new version and revalidate both
native image builds and the existing library and FRR tests.
