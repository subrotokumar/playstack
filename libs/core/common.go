package core

type Env string

const (
	EnvDevelopment Env = "dev"
	EnvStaging     Env = "stage"
	EnvProduction  Env = "prod"
)
