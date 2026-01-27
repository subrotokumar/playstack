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
    'k8s/helm',
    name='playstack',
    namespace='playstack',
    values=['k8s/helm/values.yaml', 'k8s/helm/values.secret.yaml'],
  )
)