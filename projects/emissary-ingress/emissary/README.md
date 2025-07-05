## **Emissary-ingress**
![Version](https://img.shields.io/badge/version-v3.10.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoid2Vla0ZUclhyN0dNN1ZLRnRoV3dBV241UDVVUUNmYUFyWVdDUWg3SUdqcThyblJvZ245cHN0cjFhUGNETDZhUkkySUNDZC9kdVFhemlpUjZGclRibUdnPSIsIml2UGFyYW1ldGVyU3BlYyI6IjRsempqaUZTNlB6TjF0bVQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

Official website: https://www.getambassador.io/docs/
Upstream repository: https://github.com/emissary-ingress/emissary

[Emissary-Ingress](https://github.com/emissary-ingress/emissary) is an open-source Kubernetes-native API Gateway + Layer 7 load balancer + Kubernetes Ingress built on [Envoy Proxy](https://www.envoyproxy.io/). Emissary-ingress is a CNCF incubation project (and was formerly known as Ambassador API Gateway.)

[Upstream Configuration examples](https://github.com/emissary-ingress/emissary/blob/master/charts/emissary-ingress/values.yaml.in)

### Updating

1. Review [releases notes](https://github.com/emissary-ingress/emissary/releases)
    * Any changes to the upstream configuration needs a thorough review + testing
    * Deprecation or removal of any protocol must be considered breaking 
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Verify the golang version has not changed.