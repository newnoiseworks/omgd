package config

func GameConfig(environment string, buildPath string) {
	BuildTemplatesFromPath("game", environment, buildPath)
}
