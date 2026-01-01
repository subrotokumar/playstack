resource "aws_cognito_user_pool_client" "pool_client" {
  name         = var.client_name
  user_pool_id = aws_cognito_user_pool.user_pool.id

  generate_secret = true # traditional app

  explicit_auth_flows = [
    "ALLOW_USER_AUTH",
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_USER_SRP_AUTH"
  ]

  prevent_user_existence_errors = "ENABLED"

}
