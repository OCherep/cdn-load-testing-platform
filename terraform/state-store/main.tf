resource "aws_dynamodb_table" "tests" {
  name         = "cdn-load-tests"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "test_id"

  attribute {
    name = "test_id"
    type = "S"
  }

  ttl {
    attribute_name = "expires_at"
    enabled        = true
  }

  tags = {
    Project = "cdn-load-platform"
  }
}

