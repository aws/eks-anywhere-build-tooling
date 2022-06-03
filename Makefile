BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
BUILD_LIB=${BASE_DIRECTORY}/build/lib
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION?=us-west-2
IMAGE_REPO?=$(if $(AWS_ACCOUNT_ID),$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com,localhost:5000)
ECR_PUBLIC_URI?=$(shell aws ecr-public describe-registries --region us-east-1 --query 'registries[0].registryUri' --output text)

PROJECTS?=aws_eks-anywhere aws_eks-anywhere-packages aws-containers_hello-eks-anywhere brancz_kube-rbac-proxy kubernetes-sigs_cluster-api-provider-vsphere kubernetes-sigs_cri-tools kubernetes-sigs_vsphere-csi-driver jetstack_cert-manager kubernetes_cloud-provider-vsphere plunder-app_kube-vip kubernetes-sigs_etcdadm fluxcd_helm-controller fluxcd_kustomize-controller fluxcd_notification-controller fluxcd_source-controller rancher_local-path-provisioner aws_etcdadm-bootstrap-provider aws_etcdadm-controller tinkerbell_cluster-api-provider-tinkerbell tinkerbell_hegel cloudflare_cfssl tinkerbell_boots tinkerbell_hub tinkerbell_pbnj tinkerbell_hook tinkerbell_rufio aws_cluster-api-provider-aws-snow distribution_distribution goharbor_harbor cilium_cilium grpc-ecosystem_grpc-health-probe kubernetes-sigs_cluster-api-provider-cloudstack aquasecurity_harbor-scanner-trivy aquasecurity_trivy
BUILD_TARGETS=$(addprefix build-project-, $(PROJECTS))

EKSA_TOOLS_PREREQS=kubernetes-sigs_cluster-api kubernetes-sigs_cluster-api-provider-aws kubernetes-sigs_kind fluxcd_flux2 vmware_govmomi replicatedhq_troubleshoot helm_helm tinkerbell_tink apache_cloudstack-cloudmonkey
EKSA_TOOLS_PREREQS_BUILD_TARGETS=$(addprefix build-project-, $(EKSA_TOOLS_PREREQS))

ALL_PROJECTS=$(PROJECTS) $(EKSA_TOOLS_PREREQS) aws_bottlerocket-bootstrap aws_eks-anywhere-build-tooling kubernetes-sigs_image-builder

RELEASE_BRANCH?=
GIT_HASH:=$(shell git -C $(BASE_DIRECTORY) rev-parse HEAD)
ALL_PROJECTS=$(shell $(BUILD_LIB)/all_projects.sh $(BASE_DIRECTORY))

.PHONY: clean-project-%
clean-project-%:
	$(eval PROJECT_PATH=projects/$(subst _,/,$*))
	$(MAKE) clean -C $(PROJECT_PATH)

.PHONY: clean
clean: $(addprefix clean-project-, $(ALL_PROJECTS))
	rm -rf _output

.PHONY: add-generated-help-block-project-%
add-generated-help-block-project-%:
	$(eval PROJECT_PATH=projects/$(subst _,/,$*))
	$(MAKE) add-generated-help-block -C $(PROJECT_PATH) RELEASE_BRANCH=1-21

.PHONY: add-generated-help-block
add-generated-help-block: $(addprefix add-generated-help-block-project-, $(ALL_PROJECTS))
	build/update-attribution-files/create_pr.sh

.PHONY: attribution-files-project-%
attribution-files-project-%:
	$(eval PROJECT_PATH=projects/$(subst _,/,$*))
	build/update-attribution-files/make_attribution.sh $(PROJECT_PATH) attribution

.PHONY: attribution-files
attribution-files: $(addprefix attribution-files-project-, $(ALL_PROJECTS))
	cat _output/total_summary.txt

.PHONY: checksum-files-project-%
checksum-files-project-%:
	$(eval PROJECT_PATH=projects/$(subst _,/,$*))
	build/update-attribution-files/make_attribution.sh $(PROJECT_PATH) checksums

.PHONY: checksum-files
checksum-files: $(addprefix checksum-files-project-, $(ALL_PROJECTS))
	build/update-attribution-files/create_pr.sh

.PHONY: update-attribution-files
update-attribution-files: add-generated-help-block attribution-files checksum-files
	build/lib/update_go_versions.sh
	build/update-attribution-files/create_pr.sh

.PHONY: run-target-in-docker
run-target-in-docker:
	build/lib/run_target_docker.sh $(PROJECT) $(MAKE_TARGET) $(IMAGE_REPO) "$(RELEASE_BRANCH)" $(ARTIFACTS_BUCKET)

.PHONY: update-attribution-checksums-docker
update-attribution-checksums-docker:
	build/lib/update_checksum_docker.sh $(PROJECT) $(IMAGE_REPO) $(RELEASE_BRANCH)

.PHONY: stop-docker-builder
stop-docker-builder:
	docker rm -f -v eks-a-builder

.PHONY: run-buildkit-and-registry
run-buildkit-and-registry:
	docker run -d --name buildkitd --net host --privileged moby/buildkit:v0.10.3-rootless
	docker run -d --name registry  --net host registry:2

.PHONY: stop-buildkit-and-registry
stop-buildkit-and-registry:
	docker rm -v --force buildkitd
	docker rm -v --force registry

.PHONY: generate-project-list
generate-project-list:
	build/lib/generate_projects_list.sh $(BASE_DIRECTORY)

.PHONY: generate-staging-buildspec
generate-staging-buildspec:
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "$(ALL_PROJECTS)"
	
.PHONY: generate
generate: generate-project-list generate-staging-buildspec

.PHONY: validate-generated
validate-generated: generate
	@if [ "$$(git status --porcelain -- UPSTREAM_PROJECTS.yaml release/staging-build.yml | wc -l)" -gt 0 ]; then \
		echo "Error: Generated files, UPSTREAM_PROJECTS.yaml release/staging-build.yml, do not match expected. Please run `make generate` to update"; \
		git diff -- UPSTREAM_PROJECTS.yaml release/staging-build.yml; \
		exit 1; \
	fi

.PHONY: check-project-path-exists
check-project-path-exists:
	@if ! stat $(PROJECT_PATH) &> /dev/null; then \
		echo "false"; \
	else \
		echo "true"; \
	fi
