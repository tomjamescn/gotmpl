// Package main provides ...
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/cbroglie/mustache"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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
			filePaths, err := listTemplateFilePath(path.Join(templatePath, f.Name()))
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

func StringStartWith(str string, subStrs []string, caseSensitive bool) (string, bool) {
	str = strings.Trim(str, " \n\t")
	if !caseSensitive {
		str = strings.ToLower(str)
	}
	for _, sub := range subStrs {
		prefix := sub
		if !caseSensitive {
			prefix = strings.ToLower(prefix)
		}
		if strings.HasPrefix(str, prefix) {
			return sub, true
		}
	}

	return "", false
}

func NormalizePath(path string) string {
	r := regexp.MustCompile("/{2,}")
	return string(r.ReplaceAll([]byte(path), []byte("")))
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getRealTemplatePath(templatePath, subTemplatePath string) (p string, tmpPath string, isRepository bool) {
	if _, ok := StringStartWith(templatePath, []string{"http://", "https://", "file://", "git://", "ssh://"}, true); ok {
		//repository

		tmpDir := GetMD5Hash(templatePath)
		tmpPath = path.Join(os.TempDir(), string(tmpDir))
		cloneOptions := &git.CloneOptions{
			URL:      templatePath,
			Progress: os.Stdout,
		}
		fmt.Println(tmpPath)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			panic(err)
		}
		_, err = git.PlainClone(tmpPath, false, cloneOptions)
		if err != nil {
			err = os.RemoveAll(tmpPath)
			if err != nil {
				panic(err)
			}
			s := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
			sshKey, err := ioutil.ReadFile(s)
			signer, err := ssh.ParsePrivateKey([]byte(sshKey))
			auth := &gitssh.PublicKeys{User: "git", Signer: signer}
			cloneOptions.Auth = auth
			_, err = git.PlainClone(tmpPath, false, cloneOptions)
			if err != nil {
				panic(err)
			}
		}

		return path.Join(tmpPath, subTemplatePath), tmpPath, true
	} else {
		return NormalizePath(path.Join(templatePath, subTemplatePath)), "", false
	}
}

type summary struct {
	DirPathList  []string
	FilePathList []string
}

func showPlaceholderSummary(templatePath string) {
	tmpFilePaths, err := listTemplateFilePath(templatePath)
	if err != nil {
		panic(err)
	}

	allPlaceholders := make(map[string]bool, 0)
	placeholders := make(map[string]*summary)
	comments := make(map[string]string)
	r4Placeholder := regexp.MustCompile(`\{\{\s*([^\}\s]*)\s*\}\}`)
	r4Comment := regexp.MustCompile(`\{\{\!\s*([^\}:]*)\s*:\s*([^\}]*)\s*\}\}`)
	for _, p := range tmpFilePaths {
		content, err := ioutil.ReadFile(p)
		if err != nil {
			panic(err)
		}
		retD := r4Placeholder.FindAllStringSubmatch(p, -1)
		retP := r4Placeholder.FindAllStringSubmatch(string(content), -1)
		retC := r4Comment.FindAllStringSubmatch(string(content), -1)

		for _, v := range retD {
			allPlaceholders[v[1]] = true
			s, ok := placeholders[v[1]]
			if !ok {
				s = &summary{}
				s.DirPathList = make([]string, 0)
				s.FilePathList = make([]string, 0)
				placeholders[v[1]] = s
			}
			s.DirPathList = append(s.DirPathList, p)
		}

		for _, v := range retP {
			allPlaceholders[v[1]] = true
			s, ok := placeholders[v[1]]
			if !ok {
				s = &summary{}
				s.DirPathList = make([]string, 0)
				s.FilePathList = make([]string, 0)
				placeholders[v[1]] = s
			}
			s.FilePathList = append(s.FilePathList, p)
		}

		for _, v := range retC {
			comments[v[1]] = v[2]
		}

	}

	for p, _ := range allPlaceholders {
		fmt.Println("placeholder:", p)
		if v, ok := comments[p]; ok {
			fmt.Println("comment:", v)
		} else {
			fmt.Println("comment:")
		}

		fmt.Println("dir path list:")

		if v, ok := placeholders[p]; ok {
			for _, l := range v.DirPathList {
				fmt.Println(l)
			}
		}

		fmt.Println("file list:")

		if v, ok := placeholders[p]; ok {
			for _, l := range v.FilePathList {
				fmt.Println(l)
			}
		}

		fmt.Println("")
	}
}

func main() {
	var templatePath string
	var subTemplatePath string
	var outputPath string
	var bindingDataPath string
	var summaryFlag bool
	var bindingDatasFromArgs map[string]string

	bindingDatasFromArgs = make(map[string]string)
	_, currentFilePath, _, _ := runtime.Caller(0)
	curDirPath := path.Dir(currentFilePath)
	flag.StringVar(&templatePath, "t", "", "template path, can be a folder, a file path or a git repository")
	flag.StringVar(&subTemplatePath, "s", "", "sub template path append to template path, can be a folder, a file path or a git repository")
	flag.StringVar(&outputPath, "o", curDirPath, "output path")
	flag.StringVar(&bindingDataPath, "b", "", "binding data file path, can be a yaml or json file path")
	flag.BoolVar(&summaryFlag, "p", false, "show placeholder summary")
	flag.Parse()

	if templatePath == "" {
		panic("-t can not be empty")
	}

	templatePath, tmpPath, isRepository := getRealTemplatePath(templatePath, subTemplatePath)

	if summaryFlag {
		showPlaceholderSummary(templatePath)
		return
	}

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

	if isRepository {
		fmt.Println(tmpPath)
		err = os.RemoveAll(tmpPath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("finished!")
}
