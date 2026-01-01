variable "app_name" {
  type = string
}

variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "app_version" {
  description = "AWS region"
  type        = string
}


variable "user_pool" {
  type = string
}

variable "user_pool_client" {
  type = string
}
