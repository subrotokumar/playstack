output "cognito_user_pool_id" {
  value = module.cognito.user_pool_id
}

output "cognito_client_id" {
  value = module.cognito.client_id
}

output "cognito_client_secret" {
  value     = module.cognito.client_secret
  sensitive = true
}


output "raw_media_bucket" {
  value = module.raw_media_bucket.bucket
}


output "raw_media_bucket_arn" {
  value = module.raw_media_bucket.bucket_arn
}


# output "processed_media_bucket_arn" {
#   value = module.processed_media_bucket.bucket_arn
# }


output "bucket_trigger_queue" {
  value = module.sqs.sqs_queue
}


output "bucket_trigger_queue_arn" {
  value = module.sqs.sqs_queue_arn
}

output "bucket_trigger_queue_url" {
  value = module.sqs.sqs_queue_url
}

output "image_registry_url" {
  value = module.ecr_image_registry.registry_url
}