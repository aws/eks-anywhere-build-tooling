## **cert-manager**
![Version](https://img.shields.io/badge/version-v1.1.0-blue)
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
