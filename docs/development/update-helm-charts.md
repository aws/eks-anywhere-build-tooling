# Update helm charts for EKS-A Packages

Helm charts used in EKS-A Packages are required to go through a series of modification to meet release standards. This file documents some standard procedures and modifications that are required for helm charts used in EKS-A Packages.

## Generate patch files

Helm charts modifications are done through patches. To do so, we perform the following procedures:
- clone the target helm chart repo locally;
- update helm charts locally (see details at next step); and 
- generate patch files using [`git format-patch`](https://git-scm.com/docs/git-format-patch).

## Update helm charts locally
### Update `values.yaml`

Following changes need to be made to the `values.yaml` file:

| Field                                   | Action            | Value                                     |
| :-------------------------------------: | :-------------:   | :----------------------------------------:|
| `sourceRegistry`                        | add               | `public.ecr.aws/eks-anywhere`             |
| `image:repository`                      | modify            | $(HELM_DESTINATION_REPOSITORY)[^1]        |
| `image:tag`                             | delete            | N/A       
| `image:digest`                          | add               | `{{`$(HELM_DESTINATION_REPOSITORY)[^1]`}}`|

### Update `templates` directory
Following changes need to be made to the yaml files under `templates`:

| Field                                   | Action            |  Value                                    |
| :-------------------------------------: | :---------------: | :---------------------------------------: |
| `metadata:namespace`[^2]                | add               | `{{ .Release.Namespace \| quote }}`       |
| `spec:template:spec:containers: image`  | modify            | `{{ .Values.sourceRegistry }}/{{ .Values.image.repository }}@{{ .Values.image.digest }}`        |

[^1]: Replace `$(HELM_DESTINATION_REPOSITORY)` with your project's ECR repo for the helm chart. This field should have the same value as what you define `HELM_DESTINATION_REPOSITORY` in the project Makefile.
[^2]:
    Note not all resources are in a namespace, so not all yaml files require the namespace metadata. Examples of resources not included in a namespace include `nodes`, `persistentvolumes`, `clusterrolebindings`, `clusterroles`, `csidrivers`, etc.
    You can look up if your resource is in (or not in) a namespace by running the following commands:
    ```bash
    # In a namespace
    kubectl api-resources --namespaced=true
    
    # Not in a namespace
    kubectl api-resources --namespaced=false
    ```
