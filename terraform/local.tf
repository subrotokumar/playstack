locals {
    version = "${var.app_version}-${var.environment}"
    bucket_name = "${var.app_name}-storage-${var.environment}"
    congito_name = "${var.app_name}"
}