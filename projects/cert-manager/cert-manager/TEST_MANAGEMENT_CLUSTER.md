# Testing cert-manager EKS Add-on Images on Management Cluster

## Overview

This document describes how to test the cert-manager migration to EKS add-on images
on an EKS Anywhere management cluster. The management cluster installs cert-manager
via `clusterctl init` (not Helm), using image overrides from the Bundles CRD.

## Prerequisites

- An EKS Anywhere management cluster (Docker provider is simplest for testing)
- Access to your public ECR with replicated images: `public.ecr.aws/p2x5x2t2`
- Go 1.22+ installed
- Docker running

## Image Details

| Component | Registry | Repository | Tag | Digest |
|-----------|----------|------------|-----|--------|
| controller | public.ecr.aws/p2x5x2t2 | cert-manager/eks/cert-manager-controller | v1.19.3-eksbuild.3 | sha256:13ac83e9dae34dd2ef27295ea9c76365c7a4ca247233a49fd2648266694fface |
| webhook | public.ecr.aws/p2x5x2t2 | cert-manager/eks/cert-manager-webhook | v1.19.3-eksbuild.3 | sha256:b1daed1d322b220399af5488e6ea09dd2e3ddb710abb39370e1d2d3cdc782507 |
| cainjector | public.ecr.aws/p2x5x2t2 | cert-manager/eks/cert-manager-cainjector | v1.19.3-eksbuild.3 | sha256:c54041c8224dc1e4d6dc413f89908313ac12004558401b22be1f97e1a1fd49fc |
| acmesolver | public.ecr.aws/p2x5x2t2 | cert-manager/eks/cert-manager-acmesolver | v1.19.3-eksbuild.3 | sha256:04df69da282a03c5acf8859b7a889c8099273e057b42d6ce26441aeebea0457a |

---

## Test Level 1: Unit Tests

Validates that `imageRepository()` and `buildConfig()` produce the correct clusterctl
config with the new `cert-manager/eks/` image path.

```bash
cd /Users/peirulu/peirulu-eks-anywhere-build-tooling/eks-anywhere

# Clusterctl config generation
GOSUMDB=sum.golang.org go test ./pkg/executables/ -run "TestClusterctl" -count=1 -v

# CAPI upgrader (version diff triggers upgrade correctly)
GOSUMDB=sum.golang.org go test ./pkg/clusterapi/ -count=1 -v

# Cluster manager (InstallCAPI + waitForCAPI)
GOSUMDB=sum.golang.org go test ./pkg/clustermanager/ -count=1
```

Expected: All PASS.

---

## Test Level 2: Image Pullability

Verify all images are accessible from your ECR.

```bash
for img in cert-manager-controller cert-manager-webhook cert-manager-cainjector cert-manager-acmesolver; do
  echo "=== Pulling cert-manager/eks/${img} ==="
  docker pull public.ecr.aws/p2x5x2t2/cert-manager/eks/${img}:v1.19.3-eksbuild.3
done
```

Expected: All 4 images pull successfully.

---

## Test Level 3: Patch Existing Management Cluster

Fastest way to validate images work in a running management cluster.

```bash
# Swap images on existing cert-manager deployments
kubectl set image deployment/cert-manager \
  cert-manager-controller=public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-controller@sha256:13ac83e9dae34dd2ef27295ea9c76365c7a4ca247233a49fd2648266694fface \
  -n cert-manager

kubectl set image deployment/cert-manager-webhook \
  cert-manager-webhook=public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-webhook@sha256:b1daed1d322b220399af5488e6ea09dd2e3ddb710abb39370e1d2d3cdc782507 \
  -n cert-manager

kubectl set image deployment/cert-manager-cainjector \
  cert-manager-cainjector=public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-cainjector@sha256:c54041c8224dc1e4d6dc413f89908313ac12004558401b22be1f97e1a1fd49fc \
  -n cert-manager

# Wait for rollout to complete
kubectl rollout status deployment/cert-manager -n cert-manager --timeout=120s
kubectl rollout status deployment/cert-manager-webhook -n cert-manager --timeout=120s
kubectl rollout status deployment/cert-manager-cainjector -n cert-manager --timeout=120s

# Verify images
kubectl get pods -n cert-manager -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[*].image}{"\n"}{end}'
```

Expected: All 3 deployments roll out successfully with new images.

---

## Test Level 4: Functional Validation (Certificate Issuance)

Exercises the full cert-manager stack: controller issues certs, webhook validates
requests, cainjector injects CA bundles.

```bash
# Create a self-signed ClusterIssuer and Certificate
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: test-selfsigned
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: test-cert
  namespace: cert-manager
spec:
  secretName: test-cert-tls
  issuerRef:
    name: test-selfsigned
    kind: ClusterIssuer
  commonName: test.example.com
  dnsNames:
    - test.example.com
EOF

# Wait for certificate to become Ready
kubectl wait --for=condition=Ready certificate/test-cert -n cert-manager --timeout=60s

# Verify the TLS secret was created with valid cert data
kubectl get secret test-cert-tls -n cert-manager -o jsonpath='{.data.tls\.crt}' | \
  base64 -d | openssl x509 -noout -subject -dates

# Verify webhook is working (try creating an invalid resource)
kubectl apply -f - <<EOF 2>&1 | grep -i "denied\|error\|invalid" && echo "WEBHOOK WORKING"
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: invalid-cert
  namespace: cert-manager
spec:
  secretName: invalid-cert-tls
  issuerRef:
    name: nonexistent-issuer
    kind: ClusterIssuer
  commonName: ""
EOF

# Cleanup
kubectl delete certificate test-cert -n cert-manager
kubectl delete clusterissuer test-selfsigned
kubectl delete secret test-cert-tls -n cert-manager 2>/dev/null
```

Expected:
- Certificate becomes Ready within 60s
- TLS secret contains valid X.509 certificate
- Webhook rejects or processes invalid resources correctly

---

## Test Level 5: End-to-End Management Cluster Create

Tests the full `clusterctl init` flow with your custom images by using a
`--bundles-override` that points to your ECR.

### Step 1: Create a custom Bundles manifest

```bash
# Get the current bundles from your management cluster
kubectl get bundles -A -o yaml > /tmp/current-bundles.yaml

# Edit the cert-manager section to point to your ECR images:
# Replace all cert-manager image URIs with your ECR equivalents.
# Example transformation:
#   uri: public.ecr.aws/l0g8r8j6/cert-manager/cert-manager-controller:v1.17.2-eks-a-...
#   becomes:
#   uri: public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-controller:v1.19.3-eksbuild.3
```

Create the override file at `/tmp/bundles-override.yaml`:

```yaml
# Copy the full Bundles CRD and modify only the certManager section:
# spec.versionsBundles[].certManager:
#   acmesolver:
#     uri: public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-acmesolver:v1.19.3-eksbuild.3
#   cainjector:
#     uri: public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-cainjector:v1.19.3-eksbuild.3
#   controller:
#     uri: public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-controller:v1.19.3-eksbuild.3
#   webhook:
#     uri: public.ecr.aws/p2x5x2t2/cert-manager/eks/cert-manager-webhook:v1.19.3-eksbuild.3
#   manifest:
#     uri: <path-to-your-cert-manager.yaml-manifest>
#   version: v1.19.3-eksbuild.3
```

### Step 2: Build a custom EKS-A CLI

```bash
cd /Users/peirulu/peirulu-eks-anywhere-build-tooling/eks-anywhere
make eks-anywhere-binary
```

This produces `bin/eksctl-anywhere`.

### Step 3: Create a Docker management cluster with bundles override

```bash
# Minimal Docker cluster config
cat > /tmp/mgmt-cluster.yaml <<EOF
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: Cluster
metadata:
  name: cert-manager-test
  namespace: default
spec:
  clusterNetwork:
    cniConfig:
      cilium: {}
    pods:
      cidrBlocks:
        - 192.168.0.0/16
    services:
      cidrBlocks:
        - 10.96.0.0/12
  controlPlaneConfiguration:
    count: 1
  datacenterRef:
    kind: DockerDatacenterConfig
    name: cert-manager-test
  externalEtcdConfiguration:
    count: 0
  kubernetesVersion: "1.31"
  managementCluster:
    name: cert-manager-test
  workerNodeGroupConfigurations:
    - count: 1
      name: md-0
---
apiVersion: anywhere.eks.amazonaws.com/v1alpha1
kind: DockerDatacenterConfig
metadata:
  name: cert-manager-test
  namespace: default
spec: {}
EOF

# Create cluster with bundles override
./bin/eksctl-anywhere create cluster \
  -f /tmp/mgmt-cluster.yaml \
  --bundles-override /tmp/bundles-override.yaml
```

### Step 4: Verify cert-manager on the new cluster

```bash
export KUBECONFIG=cert-manager-test/cert-manager-test-eks-a-cluster.kubeconfig

# Check cert-manager pods
kubectl get pods -n cert-manager

# Verify images are from your ECR
kubectl get pods -n cert-manager -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[*].image}{"\n"}{end}'

# Run functional validation (Level 4 above)
```

Expected:
- Cluster creates successfully
- cert-manager pods use `public.ecr.aws/p2x5x2t2/cert-manager/eks/...` images
- Certificate issuance works

### Step 5: Cleanup

```bash
./bin/eksctl-anywhere delete cluster -f /tmp/mgmt-cluster.yaml
```

---

## Test Level 6: Upgrade Path

Tests that upgrading an existing management cluster triggers cert-manager
image replacement (version change detected by `capiChangeDiff()`).

### Step 1: Create a cluster with OLD cert-manager images

Use a bundles-override with the old `cert-manager/cert-manager-controller` path and
an older version string.

### Step 2: Upgrade with NEW bundles

```bash
# Edit /tmp/bundles-override.yaml to use new images and new version string
# The version field MUST differ from the old one to trigger the upgrade

./bin/eksctl-anywhere upgrade cluster \
  -f /tmp/mgmt-cluster.yaml \
  --bundles-override /tmp/bundles-override-new.yaml
```

### Step 3: Verify

```bash
# Confirm cert-manager pods restarted with new images
kubectl get pods -n cert-manager -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[*].image}{"\n"}{end}'

# Run functional validation
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: test-after-upgrade
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: test-cert-upgrade
  namespace: default
spec:
  secretName: test-cert-upgrade-tls
  issuerRef:
    name: test-after-upgrade
    kind: ClusterIssuer
  commonName: upgrade-test.example.com
  dnsNames:
    - upgrade-test.example.com
EOF

kubectl wait --for=condition=Ready certificate/test-cert-upgrade -n default --timeout=60s
```

Expected:
- Upgrade completes successfully
- cert-manager pods have new images
- Certificate issuance works after upgrade

---

## Test Level 7: Helm Chart Install (Curated Package Path)

Tests the Helm chart for workload cluster curated package installation.

```bash
helm install cert-manager \
  /Users/peirulu/peirulu-eks-anywhere-build-tooling/eks-anywhere-build-tooling/projects/cert-manager/cert-manager/_output/helm/cert-manager \
  --namespace cert-manager-test \
  --create-namespace \
  --set sourceRegistry="public.ecr.aws/p2x5x2t2" \
  --skip-crds

# Verify
kubectl get pods -n cert-manager-test
kubectl get pods -n cert-manager-test -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.containers[*].image}{"\n"}{end}'
```

---

## How the Management Cluster Flow Works (Reference)

```
eksctl anywhere create cluster
  -> installCAPIComponentsTask
    -> ManagementComponentsFromBundles(Bundles CRD)
      -> Extracts CertManagerBundle {Controller, Webhook, Cainjector, Manifest}
    -> ClusterManager.InstallCAPI()
      -> Clusterctl.InitInfrastructure()
        -> buildOverridesLayer()
          Writes cert-manager.yaml manifest to:
            <cluster>/generated/overrides/cert-manager/<version>/cert-manager.yaml
        -> buildConfig()
          Templates clusterctl.yaml with image overrides:
            images:
              cert-manager/cert-manager-controller:
                repository: public.ecr.aws/p2x5x2t2/cert-manager/eks
                tag: v1.19.3-eksbuild.3
              cert-manager/cert-manager-cainjector:
                repository: public.ecr.aws/p2x5x2t2/cert-manager/eks
                tag: v1.19.3-eksbuild.3
              cert-manager/cert-manager-webhook:
                repository: public.ecr.aws/p2x5x2t2/cert-manager/eks
                tag: v1.19.3-eksbuild.3
            cert-manager:
              timeout: 30m
              url: "<cluster>/generated/overrides/cert-manager/<version>/cert-manager.yaml"
              version: <version>
        -> Execute: clusterctl init --config <clusterctl.yaml>
    -> waitForCAPI()
      Waits for deployments in "cert-manager" namespace:
        - cert-manager
        - cert-manager-cainjector
        - cert-manager-webhook
```

Key files in eks-anywhere repo:
- `pkg/workflows/management/create_install_capi.go` - workflow task
- `pkg/clustermanager/cluster_manager.go:380` - InstallCAPI()
- `pkg/executables/clusterctl.go:210` - InitInfrastructure()
- `pkg/executables/clusterctl.go:250` - buildConfig()
- `pkg/executables/clusterctl.go:67` - buildOverridesLayer()
- `pkg/executables/config/clusterctl.yaml` - config template
- `pkg/clustermanager/internal/deployments.go` - deployment wait list
- `release/cli/pkg/assets/config/bundle_release.go:98` - ImageRepoPrefix

---

## Code Changes Required

### eks-anywhere repo

| File | Change |
|------|--------|
| `release/cli/pkg/assets/config/bundle_release.go:98` | `ImageRepoPrefix: "cert-manager"` -> `"cert-manager/eks"` |
| `pkg/executables/clusterctl_test.go:459-468` | Update test URIs to `cert-manager/eks/cert-manager-*` |
| `pkg/executables/testdata/clusterctl_expected.yaml:50-57` | Update expected repository to `cert-manager/eks` |
| `internal/test/testdata/bundles.yaml` | Update all cert-manager image URIs |
| `release/cli/pkg/operations/testdata/main-bundle-release.yaml` | Update all cert-manager image URIs |

### eks-anywhere-build-tooling repo

Already done on `migrate-certmanager` branch:
- `projects/cert-manager/cert-manager/Makefile` - `HELM_IMAGE_LIST=eks/cert-manager-controller ...`
- `projects/cert-manager/cert-manager/EKS_ADDON_IMAGE_TAG` - `v1.19.3-eksbuild.3`
