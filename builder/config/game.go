package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"text/template"

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
		"gameVersion":            profile.Game.Version,
	}

	fmt.Println("build game config files")

	buildGameClientConfig(buildPath, config)

	// buildGameClientConfigTpl(buildPath, config)

	// buildGameBuildConfig(environment, buildPath, config)

	buildGameItemsFile(buildPath)
}

func buildGameClientConfig(buildPath string, config map[string]string) {
	fmt.Println(" >> build GameConfig.cs.tmpl >> game/Resources/Config/GameConfig.cs")

	t, err := template.ParseFiles("builder/config/templates/GameConfig.cs.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Resources/Config/GameConfig.cs", buildPath)

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

func buildGameClientConfigTpl(buildPath string, config map[string]string) {
	fmt.Println(" >> build config.tpl.tres.tmpl >> game/Resources/Config/config.tpl.tres")

	t, err := template.ParseFiles("builder/config/templates/config.tpl.tres.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Resources/Config/config.tpl.tres", buildPath)

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
	fmt.Printf(" >> build config.tpl_build.tres.tmpl >> game/Resources/Config/config.tpl_%s.tres\n", environment)
	t, err := template.ParseFiles("builder/config/templates/config.tpl_build.tres.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Resources/Config/config.tpl_%s.tres", buildPath, environment)

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

func buildGameItemsFile(buildPath string) {
	var items = utils.GetItems()

	fmt.Println(" >> build InventoryItems.cs.tmpl >> game/Data/InventoryItems.cs")

	var tmpl = "builder/config/templates/InventoryItems.cs.tmpl"
	t, err := template.New(path.Base(tmpl)).Funcs(template.FuncMap{"md5": func(text string) string {
		hash := md5.Sum([]byte(text))
		return hex.EncodeToString(hash[:])
	}}).ParseFiles(tmpl)
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Data/InventoryItems.cs", buildPath)

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, items)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
