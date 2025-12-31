resource "aws_cognito_user_pool" "pool" {
  name = local.congito_name
}