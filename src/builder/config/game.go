package config

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/newnoiseworks/tpl-fred/utils"
)

// GameConfig go brr
func GameConfig(environment string, buildPath string) {
	var profile = utils.GetProfile(environment)

	config := map[string]string{
		"realWorldSecondsPerDay": profile.Game.RealWorldSecondsPerDay,
		"nakamaHost":             profile.Nakama.Host,
		"nakamaKey":              profile.Nakama.Key,
		"nakamaPort":             strconv.Itoa(profile.Nakama.Port),
		"nakamaSecure":           strconv.FormatBool(profile.Nakama.Secure),
		"websiteHost":            profile.Website.Host,
	}

	buildGameClientConfig(buildPath, config)
	buildServerConfig(buildPath, config)
	buildGameBuildConfig(environment, buildPath, config)
}

func buildGameClientConfig(buildPath string, config map[string]string) {
	fmt.Println("build game config")

	t, err := template.ParseFiles("builder/config/templates/GameConfig.cs.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/GameConfig.cs", buildPath)

	f, err := os.Create(path)

	if err != nil {
		log.Println("create file: ", err)
		return
	}
	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}

func buildServerConfig(buildPath string, config map[string]string) {
	fmt.Println("build game config for server")
	t, err := template.ParseFiles("builder/config/templates/game_config.lua.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game_config.lua", buildPath)

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}

func buildGameBuildConfig(environment string, buildPath string, config map[string]string) {
	fmt.Println("build game config for game client build")
	t, err := template.ParseFiles("builder/config/templates/config.tpl_build.tres.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/config.tpl_%s.tres", buildPath, environment)

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
