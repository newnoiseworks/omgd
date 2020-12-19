package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
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

	buildGameGDItemsFile(buildPath)

	buildMissionListFile(buildPath)
}

func buildGameClientConfig(buildPath string, config map[string]string) {
	fmt.Println(" >> build GameConfig.gd.tmpl >> game/Utils/GameConfig.gd")

	t, err := template.ParseFiles("builder/config/templates/GameConfig.gd.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Utils/GameConfig.gd", buildPath)

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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func buildGameGDItemsFile(buildPath string) {
	var items = utils.GetItems()

	fmt.Println(" >> build InventoryItems.gd.tmpl >> game/Utils/InventoryItems.gd")

	var tmpl = "builder/config/templates/InventoryItems.gd.tmpl"
	t, err := template.New(path.Base(tmpl)).Funcs(template.FuncMap{
		"md5": func(text string) string {
			hash := md5.Sum([]byte(text))
			return hex.EncodeToString(hash[:])
		},
		"upperSnake": func(text string) string {
			snake := matchFirstCap.ReplaceAllString(text, "${1}_${2}")
			snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
			return strings.ToUpper(snake)
		},
	}).ParseFiles(tmpl)
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Utils/InventoryItems.gd", buildPath)

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

func buildMissionListFile(buildPath string) {
	var missions = utils.GetMissions()

	fmt.Println(" >> build MissionList.gd.tmpl >> game/Utils/MissionList.gd")

	var tmpl = "builder/config/templates/MissionList.gd.tmpl"
	t, err := template.New(path.Base(tmpl)).Funcs(template.FuncMap{
		"upperSnake": func(text string) string {
			snake := matchFirstCap.ReplaceAllString(text, "${1}_${2}")
			snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
			return strings.ToUpper(snake)
		},
		"md5": func(text string) string {
			hash := md5.Sum([]byte(text))
			return hex.EncodeToString(hash[:])
		}}).ParseFiles(tmpl)
	if err != nil {
		log.Print(err)
		return
	}

	path := fmt.Sprintf("%s/game/Utils/MissionList.gd", buildPath)

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
