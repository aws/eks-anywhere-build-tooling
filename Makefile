BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION?=us-west-2
IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

PROJECTS?=aws_eks-anywhere-diagnostic-collector brancz_kube-rbac-proxy kubernetes-sigs_cluster-api-provider-vsphere kubernetes-sigs_cri-tools kubernetes-sigs_vsphere-csi-driver jetstack_cert-manager kubernetes_cloud-provider-vsphere plunder-app_kube-vip kubernetes-sigs_etcdadm fluxcd_helm-controller fluxcd_kustomize-controller fluxcd_notification-controller fluxcd_source-controller rancher_local-path-provisioner mrajashree_etcdadm-bootstrap-provider mrajashree_etcdadm-controller
BUILD_TARGETS=$(addprefix build-project-, $(PROJECTS))

EKSA_TOOLS_PREREQS=kubernetes-sigs_cluster-api kubernetes-sigs_cluster-api-provider-aws kubernetes-sigs_kind fluxcd_flux2 vmware_govmomi
EKSA_TOOLS_PREREQS_BUILD_TARGETS=$(addprefix build-project-, $(EKSA_TOOLS_PREREQS))

SUPPORTED_K8S_VERSIONS=$(shell yq e 'keys | .[]' $(BASE_DIRECTORY)/projects/kubernetes-sigs/image-builder/BOTTLEROCKET_OVA_RELEASES)
OVA_TARGETS=$(addprefix release-upload-ova-ubuntu-2004-, $(SUPPORTED_K8S_VERSIONS))
OVA_TARGETS+=$(addprefix release-ova-bottlerocket-, $(SUPPORTED_K8S_VERSIONS))

.PHONY: build-all-projects
build-all-projects: $(BUILD_TARGETS) aws_bottlerocket-bootstrap aws_eks-anywhere-build-tooling

.PHONY: aws_bottlerocket-bootstrap
aws_bottlerocket-bootstrap:
	$(MAKE) release -C projects/aws/bottlerocket-bootstrap

.PHONY: aws_eks-anywhere-build-tooling
aws_eks-anywhere-build-tooling: $(EKSA_TOOLS_PREREQS_BUILD_TARGETS)
	$(MAKE) release -C projects/aws/eks-anywhere-build-tooling

.PHONY: build-project-%
build-project-%:
	$(eval PROJECT_PATH=projects/$(subst _,/,$*))
	$(MAKE) release-upload -C $(PROJECT_PATH) PROJECT_PATH=$(PROJECT_PATH)

.PHONY: release-binaries-images
release-binaries-images: build-all-projects

.PHONY: release-ovas
release-ovas:
	$(MAKE) $(OVA_TARGETS) -C projects/kubernetes-sigs/image-builder

.PHONY: clean
clean:
	make -C projects/brancz/kube-rbac-proxy clean
	make -C projects/kubernetes-sigs/cluster-api clean
	make -C projects/kubernetes-sigs/cluster-api-provider-aws clean
	make -C projects/kubernetes-sigs/cluster-api-provider-vsphere clean
	make -C projects/kubernetes-sigs/cri-tools clean
	make -C projects/kubernetes-sigs/kind clean
	make -C projects/kubernetes/cloud-provider-vsphere clean
	make -C projects/plunder-app/kube-vip clean
	make -C projects/kubernetes-sigs/etcdadm clean
	make -C projects/fluxcd/flux2 clean

	rm -rf _output

.PHONY: attribution-files
attribution-files:
	build/update-attribution-files/make_attribution.sh projects/brancz/kube-rbac-proxy
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cluster-api
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cluster-api-provider-aws
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cluster-api-provider-vsphere
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cri-tools
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/kind
	build/update-attribution-files/make_attribution.sh projects/kubernetes/cloud-provider-vsphere
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/vsphere-csi-driver
	build/update-attribution-files/make_attribution.sh projects/rancher/local-path-provisioner
	build/update-attribution-files/make_attribution.sh projects/vmware/govmomi
	build/update-attribution-files/make_attribution.sh projects/jetstack/cert-manager
	build/update-attribution-files/make_attribution.sh projects/fluxcd/flux2
	build/update-attribution-files/make_attribution.sh projects/fluxcd/helm-controller
	build/update-attribution-files/make_attribution.sh projects/fluxcd/kustomize-controller
	build/update-attribution-files/make_attribution.sh projects/fluxcd/notification-controller
	build/update-attribution-files/make_attribution.sh projects/fluxcd/source-controller
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/etcdadm
	build/update-attribution-files/make_attribution.sh projects/mrajashree/etcdadm-bootstrap-provider
	build/update-attribution-files/make_attribution.sh projects/mrajashree/etcdadm-controller
	build/update-attribution-files/make_attribution.sh projects/plunder-app/kube-vip

	cat _output/total_summary.txt

.PHONY: update-attribution-files
update-attribution-files: attribution-files
	build/update-attribution-files/create_pr.sh
