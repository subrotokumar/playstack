docker_build(
  'localhost:5001/playstack/backend',
  context='.',
  dockerfile='docker/backend.dockerfile',
  # dockerfile='docker/backend.local.dockerfile',
  live_update=[
    sync('./backend', '/app/backend')
  ],
)

k8s_yaml(
  helm(
    "../playstack-ops/helm",
    name='playstack',
    namespace='playstack',
    values=['../playstack-ops/helm/values.yaml', '../playstack-ops/helm/values.secret.yaml'],
  )
)