#!/usr/bin/env bash
set -euo pipefail

echo "🚀 CDN Load Testing Platform – FULL BOOTSTRAP"

# -------------------------------------------------
# 0. OS CHECK
# -------------------------------------------------
if [ ! -f /etc/os-release ]; then
  echo "❌ Unsupported OS"
  exit 1
fi

source /etc/os-release
echo "🖥 OS detected: $PRETTY_NAME"

# -------------------------------------------------
# 1. INSTALL PREREQUISITES
# -------------------------------------------------
echo "📦 Installing prerequisites..."

sudo apt-get update -y
sudo apt-get install -y \
  ca-certificates \
  curl \
  unzip \
  gnupg \
  lsb-release \
  jq

# ---- Docker ----
if ! command -v docker >/dev/null; then
  echo "🐳 Installing Docker"
  curl -fsSL https://get.docker.com | sudo sh
  sudo usermod -aG docker $USER
fi

# ---- Docker Compose ----
if ! command -v docker-compose >/dev/null; then
  echo "🐳 Installing Docker Compose"
  sudo curl -L \
    "https://github.com/docker/compose/releases/download/v2.25.0/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose
  sudo chmod +x /usr/local/bin/docker-compose
fi

# ---- Terraform >= 1.5 ----
if ! command -v terraform >/dev/null; then
  echo "🏗 Installing Terraform"
  curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp.gpg
  echo "deb [signed-by=/usr/share/keyrings/hashicorp.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
    | sudo tee /etc/apt/sources.list.d/hashicorp.list
  sudo apt-get update && sudo apt-get install terraform -y
fi

terraform version

# -------------------------------------------------
# 2. AWS ENV
# -------------------------------------------------
export AWS_REGION=${AWS_REGION:-eu-central-1}
export AWS_PROFILE=${AWS_PROFILE:-cdn-load}

echo "☁️ AWS_REGION=$AWS_REGION"

# -------------------------------------------------
# 3. INFRASTRUCTURE
# -------------------------------------------------
echo "🧱 Deploying infrastructure"

terraform -chdir=terraform/s3 init
terraform -chdir=terraform/s3 apply -auto-approve

terraform -chdir=terraform/controller init
terraform -chdir=terraform/controller apply -auto-approve

terraform -chdir=terraform/load-nodes init
terraform -chdir=terraform/load-nodes apply -auto-approve

# -------------------------------------------------
# 4. OBSERVABILITY
# -------------------------------------------------
echo "📊 Starting Prometheus + Grafana"
docker compose up -d prometheus grafana

# -------------------------------------------------
# 5. DONE
# -------------------------------------------------
echo "✅ Platform ready"
echo "➡️ Grafana: http://localhost:3000 (admin/admin)"
