terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "aws_instance" "load" {
  count         = var.nodes
  ami           = data.aws_ami.al2023.id
  instance_type = "c6i.large"

  user_data = templatefile("${path.module}/userdata.sh", {
    AWS_REGION     = var.aws_region
    TEST_ID        = var.test_id
    PROFILE_BUCKET = var.profile_bucket
    PROFILE_KEY    = var.profile_key
    AGENT_REGION   = var.agent_regions[count.index % length(var.agent_regions)]
  })

  tags = {
    Name = "cdn-load-agent-${count.index}"
    Role = "load-node"
    Test = var.test_id
  }
}
