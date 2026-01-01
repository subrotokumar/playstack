variable "user_pool" {
  type        = string
  description = "Cognito user pool name"
}

variable "client_name" {
  type        = string
  description = "Cognito app client name"
}

variable "tags" {
  type    = map(string)
  default = {}
}
