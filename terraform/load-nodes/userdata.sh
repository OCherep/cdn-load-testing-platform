#!/bin/bash
set -e

dnf install -y docker
systemctl enable docker
systemctl start docker

docker pull YOUR_ECR/load-agent:latest

docker run -d \
  --restart=always \
  -e AWS_REGION=${AWS_REGION} \
  -e TEST_ID=${TEST_ID} \
  -e PROFILE_BUCKET=${PROFILE_BUCKET} \
  -e PROFILE_KEY=${PROFILE_KEY} \
  -e AGENT_REGION=${AGENT_REGION} \
  -p 9090:9090 \
  YOUR_ECR/load-agent:latest
