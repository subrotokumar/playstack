docker_build(
  'localhost:5001/glitchr/backend',
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
    name='glitchr',
    namespace='glitchr',
    values=['k8s/helm/values.yaml', 'k8s/helm/values.secret.yaml'],
  )
)