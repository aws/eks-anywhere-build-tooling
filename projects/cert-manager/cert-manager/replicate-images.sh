#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

TAG="v1.20.2-eksbuild.3"
SRC="602401143452.dkr.ecr.us-west-2.amazonaws.com"
DST="public.ecr.aws/p2x5x2t2"

IMAGES=(
  "eks/cert-manager-controller"
  "eks/cert-manager-webhook"
  "eks/cert-manager-cainjector"
  "eks/cert-manager-acmesolver"
)

echo "Fetching credentials..."
SRC_CREDS="AWS:$(aws ecr get-login-password --region us-west-2)"
DST_CREDS="AWS:$(aws ecr-public get-login-password --region us-east-1)"

for image in "${IMAGES[@]}"; do
  repo="cert-manager/${image}"
  echo "Creating repo ${repo} (if not exists)..."
  aws ecr-public create-repository --repository-name "${repo}" --region us-east-1 2>/dev/null || true

  echo "Copying ${SRC}/${image}:${TAG} -> ${DST}/${repo}:${TAG}"
  skopeo copy --all \
    --src-creds "${SRC_CREDS}" \
    --dest-creds "${DST_CREDS}" \
    "docker://${SRC}/${image}:${TAG}" \
    "docker://${DST}/${repo}:${TAG}"
done

echo "Done! All images replicated to ${DST}/cert-manager/eks/*:${TAG}"
