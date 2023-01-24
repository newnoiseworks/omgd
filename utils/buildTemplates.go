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

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func getData(environment string, buildPath string) *map[interface{}]interface{} {
	fp := make(map[interface{}]interface{})
	fp["profile"] = GetProfile(environment).GetProfileAsMap()

	resourceDir := buildPath

	if strings.HasPrefix(environment, "..") {
		paths := strings.Split(environment, "./")

		resourceDir = ""

		for _, p := range paths {
			if p == "." {
				resourceDir += "../"
			}
		}
	}

	err := filepath.Walk(fmt.Sprintf("%s/resources/", resourceDir), func(tmpl string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("no resources directory found in %v. in buildpath %v", resourceDir, buildPath)
			return nil
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

func BuildTemplatesFromPath(path string, environment string, buildPath string, templateExtension string, verbose bool) {
	fp := getData(environment, buildPath)

	if verbose {
		log.Println(fmt.Sprintf("build %s config files", path))
		log.Println(fmt.Sprintf("%s/%s", buildPath, path))
	}

	err := filepath.Walk(fmt.Sprintf("%s/%s", buildPath, path), func(tmpl string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		name := info.Name()

		if info.IsDir() == false && strings.HasSuffix(name, "."+templateExtension) {
			processTemplate(tmpl, fp, templateExtension, verbose)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func BuildTemplateFromPath(path string, environment string, buildPath string, templateExtension string, verbose bool) {
	fp := getData(environment, buildPath)
	processTemplate(path, fp, templateExtension, verbose)
}

func processTemplate(tmpl string, fp *map[interface{}]interface{}, templateExtension string, verbose bool) {
	final_path := strings.ReplaceAll(tmpl, "."+templateExtension, "")

	if verbose {
		log.Println(fmt.Sprintf(" >> build %s >> %s", tmpl, final_path))
	}

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
		"camel": func(text string) string {
			return strcase.ToCamel(text)
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
