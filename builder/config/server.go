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

	"github.com/logrusorgru/aurora"
)

// ServerConfig go brr
func ServerConfig(environment string, buildPath string) {
	var profile = utils.GetProfile(environment)

	config := map[string]string{
		"realWorldSecondsPerDay": profile.Game.RealWorldSecondsPerDay,
		"nakamaHost":             profile.Nakama.Host,
		"nakamaKey":              profile.Nakama.Key,
		"nakamaPort":             strconv.Itoa(profile.Nakama.Port),
		"nakamaSecure":           strconv.FormatBool(profile.Nakama.Secure),
		"gameVersion":            profile.Game.Version,
	}

	fmt.Println(aurora.Green("building server config files"))

	buildServerGameConfig(buildPath, config)
	buildServerConfig(buildPath, config)
	buildServerVersionFile(buildPath, config)
	buildServerItemsFile(buildPath)
	buildServerMissionList(buildPath)
}

func buildServerGameConfig(buildPath string, config map[string]string) {
	fmt.Println(aurora.Yellow(" >> building game_config.lua.tmpl >> server/nakama/data/modules/game_config.lua"))

	t, err := template.ParseFiles("builder/config/templates/game_config.lua.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/server/nakama/data/modules/game_config.lua", buildPath)

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

func buildServerItemsFile(buildPath string) {
	var items = utils.GetItems()

	fmt.Println(" >> build inventory_items.lua.tmpl >> /server/nakama/data/modules/inventory_items.lua")

	var tmpl = "builder/config/templates/inventory_items.lua.tmpl"
	t, err := template.New(path.Base(tmpl)).Funcs(template.FuncMap{"md5": func(text string) string {
		hash := md5.Sum([]byte(text))
		return hex.EncodeToString(hash[:])
	}}).ParseFiles(tmpl)
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/server/nakama/data/modules/inventory_items.lua", buildPath)

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

func buildServerConfig(buildPath string, config map[string]string) {
	fmt.Println(" >> build config.yml.tmpl >> server/nakama/data/config.yml")
	t, err := template.ParseFiles("builder/config/templates/config.yml.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/server/nakama/data/config.yml", buildPath)

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

func buildServerVersionFile(buildPath string, config map[string]string) {
	fmt.Println(" >> build version.lua.tmpl >> /server/nakama/data/modules/version.lua")

	t, err := template.ParseFiles("builder/config/templates/version.lua.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/server/nakama/data/modules/version.lua", buildPath)

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

func buildServerMissionList(buildPath string) {
	var missions = utils.GetMissions()

	fmt.Println(" >> build mission_list.lua.tmpl >> /server/nakama/data/modules/mission_list.lua")

	var tmpl = "builder/config/templates/mission_list.lua.tmpl"
	t, err := template.New(path.Base(tmpl)).Funcs(template.FuncMap{"md5": func(text string) string {
		hash := md5.Sum([]byte(text))
		return hex.EncodeToString(hash[:])
	}}).ParseFiles(tmpl)
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/server/nakama/data/modules/mission_list.lua", buildPath)

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, missions)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}
