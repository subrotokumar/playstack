module "cognito" {
  source = "../modules/cognito"

  user_pool   = var.user_pool
  client_name = var.user_pool_client
}
