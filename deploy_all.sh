#!/usr/bin/env bash
set -euo pipefail

echo "🚀 CDN Load Testing Platform – UNIVERSAL BOOTSTRAP"

# -------------------------------------------------
# 0. OS / PACKAGE MANAGER DETECTION
# -------------------------------------------------
PM=""
INSTALL_CMD=""
UPDATE_CMD=""

if command -v apt-get >/dev/null 2>&1; then
  PM="apt"
  UPDATE_CMD="sudo apt-get update -y"
  INSTALL_CMD="sudo apt-get install -y"
elif command -v dnf >/dev/null 2>&1; then
  PM="dnf"
  UPDATE_CMD="sudo dnf makecache -y"
  INSTALL_CMD="sudo dnf install -y"
elif command -v yum >/dev/null 2>&1; then
  PM="yum"
  UPDATE_CMD="sudo yum makecache -y"
  INSTALL_CMD="sudo yum install -y"
else
  echo "❌ Unsupported Linux distribution"
  exit 1
fi

echo "🖥 Package manager detected: $PM"

# -------------------------------------------------
# 1. BASE PACKAGES (SAFE)
# -------------------------------------------------
echo "📦 Installing base packages..."
$UPDATE_CMD

BASE_PKGS=(unzip jq ca-certificates gnupg tar)

for pkg in "${BASE_PKGS[@]}"; do
  if ! rpm -q "$pkg" >/dev/null 2>&1 && [[ "$PM" != "apt" ]]; then
    $INSTALL_CMD "$pkg"
  elif [[ "$PM" == "apt" ]]; then
    $INSTALL_CMD "$pkg"
  fi
done

# ---- CURL (SPECIAL HANDLING FOR AMAZON LINUX)
if ! command -v curl >/dev/null 2>&1; then
  echo "📦 Installing curl"
  if [[ "$PM" == "apt" ]]; then
    $INSTALL_CMD curl
  else
    echo "ℹ️ curl-minimal already provided by system"
  fi
else
  echo "✔ curl already present"
fi

# -------------------------------------------------
# 2. DOCKER
# -------------------------------------------------
if ! command -v docker >/dev/null 2>&1; then
  echo "🐳 Installing Docker"

  if [[ "$PM" == "apt" ]]; then
    curl -fsSL https://get.docker.com | sudo sh
  else
    sudo $INSTALL_CMD docker
    sudo systemctl enable docker
    sudo systemctl start docker
  fi
fi

docker --version

# -------------------------------------------------
# 3. DOCKER COMPOSE (v2)
# -------------------------------------------------
if ! command -v docker-compose >/dev/null 2>&1; then
  echo "🐳 Installing Docker Compose"
  sudo curl -L \
    "https://github.com/docker/compose/releases/download/v2.25.0/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose
  sudo chmod +x /usr/local/bin/docker-compose
fi

docker-compose version

# -------------------------------------------------
# 4. TERRAFORM >= 1.5 (AMAZON LINUX SAFE)
# -------------------------------------------------
if /usr/bin/terraform version >/dev/null 2>&1; then
  echo "✔ Terraform already installed: $(/usr/bin/terraform version | head -n1)"
else
  echo "🏗 Installing Terraform"

  TF_VERSION="1.6.6"
  TMP_DIR=$(mktemp -d)

  curl -fsSL \
    "https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip" \
    -o "${TMP_DIR}/terraform.zip"

  unzip -q "${TMP_DIR}/terraform.zip" -d "${TMP_DIR}"

  sudo install -m 0755 "${TMP_DIR}/terraform" /usr/bin/terraform

  rm -rf "${TMP_DIR}"

  echo "✔ Terraform installed: $(/usr/bin/terraform version | head -n1)"
fi

# Make terraform available in current shell
export PATH="/usr/bin:$PATH"


terraform version

# -------------------------------------------------
# 5. AWS ENV
# -------------------------------------------------
# export AWS_REGION=${AWS_REGION:-eu-central-1}
# export AWS_PROFILE=${AWS_PROFILE:-default}
#
# echo "☁️ AWS_REGION=$AWS_REGION"
# aws sts get-caller-identity >/dev/null || {
#   echo "❌ AWS credentials not available (IAM role missing?)"
#   exit 1
#}
# AWS
unset AWS_PROFILE
export AWS_REGION=${AWS_REGION:-eu-central-1}

echo "☁️ AWS_REGION=$AWS_REGION"
aws sts get-caller-identity || {
  echo "❌ AWS credentials not available (IAM role missing?)"
  exit 1
}

# -------------------------------------------------
# 6. INFRASTRUCTURE DEPLOY
# -------------------------------------------------
echo "🧱 Deploying infrastructure"

terraform -chdir=terraform/s3 init
terraform -chdir=terraform/s3 apply -auto-approve

terraform -chdir=terraform/controller init
terraform -chdir=terraform/controller apply -auto-approve

terraform -chdir=terraform/load-nodes init
terraform -chdir=terraform/load-nodes apply -auto-approve

# -------------------------------------------------
# 7. OBSERVABILITY
# -------------------------------------------------
echo "📊 Starting Prometheus & Grafana"
docker compose up -d prometheus grafana

# -------------------------------------------------
# 8. DONE
# -------------------------------------------------
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)

echo ""
echo "✅ DEPLOYMENT COMPLETE"
echo "📊 Grafana: http://${PUBLIC_IP}:3000"
echo "   login: admin / admin"
echo ""
echo "🧪 Controller API: http://localhost:8080"
