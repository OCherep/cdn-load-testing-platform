environment = {
  AGENT_REGION = var.agent_regions[count.index]
}

resource "aws_instance" "load" {
  count         = var.nodes
  instance_type = "c6i.large"
  ami           = data.aws_ami.al2023.id

  user_data = <<EOF
#!/bin/bash
docker run -d agent
EOF

  tags = { Role = "load-node" }
}

