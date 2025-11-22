ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

RUN yum update -y && yum install -y iproute util-linux && yum clean all && rm -rf /var/cache/yum
COPY docker/linux/bottlerocket-bootstrap-vsphere-multi-network/bootstrap-container.sh ./
RUN chmod +x ./bootstrap-container.sh
ENTRYPOINT ["sh", "bootstrap-container.sh"]