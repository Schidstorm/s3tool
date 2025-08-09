variable "bucket_name" {
  description = "The ARN of the S3 bucket"
  type        = string
  
}

variable "object_count" {
  description = "The number of objects to create in the S3 bucket"
  type        = number
  default     = 0
}