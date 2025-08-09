resource "random_string" "bucket_name" {
  length           = 60
  special          = false
  upper = false
}

locals {
  bucket_name = "b${random_string.bucket_name.result}"
}

module "test_bucket" {
  source = "../bucket"

  bucket_name = local.bucket_name
}

module "test_objects" {
  source = "../test_objects"

  bucket_name = module.test_bucket.bucket_name
  object_count = 100
}


resource "local_sensitive_file" "profile_parameters" {
    content  = <<EOF
access_key_id: "${module.test_bucket.access_key_id}"
secret_access_key: "${module.test_bucket.secret_access_key}"
region: "${module.test_bucket.region}"
EOF
    filename = pathexpand("~/.s3tool/profile_parameters.yaml")
}