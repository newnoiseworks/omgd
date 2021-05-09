package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func getData(environment string, buildPath string) *map[interface{}]interface{} {
	fp := make(map[interface{}]interface{})
	fp["profile"] = GetProfile(environment).GetProfileAsMap()

	resourceDir := buildPath

	if strings.HasPrefix(environment, "..") {
		resourceDir = strings.Split(environment, "/profiles")[0]
	}

	err := filepath.Walk(fmt.Sprintf("%s/resources/", resourceDir), func(tmpl string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		name := info.Name()

		if info.IsDir() == false && strings.HasSuffix(name, ".yml") {
			c := make(map[interface{}]interface{})

			yamlFile, err := ioutil.ReadFile(tmpl)
			if err != nil {
				log.Printf("yamlFile Get err: #%v ", err)
			}

			err = yaml.Unmarshal(yamlFile, &c)
			if err != nil {
				log.Fatalf("Unmarshal err: %v", err)
			}

			for k, v := range c {
				fp[k] = v
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return &fp
}

func BuildTemplatesFromPath(path string, environment string, buildPath string) {
	fp := getData(environment, buildPath)

	fmt.Println(fmt.Sprintf("build %s config files", path))
	fmt.Println(fmt.Sprintf("%s/%s", buildPath, path))

	err := filepath.Walk(fmt.Sprintf("%s/%s", buildPath, path), func(tmpl string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		name := info.Name()

		if info.IsDir() == false && strings.HasSuffix(name, ".tmpl") {
			processTemplate(tmpl, name, fp)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func processTemplate(tmpl string, name string, fp *map[interface{}]interface{}) {
	final_path := strings.ReplaceAll(tmpl, ".tmpl", "")

	fmt.Println(fmt.Sprintf(" >> build %s >> %s", tmpl, final_path))

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
		log.Fatal(err)
	}

	f, err := os.Create(final_path)
	if err != nil {
		log.Fatal("create file: ", err)
	}

	err = t.Execute(f, fp)
	if err != nil {
		log.Fatal("execute: ", err)
	}
}
