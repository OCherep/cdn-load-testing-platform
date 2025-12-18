#!/bin/bash
dnf install -y docker
systemctl enable docker
systemctl start docker

docker pull YOUR_ECR/load-agent:latest
docker run -d \
  -e PROFILE_BUCKET=cdn-load-profiles \
  -e PROFILE_KEY=example.json \
  -p 9090:9090 \
  YOUR_ECR/load-agent:latest
