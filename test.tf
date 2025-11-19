# Test file for verifying the uncommenter action

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

# These blocks should be automatically commented out by the action
moved {
  from = aws_instance.old_name
  to   = aws_instance.example
}

import {
  to = aws_s3_bucket.test_bucket
  id = "my-test-bucket"
}

removed {
  from = aws_security_group.deprecated
  lifecycle {
    destroy = false
  }
}

resource "aws_s3_bucket" "test_bucket" {
  bucket = "my-test-bucket"
}
