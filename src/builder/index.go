package builder

// Builder it's great
func Builder(project string, environment string) {
	switch project {
	case "game":
		BuildGame(environment)
		break
	}
}
