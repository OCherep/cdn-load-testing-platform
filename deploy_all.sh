#!/usr/bin/env bash
set -euo pipefail

### ================================
### CONFIG
### ================================
AWS_PROFILE=${AWS_PROFILE:-cdn-load}
AWS_REGION=${AWS_REGION:-eu-central-1}

CONTROLLER_IMAGE=cdn-load-controller
AGENT_IMAGE=cdn-load-agent

ROOT_DIR=$(pwd)

echo "🚀 CDN Load Platform – ONE CLICK DEPLOY"
echo "AWS_PROFILE=$AWS_PROFILE"
echo "AWS_REGION=$AWS_REGION"
echo

### ================================
### PREREQUISITES
### ================================
echo "🔍 Checking dependencies..."

for cmd in aws docker terraform go; do
  if ! command -v $cmd >/dev/null 2>&1; then
    echo "❌ Missing dependency: $cmd"
    exit 1
  fi
done

aws sts get-caller-identity --profile $AWS_PROFILE >/dev/null \
  || { echo "❌ AWS auth failed"; exit 1; }

echo "✅ Dependencies OK"
echo

### ================================
### DOCKER BUILD
### ================================
echo "🐳 Building Docker images..."

docker build -t $AGENT_IMAGE -f docker/agent.Dockerfile .
docker build -t $CONTROLLER_IMAGE -f docker/controller.Dockerfile .

echo "✅ Docker images built"
echo

### ================================
### TERRAFORM: S3 + DYNAMO
### ================================
echo "🗄️  Deploying S3 + DynamoDB..."

cd terraform/s3
terraform init -input=false
terraform apply -auto-approve \
  -var "aws_region=$AWS_REGION"

cd "$ROOT_DIR"

### ================================
### TERRAFORM: CONTROLLER
### ================================
echo "🎛️  Deploying Controller..."

cd terraform/controller
terraform init -input=false
terraform apply -auto-approve \
  -var "aws_region=$AWS_REGION" \
  -var "controller_image=$CONTROLLER_IMAGE"

cd "$ROOT_DIR"

### ================================
### TERRAFORM: LOAD NODES (IDLE)
### ================================
echo "⚙️  Preparing Load Nodes module..."

cd terraform/load-nodes
terraform init -input=false

cd "$ROOT_DIR"

### ================================
### GRAFANA
### ================================
echo "📊 Starting Grafana..."

docker rm -f grafana-cdn-load >/dev/null 2>&1 || true

docker run -d \
  --name grafana-cdn-load \
  -p 3000:3000 \
  grafana/grafana:latest

echo
echo "======================================"
echo "✅ DEPLOYMENT COMPLETE"
echo "======================================"
echo
echo "Controller API : http://<controller_public_ip>:8080"
echo "Grafana UI     : http://localhost:3000"
echo "Grafana login  : admin / admin"
echo
echo "👉 NEXT STEPS:"
echo "1. POST /auth/login"
echo "2. POST /tests"
echo "3. POST /tests/:id/start"
echo
