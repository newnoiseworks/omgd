package config

func ServerConfig(environment string, buildPath string) {
	BuildTemplatesFromPath("server", environment, buildPath)
}
