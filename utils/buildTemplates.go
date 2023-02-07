package utils

import (
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

func getData(environment string, buildPath string, verbose bool) *map[interface{}]interface{} {
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

	resourceDir = fmt.Sprintf("%s/resources/", resourceDir)
	_, err := os.Stat(resourceDir)

	if !os.IsNotExist(err) {
		err = filepath.Walk(resourceDir, func(tmpl string, info fs.FileInfo, err error) error {
			if err != nil {
				if verbose {
					log.Printf("no resources directory found in %v", resourceDir)
				}
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
			if verbose {
				log.Println(err)
			}
		}
	}

	return &fp
}

func BuildTemplatesFromPath(environment string, buildPath string, templateExtension string, removeTemplateAfterProcessing bool, verbose bool) {
	fp := getData(environment, buildPath, verbose)

	if verbose {
		log.Println(fmt.Sprintf("building template files in %s", buildPath))
	}

	err := filepath.Walk(buildPath, func(tmpl string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		name := info.Name()

		if info.IsDir() == false && strings.HasSuffix(name, "."+templateExtension) {
			processTemplate(tmpl, fp, templateExtension, removeTemplateAfterProcessing, verbose)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func BuildTemplateFromPath(tmplPath string, environment string, buildPath string, templateExtension string, removeTemplateAfterProcessing bool, verbose bool) {
	fp := getData(environment, buildPath, verbose)
	processTemplate(tmplPath, fp, templateExtension, removeTemplateAfterProcessing, verbose)
}

func processTemplate(tmpl string, fp *map[interface{}]interface{}, templateExtension string, removeTemplateAfterProcessing bool, verbose bool) {
	final_path := strings.ReplaceAll(tmpl, "."+templateExtension, "")

	if verbose {
		log.Println(fmt.Sprintf(" >> build %s >> %s", tmpl, final_path))
	}

	t := template.New(path.Base(tmpl)).Funcs(template.FuncMap{
		"md5":        StrToMd5,
		"upperSnake": StrToUpperSnake,
		"camel":      StrToCamel,
	})

	if templateExtension == "omgdtpl" {
		t.Delims("{*", "*}")
	}

	t, err := t.ParseFiles(tmpl)

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

	if removeTemplateAfterProcessing {
		err = os.Remove(tmpl)
		if err != nil {
			log.Fatal(err)
		}
	}
}
