# FRR compatibility patch

`0001-libyang6-compatibility.patch` adapts FRR 10.4.5 to libyang 6.4.13.
It changes existing libyang call sites, not the FRR routing protocol version.

The patch selects and adapts the following upstream FRR changes:

- [libyang 4 API migration](https://github.com/FRRouting/frr/commit/908b162fcae696eac6a3b923f6d2004a2a74931f)
- [libyang 5 API migration](https://github.com/FRRouting/frr/commit/ca824ba98a274d335e53303ed3ee367f6dea2e2b)
- [LYB printing options](https://github.com/FRRouting/frr/commit/1e4957e03f845f4eea8cda232d8f30db1424ebda)
- [Canonical integer term construction](https://github.com/FRRouting/frr/commit/e5da38bfc92fb08e990b5dbf061787c08616c976)

It uses the public parent/child API, updated schema defaults and parser flags,
and canonical integer strings for the removed binary term constructor. The
existing ordinary northbound test fixture uses the same supported interface.
FRR 10.4.5's path constructors pass zero value lengths, so the newer byte-to-bit
length convention does not change those calls.

Review this patch when either the FRR or libyang pin changes. Validate the
complete native images, management/configuration operations and the frr-k8s
startup/reload/BGP/BFD flows.
