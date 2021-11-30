## **cert-manager**
![Version](https://img.shields.io/badge/version-v1.5.3-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiUkphQkhWTUpOOVE1OFVLU0dHQmVFUXZJV0dJaGVLYmtEZHp0aGtDRnJBQUxtaHVqOWp3S0l6d0NlTytqNWpwc2tNTmF6RnNhMTZ3d1J1RXErR0lWcldZPSIsIml2UGFyYW1ldGVyU3BlYyI6IlQyU2lIcVVtU3ozZVZSVTgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[cert-manager](https://github.com/jetstack/cert-manager) is a Kubernetes add-on to automate the management and issuance of TLS certificates from various issuing sources, such as [Let’s Encrypt](https://letsencrypt.org), [HashiCorp Vault](https://www.vaultproject.io), [Venafi](https://www.venafi.com/), a simple signing key pair, or self signed. It periodically ensures that certificates are valid and up-to-date, and attempts to renew certificates at an appropriate time before expiry.

cert-manager runs within Kubernetes clusters as a series of Deployment resources, its components involving a main cert-manager controller, CA injector and a webhook. 
* The controller is in charge of requesting issuance of signed certificates, leader election, approval and denial of signed certificate requests, etc.
* The CA injector helps to configure the CA certificates for various types of webhooks. It copies CA data from one of three sources: a Kubernetes Secret, a cert-manager Certificate, or from the Kubernetes API server CA certificate.
* The webhook server component is deployed as another pod that runs alongside the cert-manager controller and CA injector components. It has three main functions: `ValidatingAdmissionWebhook`, `MutatingAdmissionWebhook` and `CustomResourceConversionWebhook`.

In addition, cert-manager supports requesting certificates from ACME servers, including from Let’s Encrypt, with use of the ACME Issuer. These certificates are typically trusted on the public Internet by most computers. To successfully request a certificate, cert-manager must solve ACME Challenges which are completed in order to prove that the client owns the DNS addresses that are being requested. The component that helps to do this is the ACME Solver.

You can find the latest versions of these images on ECR Public Gallery.

[ACME Solver](https://gallery.ecr.aws/eks-anywhere/jetstack/cert-manager-acmesolver) | 
[cert-manager Controller](https://gallery.ecr.aws/eks-anywhere/jetstack/cert-manager-controller) | 
[CA injector](https://gallery.ecr.aws/eks-anywhere/jetstack/cert-manager-cainjector) | 
[cert-manager Webhook Server](https://gallery.ecr.aws/eks-anywhere/jetstack/cert-manager-webhook)

### Updating

1. Update cert-manager tag when updating cluster-api tag if cluster-api is using a newer tag.
   Use the same tag that cluster-api uses by default. For instance [cluster-api v1.0.1 uses cert-manager v1.5.3 by default]
   (https://github.com/kubernetes-sigs/cluster-api/blob/v1.0.1/cmd/clusterctl/client/config/cert_manager_client.go#L30) so when updating
   to cluster-api tag v1.0.1, update cert-manager tag to v1.5.3
1. Review releases and changelogs in upstream [repo](https://github.com/jetstack/cert-manager) and decide on new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jgw or @mrajashree.
1. Review the patches under patches/ folder and remove any that are either merged upstream or no longer needed.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes.
   ex: [1.1.0 compared to 1.5.3](https://github.com/jetstack/cert-manager/compare/v1.1.0...v1.5.3).
1. Check the go.mod file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Update checksums and attribution using `make update-attribution-checksums-docker PROJECT=jetstack/cert-manager` from the root of the repo.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
