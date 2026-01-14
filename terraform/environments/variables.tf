variable "app_name" {
  description = "Name of the application."
  type        = string

  validation {
    condition     = length(var.app_name) > 0
    error_message = "app_name must not be empty."
  }
}

variable "environment" {
  description = "Deployment environment (e.g. dev, stage, prod)."
  type        = string

  validation {
    condition     = contains(["dev", "stage", "prod"], var.environment)
    error_message = "environment must be one of: dev, stage, prod."
  }
}

variable "aws_region" {
  description = "AWS region where all resources will be created."
  type        = string

  validation {
    condition     = length(var.aws_region) > 0
    error_message = "aws_region must not be empty."
  }
}

variable "app_version" {
  description = "Application version identifier used for tagging and deployments (e.g. v1.0.0)."
  type        = string

  validation {
    condition     = length(var.app_version) > 0
    error_message = "app_version must not be empty."
  }
}

variable "user_pool" {
  description = "AWS Cognito User Pool ID used for authentication."
  type        = string

  validation {
    condition     = length(var.user_pool) > 0
    error_message = "user_pool must not be empty."
  }
}

variable "user_pool_client" {
  description = "AWS Cognito User Pool Client ID associated with the user pool."
  type        = string

  validation {
    condition     = length(var.user_pool_client) > 0
    error_message = "user_pool_client must not be empty."
  }
}

variable "raw_media_bucket" {
  description = "S3 bucket name for storing raw uploaded media files."
  type        = string

  validation {
    condition     = length(var.raw_media_bucket) > 0
    error_message = "raw_media_bucket must not be empty."
  }
}

variable "processed_media_bucket" {
  description = "S3 bucket name for storing processed/transcoded media files."
  type        = string

  validation {
    condition     = length(var.processed_media_bucket) > 0
    error_message = "processed_media_bucket must not be empty."
  }
}

variable "storage_trigger_queue" {
  description = "SQS queue name used to receive S3 storage event notifications."
  type        = string

  validation {
    condition     = length(var.storage_trigger_queue) > 0
    error_message = "storage_trigger_queue must not be empty."
  }
}


variable "image_registry" {
  type        = string
  description = "Name of the container image registry (for example, an Amazon ECR repository name)."

  validation {
    condition     = length(var.image_registry) > 0
    error_message = "Registry name must not be empty."
  }
}

variable "image_mutable" {
  type        = string
  description = "Image tag mutability setting for the registry. Allowed values are MUTABLE or IMMUTABLE."
  default     = "MUTABLE"

  validation {
    condition     = contains(["MUTABLE", "IMMUTABLE"], var.image_mutable)
    error_message = "image_mutable must be either 'MUTABLE' or 'IMMUTABLE'."
  }
}

variable "image_scan_on_push" {
  type        = bool
  description = "Whether to enable automatic image vulnerability scanning when an image is pushed to the registry."
  default     = true
}
