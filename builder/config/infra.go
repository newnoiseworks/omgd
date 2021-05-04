package config

func InfraConfig(environment string, buildPath string) {
	BuildTemplatesFromPath("server/infra", environment, buildPath)
}
