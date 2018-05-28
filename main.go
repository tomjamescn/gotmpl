// Package main provides ...
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cbroglie/mustache"
	yaml "gopkg.in/yaml.v2"
)

func listTemplateFilePath(templatePath string) ([]string, error) {
	templateFilePaths := make([]string, 0)
	fileInfo, err := os.Stat(templatePath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		fileInfoes, err := ioutil.ReadDir(templatePath)
		if err != nil {
			panic(err)
		}
		for _, f := range fileInfoes {
			filePaths, err := listTemplateFilePath(templatePath + string(filepath.Separator) + f.Name())
			if err != nil {
				return nil, err
			}
			templateFilePaths = append(templateFilePaths, filePaths...)
		}
	} else {
		templateFilePaths = append(templateFilePaths, templatePath)
	}

	return templateFilePaths, nil
}

func renderAndSave(templateFilePath, templatePath, outputFilePath string, bindingDatas ...interface{}) error {
	text, err := mustache.RenderFile(templateFilePath, bindingDatas...)
	if err != nil {
		return err
	}
	//maybe file path string contains template variable
	templateFilePath = strings.Replace(templateFilePath, filepath.Dir(templatePath), outputFilePath, 1)
	templateFilePath, err = mustache.Render(templateFilePath, bindingDatas...)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(templateFilePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		//be careful, must be 07xx
		err = os.MkdirAll(dirPath, 0744)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(templateFilePath, []byte(text), 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseBindingData(bindingDataPath string) (interface{}, error) {
	b, err := ioutil.ReadFile(bindingDataPath)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if strings.HasSuffix(strings.ToLower(bindingDataPath), "yaml") {
		if err := yaml.Unmarshal(b, &data); err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(strings.ToLower(bindingDataPath), "json") {
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Not Support binding data format")
	}
	return data, nil
}

func main() {
	var templatePath string
	var outputPath string
	var bindingDataPath string
	var bindingDatasFromArgs map[string]string

	bindingDatasFromArgs = make(map[string]string)
	_, currentFilePath, _, _ := runtime.Caller(0)
	curDirPath := path.Dir(currentFilePath)
	flag.StringVar(&templatePath, "templatePath", "", "template path, can be a folder or a file path")
	flag.StringVar(&outputPath, "outputPath", curDirPath, "output path")
	flag.StringVar(&bindingDataPath, "bindingDataPath", "", "data for variables in template, can be a yaml or json file path")
	flag.Parse()

	context := make([]interface{}, 0)
	if bindingDataPath != "" {
		bindingData, err := parseBindingData(bindingDataPath)
		if err != nil {
			panic(err)
		}
		context = append(context, bindingData)
	}

	for _, v := range flag.Args() {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) != 2 {
			panic("bind data must match this pattern:xxx=yyy")
		}
		bindingDatasFromArgs[kv[0]] = kv[1]
	}

	if len(bindingDatasFromArgs) > 0 {
		context = append(context, bindingDatasFromArgs)
	}

	templateFilePaths, err := listTemplateFilePath(templatePath)
	if err != nil {
		panic(err)
	}
	for _, p := range templateFilePaths {
		if err = renderAndSave(p, templatePath, outputPath, context...); err != nil {
			panic(err)
		}
	}

	fmt.Println("finished!")
}
