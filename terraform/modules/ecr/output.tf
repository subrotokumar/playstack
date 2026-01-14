output "registry" {
  value = aws_ecr_repository.registry.registry_id
}

output "registry_url" {
  value = aws_ecr_repository.registry.repository_url
}