locals {
  directories = [
    "",
    "directory1/",
    "directory2/",
    "directory3/",
    "directory4/level1/",
    "directory4/level1/level2/",
  ]
}

resource "aws_s3_object" "object" {
  bucket = var.bucket_name
  key    = "${local.directories[tonumber(each.value) % length(local.directories)]}example-object-${each.value}.txt"
  content = "This is an example object content for object ${each.value}."

  for_each = toset([for i in range(var.object_count) : tostring(i)])
}