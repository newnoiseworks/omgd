package config

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/newnoiseworks/tpl-fred/utils"
)

// GameConfig go brr
func GameConfig(environment string) {
	fmt.Println("build game config")

	var profile = utils.GetProfile(environment)

	fmt.Println(profile)

	t, err := template.ParseFiles("builder/config/templates/GameConfig.cs.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := ".tmp/GameConfig.cs"

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	config := map[string]string{
		"realWorldSecondsPerDay": profile.Game.RealWorldSecondsPerDay,
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
