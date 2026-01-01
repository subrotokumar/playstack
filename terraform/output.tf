# output "cognito" {
#   value = aws_cognito_user_pool.pool.name
# }

output "storage" {
  value = {
    id =  aws_s3_bucket.s3_bucket.id
    bucket = aws_s3_bucket.s3_bucket.bucket
    bucket_domain_name = aws_s3_bucket.s3_bucket.bucket_domain_name
  }
}