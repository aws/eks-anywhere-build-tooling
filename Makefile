BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION?=us-west-2
IMAGE_REPO=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

PROJECTS?=aws_eks-anywhere brancz_kube-rbac-proxy kubernetes-sigs_cluster-api-provider-vsphere kubernetes-sigs_cri-tools kubernetes-sigs_vsphere-csi-driver jetstack_cert-manager kubernetes_cloud-provider-vsphere plunder-app_kube-vip kubernetes-sigs_etcdadm fluxcd_helm-controller fluxcd_kustomize-controller fluxcd_notification-controller fluxcd_source-controller rancher_local-path-provisioner mrajashree_etcdadm-bootstrap-provider mrajashree_etcdadm-controller tinkerbell_cluster-api-provider-tinkerbell
BUILD_TARGETS=$(addprefix build-project-, $(PROJECTS))

EKSA_TOOLS_PREREQS=kubernetes-sigs_cluster-api kubernetes-sigs_cluster-api-provider-aws kubernetes-sigs_kind fluxcd_flux2 vmware_govmomi
EKSA_TOOLS_PREREQS_BUILD_TARGETS=$(addprefix build-project-, $(EKSA_TOOLS_PREREQS))

RELEASE_BRANCH?=

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
	$(MAKE) release -C $(PROJECT_PATH) PROJECT_PATH=$(PROJECT_PATH)

.PHONY: release-binaries-images
release-binaries-images: build-all-projects

.PHONY: release-ovas
release-ovas:
	$(MAKE) release -C projects/kubernetes-sigs/image-builder

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
	make -C projects/kubernetes/cloud-provider-vsphere clean
	make -C projects/kubernetes-sigs/vsphere-csi-driver clean
	make -C projects/rancher/local-path-provisioner clean
	make -C projects/vmware/govmomi clean
	make -C projects/jetstack/cert-manager clean
	make -C projects/fluxcd/helm-controller clean
	make -C projects/fluxcd/kustomize-controller clean
	make -C projects/fluxcd/notification-controller clean
	make -C projects/fluxcd/source-controller clean
	make -C projects/mrajashree/etcdadm-bootstrap-provider clean
	make -C projects/mrajashree/etcdadm-controller clean
	make -C projects/aws/bottlerocket-bootstrap clean
	make -C projects/replicatedhq/troubleshoot clean
	make -C projects/tinkerbell/cluster-api-provider-tinkerbell clean

	make -C projects/kubernetes-sigs/image-builder clean
	make -C projects/aws/eks-anywhere-test clean
	make -C projects/aws/eks-anywhere-build-tooling clean
	make -C projects/cilium/cilium clean

	rm -rf _output

.PHONY: attribution-files
attribution-files:
	build/update-attribution-files/make_attribution.sh projects/brancz/kube-rbac-proxy
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cluster-api
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cluster-api-provider-aws
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cluster-api-provider-vsphere
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/cri-tools
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/kind
	build/update-attribution-files/make_attribution.sh projects/plunder-app/kube-vip
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/etcdadm
	build/update-attribution-files/make_attribution.sh projects/fluxcd/flux2
	build/update-attribution-files/make_attribution.sh projects/kubernetes/cloud-provider-vsphere
	build/update-attribution-files/make_attribution.sh projects/kubernetes-sigs/vsphere-csi-driver
	build/update-attribution-files/make_attribution.sh projects/rancher/local-path-provisioner
	build/update-attribution-files/make_attribution.sh projects/vmware/govmomi
	build/update-attribution-files/make_attribution.sh projects/jetstack/cert-manager
	build/update-attribution-files/make_attribution.sh projects/fluxcd/helm-controller
	build/update-attribution-files/make_attribution.sh projects/fluxcd/kustomize-controller
	build/update-attribution-files/make_attribution.sh projects/fluxcd/notification-controller
	build/update-attribution-files/make_attribution.sh projects/fluxcd/source-controller
	build/update-attribution-files/make_attribution.sh projects/mrajashree/etcdadm-bootstrap-provider
	build/update-attribution-files/make_attribution.sh projects/mrajashree/etcdadm-controller
	build/update-attribution-files/make_attribution.sh projects/aws/bottlerocket-bootstrap
	build/update-attribution-files/make_attribution.sh projects/replicatedhq/troubleshoot
	build/update-attribution-files/make_attribution.sh projects/tinkerbell/cluster-api-provider-tinkerbell

	cat _output/total_summary.txt

.PHONY: update-attribution-files
update-attribution-files: attribution-files
	build/update-attribution-files/create_pr.sh

.PHONY: run-target-in-docker
run-target-in-docker:
	build/lib/run_target_docker.sh $(PROJECT) $(MAKE_TARGET) $(IMAGE_REPO) $(RELEASE_BRANCH) $(ARTIFACTS_BUCKET)

.PHONY: update-attribution-checksums-docker
update-attribution-checksums-docker:
	build/lib/update_checksum_docker.sh $(PROJECT) $(IMAGE_REPO) $(RELEASE_BRANCH)

.PHONY: stop-docker-builder
stop-docker-builder:
	docker rm -f -v eks-a-builder

.PHONY: run-buildkit-and-registry
run-buildkit-and-registry:
	docker run -d --name buildkitd --net host --privileged moby/buildkit:v0.9.0-rootless
	docker run -d --name registry  --net host registry:2

.PHONY: stop-buildkit-and-registry
stop-buildkit-and-registry:
	docker rm -v --force buildkitd
	docker rm -v --force registry

.PHONY: generate
generate:
	build/lib/generate_projects_list.sh $(BASE_DIRECTORY)
