## **Kustomize Controller**
![Version](https://img.shields.io/badge/version-v0.11.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoibldOWFUyd2ZXOXR1WkNhSVZDZkprbEowWi9nNEZrN2RMcCtRK3EvQW9qbWUzQjcxVEZvTEZ6VUw3M004WHNKQ0M1MGJ4SlU0RUJvVE1YQ0hFT0hzZ21nPSIsIml2UGFyYW1ldGVyU3BlYyI6Ing4cTAwdG9pc1I0Qk81MlQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [kustomize-controller](https://github.com/fluxcd/kustomize-controller) is a Kubernetes operator, specialized in running continuous delivery pipelines for infrastructure and workloads defined with Kubernetes manifests and assembled with Kustomize.

Some of the features of the Kustomize controller are:

* Watches for `Kustomization` objects
* Fetches artifacts produced by source-controller from `Source` objects 
* Watches `Source` objects for revision changes 
* Generates the `kustomization.yaml` file if needed
* Generates Kubernetes manifests with kustomize build
* Gecrypts Kubernetes secrets with Mozilla SOPS
* Validates the build output with client-side or APIServer dry-run
* Applies the generated manifests on the cluster
* Prunes the Kubernetes objects removed from source
* Checks the health of the deployed workloads
* Runs `Kustomizations` in a specific order, taking into account the depends-on relationship 
* Notifies whenever a `Kustomization` status changes

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/fluxcd/kustomize-controller).
