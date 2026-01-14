resource "aws_ecr_repository" "registry" {
  name                 = var.image_registry
  image_tag_mutability = var.image_mutable

  image_scanning_configuration {
    scan_on_push = var.image_scan_on_push
  }
}