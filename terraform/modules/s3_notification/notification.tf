resource "aws_s3_bucket_notification" "s3_notification" {
  bucket = var.bucket

  queue {
    queue_arn = var.queue_arn
    events    = ["s3:ObjectCreated:*"]
    filter_prefix = "videos/"
  }
}


