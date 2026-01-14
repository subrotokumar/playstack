module "cognito" {
  source = "../modules/cognito"

  user_pool   = var.user_pool
  client_name = var.user_pool_client
}

# module "processed_media_bucket" {
#   source = "../modules/s3"

#   bucket = var.processed_media_bucket
# }


module "raw_media_bucket" {
  source = "../modules/s3"

  bucket = var.raw_media_bucket
}

module "sqs" {
  source = "../modules/sqs"

  sqs_queue  = var.storage_trigger_queue
  bucket_arn = module.raw_media_bucket.bucket_arn
}

module "s3_notification" {
  source = "../modules/s3_notification"

  queue_arn = module.sqs.sqs_queue_arn
  bucket    = module.raw_media_bucket.bucket
}


module "ecr_image_registry" {
  source = "../modules/ecr"

  image_registry = var.image_registry
} 