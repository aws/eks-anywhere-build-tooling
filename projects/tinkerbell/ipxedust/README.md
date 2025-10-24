## **ipxedust**
![Version](https://img.shields.io/badge/version-3c29a914f8be9b139505bfa57fffc7330e263272-blue)
![Build Status]()

[ipxedust](https://github.com/tinkerbell/ipxedust) is a TFTP and HTTP library and cli for serving iPXE binaries. ipxedust is a go mod
which [smee](https://github.com/tinkerbell/smee) depends on. The built ipxe binaries and other various image formats are embeded in
ipxedust and the resulting smee binary. This project exists to produce these built ipxe binaries. The smee project in this repo
pulls in the built ipxe binaries and overwrites the upstream vendored binaries. This allows us to control the build
of ipxe directly instead of pulling it as an implicit dependency.

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/ipxedust) and decide on release tag or commit hash to track. 
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Generally no changes will be required to eks-a build tooling, take note of the version of ipxe provided from ipxedust upstream
via [ipxe.commit](https://github.com/tinkerbell/ipxedust/blob/main/binary/script/ipxe.commit).
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
