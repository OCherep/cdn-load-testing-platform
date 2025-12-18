#!/usr/bin/env bash
set -e

export AWS_PROFILE=cdn-load

echo "🚀 Deploying CDN Load Platform"

cd terraform/s3
terraform init
terraform apply -auto-approve

cd ../controller
terraform init
terraform apply -auto-approve

cd ../load-nodes
terraform init
terraform apply -auto-approve

echo "✅ Deployment complete"


### old command without single-command start ###
#
#docker build -t load-agent -f docker/agent.Dockerfile .
#docker build -t controller -f docker/controller.Dockerfile .
#
#terraform -chdir=terraform/s3 apply -auto-approve
#terraform -chdir=terraform/load-nodes apply -auto-approve
#terraform -chdir=terraform/controller apply -auto-approve
#