#!/bin/bash
set -e

echo "=== CDN Load Platform Deploy ==="

# ========= CONFIG =========
AWS_REGION=eu-central-1
PROFILE_BUCKET=cdn-load-profiles
STATE_TABLE=cdn-load-tests
JWT_SECRET=supersecret
# ==========================

export AWS_REGION

echo "[1/6] Terraform state-store"
cd terraform/state-store
terraform init
terraform apply -auto-approve
cd -

echo "[2/6] Terraform S3"
cd terraform/s3
terraform init
terraform apply -auto-approve
cd -

echo "[3/6] Upload profiles"
aws s3 cp profiles/example.json s3://$PROFILE_BUCKET/example.json

echo "[4/6] Terraform controller"
cd terraform/controller
terraform init
terraform apply -auto-approve
cd -

echo "[5/6] Terraform agents"
cd terraform/load-nodes
terraform init
terraform apply -auto-approve
cd -

echo "[6/6] Build controller docker"
docker build -f docker/controller.Dockerfile -t cdn-controller .

echo "DEPLOY COMPLETE"



### old command without single-command start ###
#
#docker build -t load-agent -f docker/agent.Dockerfile .
#docker build -t controller -f docker/controller.Dockerfile .
#
#terraform -chdir=terraform/s3 apply -auto-approve
#terraform -chdir=terraform/load-nodes apply -auto-approve
#terraform -chdir=terraform/controller apply -auto-approve
#