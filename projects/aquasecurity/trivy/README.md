## **trivy**
![Version](https://img.shields.io/badge/version-v0.68.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiMVBvZE5FTEtYaVpuWUJ3eGd2Tis1dHAxT0ZKcXBuWkNVUmpjL0pRVnduRUl2Qm1XZ29xbHBENU5wVGM3TzVTTXhFTS83VUtrWGdCVU9lVkVxSmFhUnBFPSIsIml2UGFyYW1ldGVyU3BlYyI6IkQzTU9tSEd0YWZDc0NVYkIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Trivy](https://github.com/aquasecurity/trivy/) is a simple and comprehensive scanner for vulnerabilities in container images, file systems, and Git repositories, as well as for configuration issues. Trivy detects vulnerabilities of OS packages (Alpine, RHEL, CentOS, etc.) and language-specific packages (Bundler, Composer, npm, yarn, etc.). In addition, Trivy scans Infrastructure as Code (IaC) files such as Terraform, Dockerfile and Kubernetes, to detect potential configuration issues that expose your deployments to the risk of attack. Trivy also scans hardcoded secrets like passwords, API keys and tokens.

### Updating

1. Update trivy tag when updating harbor tag if harbor is using a newer tag. Use the same tag that harbor uses by default. For instance [harbor v2.5.1 uses trivy v0.26.0 by default](https://github.com/goharbor/harbor/blob/v2.5.1/Makefile#L114) so when updating to harbor tag v2.5.1, update trivy tag to v0.26.0 or higher if security patching requires.
1. Review releases and changelogs in upstream [repo](https://github.com/aquasecurity/trivy) and decide on new version.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. Check the `build` target for any build flag changes, tag changes, dependencies, etc. Check that the manifest target has not changed, this is called from our Makefile.
1. Check the `go.mod` file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in Makefile to match the version upstream.
1. Update checksums and attribution using make `run-attribution-checksums-in-docker`.
1. Update the version at the top of this `README`.
1. Run `make generate` to update the `UPSTREAM_PROJECTS.yaml` file.
