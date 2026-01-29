## **Project Version Tracker Tool**

The `version-tracker` CLI is a command-line tool that is used to track and upgrade the versions of upstream projects in this repository. The projects targetted by this tool are simple upstream projects that don't have patches and can be automatically upgraded to the latest revision without the need for conflict resolution. In the scope of this tool, a _revision_ is defined as _a Git tag or commit object being tracked for an upstream project_.

The CLI has three subcommands namely, `display`, `list-projects` and `upgrade`. Their functionality and usage are described in the sections below.

### The `display` subcommand

The `display` subcommand is used to tabulate the current and latest revision for a particular project or all projects in the repository. It takes in an optional `print-latest-revision` flag to display only the latest revision for a particular project instead of a tabular output.

#### Usage

```
$ version-tracker display --help                  
Use this command to display the version information for a particular project or for all projects

Usage:
  version-tracker display --project <project name> [flags]

Flags:
  -h, --help                   help for display
      --project string         Specify the project name to track versions for

Global Flags:
  -v, --verbosity int   Set the logging verbosity level
```

#### Sample output

* Displaying all project versions
```
$ version-tracker display
ORGANIZATION          REPOSITORY                       CURRENT VERSION                           LATEST VERSION                                 
apache                cloudstack-cloudmonkey           6.3.0                                     6.3.0                                          
aquasecurity          trivy                            v0.37.3                                   v0.47.0                                        
aws                   etcdadm-bootstrap-provider       v1.0.10                                   v1.0.10                                        
aws                   etcdadm-controller               v1.0.15                                   v1.0.15                                        
aws                   rolesanywhere-credential-helper  v1.0.4                                    v1.1.1                                         
aws-observability     aws-otel-collector               v0.25.0                                   v0.35.0                                        
brancz                kube-rbac-proxy                  v0.14.2                                   v0.15.0                                        
cert-manager          cert-manager                     v1.13.1                                   v1.13.2                                        
cilium                cilium                           v1.12.15-eksa.1                           v1.14.4                                        
containerd            containerd                       v1.6.21                                   v1.6.25                                        
distribution          distribution                     v2.8.1                                    v2.8.3                                         
emissary-ingress      emissary                         v3.8.1                                    v3.9.1                                         
envoyproxy            envoy                            v1.22.2.0-prod                            v1.28.0                                        
fluxcd                flux2                            v2.0.0                                    v2.1.2                                         
fluxcd                helm-controller                  v0.35.0                                   v0.36.2                                        
fluxcd                kustomize-controller             v1.0.0                                    v1.1.1                                         
fluxcd                notification-controller          v1.0.0                                    v1.1.0                                         
fluxcd                source-controller                v1.0.0                                    v1.1.2                                         
goharbor              harbor                           v2.9.1                                    v2.9.1                                         
goharbor              harbor-scanner-trivy             v0.30.7                                   v0.30.19                                       
helm                  helm                             v3.12.1                                   v3.13.2                                        
kube-vip              kube-vip                         v0.6.0                                    v0.6.3                                         
kubernetes            autoscaler                       5bcb526e08c17ff93cc6093ee89a95730a90e45b  cluster-autoscaler-chart-9.32.1                
kubernetes            cloud-provider-aws               v1.28.1                                   helm-chart-aws-cloud-controller-manager-0.0.8  
kubernetes            cloud-provider-vsphere           v1.28.0                                   v1.23.5                                        
kubernetes-sigs       cluster-api                      v1.5.3                                    v1.5.3                                         
kubernetes-sigs       cluster-api-provider-cloudstack  v0.4.9-rc8                                v0.4.8                                         
kubernetes-sigs       cluster-api-provider-vsphere     v1.7.0                                    v1.8.4                                         
kubernetes-sigs       cri-tools                        v1.28.0                                   v1.28.0                                        
kubernetes-sigs       etcdadm                          f089d308442c18f487a52d09fd067ae9ac7cd8f2  v0.1.5                                         
kubernetes-sigs       image-builder                    v0.1.19                                   v0.1.21                                        
kubernetes-sigs       kind                             v0.20.0                                   v0.20.0                                        
metallb               metallb                          v0.13.7                                   v0.13.12                                       
nutanix-cloud-native  cluster-api-provider-nutanix     v1.2.3                                    v1.2.4                                         
opencontainers        runc                             v1.1.7                                    v1.1.10                                        
prometheus            node_exporter                    v1.5.0                                    v1.7.0                                         
prometheus            prometheus                       v2.43.0                                   v2.48.0                                        
rancher               local-path-provisioner           v0.0.24                                   v0.0.25                                        
redis                 redis                            6.2.6                                     7.2.3                                          
replicatedhq          troubleshoot                     v0.69.2                                   v0.78.1                                        
tinkerbell            boots                            v0.8.1                                    v0.10.1                                        
tinkerbell            cluster-api-provider-tinkerbell  v0.4.0                                    v0.4.0                                         
tinkerbell            hegel                            v0.10.1                                   v0.11.1                                        
tinkerbell            hook                             9d54933a03f2f4c06322969b06caa18702d17f66  v0.8.1                                         
tinkerbell            hub                              404dab73a8a7f33e973c6e71782f07e82b125da9  404dab73a8a7f33e973c6e71782f07e82b125da9       
tinkerbell            rufio                            afd7cd82fa08dae8f9f3ffac96eb030176f3abbd  v0.3.2                                         
tinkerbell            tink                             v0.8.0                                    v0.9.0                                         
torvalds              linux                            v5.17                                     v6.7-rc3                                       
vmware                govmomi                          v0.30.5                                   v0.33.0
```

* Displaying the current and latest versions for a particular project
```
$ version-tracker display --project cilium/cilium
ORGANIZATION  REPOSITORY  CURRENT VERSION  LATEST VERSION  UP-TO-DATE  
cilium        cilium      v1.15.16-eksa.1  v1.17.6-0       false
```

### The `list-projects` subcommand

The `list-projects` subcommand is used to tabulate the various GitHub projects being built from this repository. The projects are grouped by organization or owner and presented in a tabular format.

#### Usage

```
$ version-tracker list-projects --help
Use this command to list the upstream projects that are built from the eks-anywhere-build-tooling repository

Usage:
  version-tracker list-projects [flags]

Flags:
  -h, --help   help for list-projects

Global Flags:
  -v, --verbosity int   Set the logging verbosity level
```

#### Sample output

```
$ version-tracker list-projects
ORGANIZATION          REPOSITORY                       
--------------------  -------------------------------  
apache                cloudstack-cloudmonkey           
--------------------  -------------------------------  
aquasecurity          trivy                            
--------------------  -------------------------------  
                      etcdadm-bootstrap-provider       
aws                   etcdadm-controller               
                      rolesanywhere-credential-helper  
--------------------  -------------------------------  
aws-observability     aws-otel-collector               
--------------------  -------------------------------  
brancz                kube-rbac-proxy                  
--------------------  -------------------------------  
cert-manager          cert-manager                     
--------------------  -------------------------------  
cilium                cilium                           
--------------------  -------------------------------  
containerd            containerd                       
--------------------  -------------------------------  
distribution          distribution                     
--------------------  -------------------------------  
emissary-ingress      emissary                         
--------------------  -------------------------------  
envoyproxy            envoy                            
--------------------  -------------------------------  
                      flux2                            
                      helm-controller                  
fluxcd                kustomize-controller             
                      notification-controller          
                      source-controller                
--------------------  -------------------------------  
goharbor              harbor                           
                      harbor-scanner-trivy             
--------------------  -------------------------------  
helm                  helm                             
--------------------  -------------------------------  
isc-projects          dhcp                             
--------------------  -------------------------------  
kube-vip              kube-vip                         
--------------------  -------------------------------  
                      autoscaler                       
kubernetes            cloud-provider-aws               
                      cloud-provider-vsphere           
--------------------  -------------------------------  
                      cluster-api                      
                      cluster-api-provider-cloudstack  
                      cluster-api-provider-vsphere     
kubernetes-sigs       cri-tools                        
                      etcdadm                          
                      image-builder                    
                      kind                             
--------------------  -------------------------------  
linuxkit              linuxkit                             
--------------------  -------------------------------  
metallb               metallb                          
--------------------  -------------------------------  
nutanix-cloud-native  cloud-provider-nutanix
                      cluster-api-provider-nutanix     
--------------------  -------------------------------  
opencontainers        runc                             
--------------------  -------------------------------  
prometheus            node_exporter                    
                      prometheus                       
--------------------  -------------------------------  
rancher               local-path-provisioner           
--------------------  -------------------------------  
redis                 redis                            
--------------------  -------------------------------  
replicatedhq          troubleshoot                     
--------------------  -------------------------------  
                      actions
                      boots                          
                      charts  
                      cluster-api-provider-tinkerbell  
                      hegel                            
tinkerbell            hook                             
                      ipxedust                              
                      rufio                            
                      tink                             
--------------------  -------------------------------  
torvalds              linux                            
--------------------  -------------------------------  
vmware                govmomi                          
--------------------  -------------------------------
```

### The `upgrade` subcommand

The `upgrade` subcommand is used to upgrade the Git revision of a particular project. This command takes in a project name as input and updates the various version files pertaining to the project, such as Git tag, Go version, checksums, etc. Then it creates a PR with these changes from a fork of the build-tooling repository. The PR can then be reviewed and merged by a repository maintainer.

When patches fail to apply during an upgrade, the command can optionally publish an event to AWS EventBridge (if `ENABLE_AUTO_PATCH_FIX=true`) for downstream automation.

#### Usage

```
$ version-tracker upgrade --help
Use this command to upgrade the Git tag and related versions for a particular project in the EKS-A build-tooling repository

Usage:
  version-tracker upgrade --project <project name> [flags]

Flags:
      --dry-run          Upgrade the project locally but do not push changes and create PR
  -h, --help             help for upgrade
      --project string   Specify the project name to upgrade versions for

Global Flags:
  -v, --verbosity int   Set the logging verbosity level
```

#### Sample output

```
$ version-tracker upgrade --project vmware/govmomi
Project is out of date. {"Current version": "v0.30.5", "Latest version": "v0.33.0"}
Updating Git tag file corresponding to the project.
Project Go version needs to be updated. {"Current Go version": "1.18", "Latest Go version": "1.20"}
Updating Go version file corresponding to the project
Updating Git tag and Go version in upstream projects tracker file
Updating project checksums and attribution files
Updating project readme
Creating pull request with updated files
```

### The `fix-patches` subcommand

The `fix-patches` subcommand is used to automatically fix patches that fail to apply during version upgrades. This command uses AI (AWS Bedrock with Claude) to analyze patch failures, understand the context, and generate corrected patches.

#### How it works

1. **Analyzes patch failures** - Extracts information about which patches failed and why
2. **Gathers context** - Fetches relevant code from GitHub (PR diff, failed files, patch content)
3. **Generates fixes** - Uses Claude AI to understand the changes and create corrected patches
4. **Validates fixes** - Applies patches and runs builds to ensure they work
5. **Updates metadata** - Regenerates checksums and attribution files
6. **Pushes changes** - Commits and pushes fixed patches back to the PR

#### Usage

```
$ version-tracker fix-patches --help
Automatically fix patches that fail to apply during version upgrades

Usage:
  version-tracker fix-patches --project <project name> --pr <pr number> [flags]

Flags:
  -h, --help                help for fix-patches
      --max-attempts int    Maximum number of fix attempts per patch (default 3)
      --pr int              PR number where patches failed
      --project string      Project name (e.g., kubernetes-sigs/image-builder)

Global Flags:
  -v, --verbosity int   Set the logging verbosity level
```

#### Required Environment Variables

```bash
# GitHub access
export GITHUB_TOKEN=<your-github-token>

# AWS Bedrock for AI-powered patch fixing
export AWS_BEDROCK_REGION=us-west-2
export BEDROCK_MODEL_ID=anthropic.claude-3-5-sonnet-20241022-v2:0
```

#### Sample output

```
$ version-tracker fix-patches --project kubernetes-sigs/image-builder --pr 5005
Starting patch fixing workflow  {"project": "kubernetes-sigs/image-builder", "pr": 5005}
Analyzing patch failures...
Found 1 failed patch: 0001-EKS-A-AMI-changes.patch
Extracting context from GitHub PR #5005...
Fetching PR diff (2847 lines)...
Fetching failed files content...
Building prompt for LLM (estimated 45000 tokens)...
Calling Claude to generate patch fix...
Received fix from LLM (3421 tokens)
Applying fixed patch...
Patch applied successfully!
Running build validation...
Build succeeded!
Updating checksums...
Updating attribution files...
Committing changes...
Pushing to PR branch...
Successfully fixed patches for kubernetes-sigs/image-builder PR #5005
```

#### Special Cases

Some projects require special handling:

- **kubernetes/autoscaler**: Patches are in `projects/kubernetes/autoscaler/patches/` and need to be regenerated using `git format-patch` after fixes are applied to the source code.

#### Running Locally

Run `fix-patches` manually when you see a PR with patch failures:

```bash
# Set up environment
export GITHUB_TOKEN=$(gh auth token)
export AWS_BEDROCK_REGION=us-west-2
export BEDROCK_MODEL_ID=anthropic.claude-3-5-sonnet-20241022-v2:0

# Run from the repository root
cd /path/to/eks-anywhere-build-tooling
version-tracker fix-patches --project fluxcd/source-controller --pr 4883
```