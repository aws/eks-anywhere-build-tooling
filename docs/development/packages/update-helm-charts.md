# Update helm charts for EKS-A Packages

Helm charts used in EKS-A Packages are required to go through a series of modification to meet release standards. This file documents some standard modification procedures that apply to (almost) all packages.

Note helm chart structure varies by repo, so use judgment while applying following changes.

## Generate patch files

Helm charts modifications are done through patches. To do so, we perform the following procedures:
- clone the target helm chart repo locally;
- update helm charts locally (see details at next step); and 
- generate patch files using [`git format-patch`](https://git-scm.com/docs/git-format-patch).

## Update helm charts locally
### Update `values.yaml` file

Following changes need to be made to the `values.yaml` file:

- `sourceRegistry` or similar field
    - Add / modify this field with value `public.ecr.aws/eks-anywhere`. 
    - To test the helm chart locally, you can call `helm install` with flag `--set sourceRegistry=${YourECRRegistry}` to override its value.
- `image:repository` or similar field
    - Add / modify this field with value `${Image}`, which is the name of the image repo in ECR.
    - This field should be an enumeration of `HELM_IMAGE_LIST` in the project `Makefile`. Taking [project metallb](https://github.com/aws/eks-anywhere-build-tooling/tree/main/projects/metallb/metallb) as an example. As `HELM_IMAGE_LIST` was specified as `metallb/controller metallb/speaker` in the project `Makefile`, you should use `metallb/controller` and `metallb/speaker` as values for `controller:image:repository` and `speaker:image:repository` definitions respectively in the helm chart.
- `image:tag` or similar field
    - Delete this field (if exists) as we use `image:digest` instead.
- `image:digest` or similar field
    - Add / modify  this field with value `{{${Image}}}`.
    - As part of the helm chart build process, [helm_require.sh](https://github.com/aws/eks-anywhere-build-tooling/blob/main/build/lib/helm_require.sh) will replace the `{{${IMAGE}}}` with `${IMAGE_SHASUM}`. In the example of [project metallb](https://github.com/aws/eks-anywhere-build-tooling/tree/main/projects/metallb/metallb), `{{metallb/controller}}` will be replaced with the shasum of image `metallb/controller` before packaging the helm chart. You can verify if this update is performed successfully by reviewing the generated `sedfile` under `_output/helm`.
    - To test the helm chart outside of the `eks-anywhere-build-tooling`, you can hardcode this value.
- `imagePullPolicy` or similar field
    - Add / modify this field with value `IfNotPresent`.
- `imagePullSecrets` or similar field
    - Add field if it doesnt exist. No value is needed.
- `defaultNamespace` or similar field
    - Add/ modify this field with value of where the default namespace for the project installation

### Update `templates` directory
Following changes need to be made to the `yaml` files under `templates`:

- `metadata:namespace` or similar field
    - Add / modify this field with value `{{ .Release.Namespace | default .Values.defaultNamespace | quote }}`.
    - Note not all resources are in a namespace, so not all yaml files require the namespace metadata. Examples of resources not included in a namespace include `nodes`, `persistentvolumes`, `clusterrolebindings`, `clusterroles`, `csidrivers`, etc.
    You can look up if your resource is in (or not in) a namespace by running the following commands:
        ```bash
        # In a namespace
        kubectl api-resources --namespaced=true
        
        # Not in a namespace
        kubectl api-resources --namespaced=false
        ```
- `spec:template:spec:containers:image` or similar field
    - Add / modify this field with value `{{ .Values.sourceRegistry }}/{{ .Values.image.repository }}@{{ .Values.image.digest }}`.
- `spec:template:spec:containers:imagePullSecrets` or similar field
    - Add / modify this field with value
    `{{- with .Values.imagePullSecrets }} / imagePullSecrets: / {{- toYaml . | nindent 8 }} / {{- end }}`.

Note in some helm charts, fields above in `yaml` files are not hardcoded values but rather references to definitions in `tpl` files (also under the `templates` directory). In this case, you should update the `tpl` files directly while keeping the `yaml` files intact.

### Dealing with CRDs

For packages that include CRDs as well as custom resources, the CRDs must be deployed before the rest of the resources. CRDs can't be included in the `templates` directory because the result is a single yaml file applied once. To overcome this issue, CRDs must be defined in their own package under the `templates` directory. Once a CRDs package is ready, a dependency to that package can be declared in the package bundle definition. To add the dependency to the resulting bundle file, add `PACKAGE_DEPENDENCIES=X` to your package `Makefile` in build tooling. If both the CRDs and the actual chart are built from the same project, you'll have to use the same workaround as used in metallb which involves redefining the helm/build and helm/push targets.

