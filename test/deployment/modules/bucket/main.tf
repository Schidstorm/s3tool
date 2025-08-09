resource "aws_s3_bucket" "test_bucket" {
  bucket = var.bucket_name
  force_destroy = true

  tags = {
    Name        = var.bucket_name
  }
}


resource "aws_iam_user" "s3_user" {
  name = var.bucket_name
}

resource "aws_iam_access_key" "s3_user_key" {
  user = aws_iam_user.s3_user.name
}

data "aws_iam_policy_document" "s3_bucket_access" {
  statement {
    actions = [
      "s3:*"
    ]
    resources = [
      aws_s3_bucket.test_bucket.arn,
      "${aws_s3_bucket.test_bucket.arn}/*"
    ]
    effect = "Allow"
  }
}

resource "aws_iam_policy" "s3_user_policy" {
  name   = "${var.bucket_name}-s3-access"
  policy = data.aws_iam_policy_document.s3_bucket_access.json
}

resource "aws_iam_user_policy_attachment" "attach_policy" {
  user       = aws_iam_user.s3_user.name
  policy_arn = aws_iam_policy.s3_user_policy.arn
}

output "access_key_id" {
  value       = aws_iam_access_key.s3_user_key.id
  description = "Access key ID for the S3 IAM user"
  sensitive   = true
}

output "secret_access_key" {
  value       = aws_iam_access_key.s3_user_key.secret
  description = "Secret access key for the S3 IAM user"
  sensitive   = true
}

output "region" {
  value       = aws_s3_bucket.test_bucket.region
  description = "AWS region where the S3 bucket is created"
}

output "bucket_name" {
  value       = aws_s3_bucket.test_bucket.bucket
  description = "ARN of the S3 bucket"
}