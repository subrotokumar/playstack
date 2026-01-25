docker_build(
  'kind-registry:5000/glitchr/backend',
  context='../..',
  dockerfile='docker/backend.local.dockerfile',
  live_update=[
    sync('./backend', '/app/backend')
  ],
)

k8s_yaml(
  helm(
    '../helm',
    name='glitchr',
    namespace='glitchr',
    values=['k8s/helm/values.local.yaml'],
  )
)