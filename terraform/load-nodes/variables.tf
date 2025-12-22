variable "nodes" {
  type        = number
  description = "Number of load agent instances"
  default     = 2
}

variable "aws_region" {
  type        = string
  default     = "eu-central-1"
}

variable "test_id" {
  type        = string
  description = "Test ID"
}

variable "profile_bucket" {
  type        = string
  default     = "cdn-load-profiles"
}

variable "profile_key" {
  type        = string
  description = "Load profile key in S3"
}

variable "agent_regions" {
  type    = list(string)
  default = ["EU", "US"]
}
