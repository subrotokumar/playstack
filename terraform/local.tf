locals {
    version = "${var.app_version}-${var.environment}"
    bucket_name = "${var.app_name}-storage-${var.environment}"
    congito_name = "${var.app_name}"
}


locals {
  common_tags = {
    name = var.app_name
    environment = var.environment
  }
}