## **cert-manager**
![Version](https://img.shields.io/badge/version-v1.18.6-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiUkphQkhWTUpOOVE1OFVLU0dHQmVFUXZJV0dJaGVLYmtEZHp0aGtDRnJBQUxtaHVqOWp3S0l6d0NlTytqNWpwc2tNTmF6RnNhMTZ3d1J1RXErR0lWcldZPSIsIml2UGFyYW1ldGVyU3BlYyI6IlQyU2lIcVVtU3ozZVZSVTgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[cert-manager](https://github.com/cert-manager/cert-manager) is a Kubernetes add-on to automate the management and issuance of TLS certificates from various issuing sources, such as [Let’s Encrypt](https://letsencrypt.org), [HashiCorp Vault](https://www.vaultproject.io), [Venafi](https://www.venafi.com/), a simple signing key pair, or self signed. It periodically ensures that certificates are valid and up-to-date, and attempts to renew certificates at an appropriate time before expiry.

cert-manager runs within Kubernetes clusters as a series of Deployment resources, its components involving a main cert-manager controller, CA injector and a webhook. 
* The controller is in charge of requesting issuance of signed certificates, leader election, approval and denial of signed certificate requests, etc.
* The CA injector helps to configure the CA certificates for various types of webhooks. It copies CA data from one of three sources: a Kubernetes Secret, a cert-manager Certificate, or from the Kubernetes API server CA certificate.
* The webhook server component is deployed as another pod that runs alongside the cert-manager controller and CA injector components. It has three main functions: `ValidatingAdmissionWebhook`, `MutatingAdmissionWebhook` and `CustomResourceConversionWebhook`.

In addition, cert-manager supports requesting certificates from ACME servers, including from Let’s Encrypt, with use of the ACME Issuer. These certificates are typically trusted on the public Internet by most computers. To successfully request a certificate, cert-manager must solve ACME Challenges which are completed in order to prove that the client owns the DNS addresses that are being requested. The component that helps to do this is the ACME Solver.

You can find the latest versions of these images on ECR Public Gallery.

[ACME Solver](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-acmesolver) | 
[cert-manager Controller](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-controller) | 
[CA injector](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-cainjector) | 
[cert-manager Webhook Server](https://gallery.ecr.aws/eks-anywhere/cert-manager/cert-manager-webhook)

### Helm Chart

The helm chart is a modified version of the source helm chart located in the jetstck/cert-manager repo at `deploy/charts/cert-manager/`.

If there are any patches to the make file, use `git format-patch $(cat ../GIT_TAG)` and add them to the `helm/patches` directory.

### Cert manager manifest

The cert-manager.yaml manifest is currently stored in the build/ repo. This is the cert-manager.yaml from the assets of the current GIT_TAG(v1.5.3)
The later tags of cert-manager (v1.7.0-alpha.0 onwards) have a make target that helps create this static cert-manager.yaml from the helm chart.
But right now we're not upgrading the cert-manager tag beyond v1.5.3 since that's what the currently used tag of Cluster API(v1.0.1) uses.
So till we use cert-manager tags lesser between v1.5.3 and v1.7.0, we need to get the cert-manager.yaml for each release from the assets section
and replace build/cert-manager.yaml with that. The reason we are doing this instead of fetching the file from github is to avoid getting the file
from github during each build, and so we're sure nothing changes in the file even if something changes later in the release assets.

### Updating

1. Update cert-manager tag when updating cluster-api tag if cluster-api is using a newer tag.
   Use the same tag that cluster-api uses by default. For instance [cluster-api v1.0.1 uses cert-manager v1.5.3 by default]
   (https://github.com/kubernetes-sigs/cluster-api/blob/v1.0.1/cmd/clusterctl/client/config/cert_manager_client.go#L30) so when updating
   to cluster-api tag v1.0.1, update cert-manager tag to v1.5.3
1. Review releases and changelogs in upstream [repo](https://github.com/cert-manager/cert-manager) and decide on new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn or @g-gaston.
1. Review the patches under patches/ folder and remove any that are either merged upstream or no longer needed.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Changes to cert-manager CRs:
   1. Usually we will update cert-manager tag only when updating CAPI tag and if the new CAPI tag uses a new cert-manager tag.
   1. If the updated cert-manager tag introduces a new API version for the cert-manager CRDs, the updated tags of upstream cluster-api providers 
      (including CAPI, CAPBK, KCP, CAPD and CAPV) will already be using the new API version for cert-manager CRs so we won't have to make any changes there.
   1. But we also use cert-manager in our custom providers like the [etcdadm-bootstrap-provider](https://github.com/aws/etcdadm-bootstrap-provider/tree/v1beta1/config/certmanager)
   and [etcdadm-controller](https://github.com/aws/etcdadm-controller/tree/v1beta1/config/certmanager) and we should use the same API version for cert-manager in these providers
   as used by the upstream providers. To make the required changes to cert-manager CRs in our providers, checkout the CAPI book's [Provider Implementers](https://cluster-api.sigs.k8s.io/developer/providers/implementers.html)
   section and review the page containing details for upgrading to the desired capi version. 
   1. For instance, when updating CAPI from v1alpha3 to v1beta1, cert-manager tag changed from v1.1.0 to v1.5.3, and upstream CAPI providers made [these](https://cluster-api.sigs.k8s.io/developer/providers/v1alpha3-to-v1alpha4.html#upgrade-cert-manager-to-v110)
   changes to their cert-manager CRs. So we made the same changes to the etcdadm providers. Similarly, check the instructions corresponding to the new
   capi version you are updating to.
1. Check the go.mod file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Update the cert-manager.yaml manifest running `make update-cert-manager-manifest` from this directory. That will download the new manifest and update `manifests/cert-manager.yaml`.
1. Update checksums and attribution using `make attribution checksums` in this folder.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
