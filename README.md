# gotmpl
Template generation tools stands for GOlang TeMPLate.

Maybe [Mustache for Golang](https://github.com/cbroglie/mustache) is already a good solution for template generation, BUT, it is not perfect because it has not folder name including tempalte palceholder, more binding data input format, comment for binding data and etc.

#feature
- using [Mustache](https://mustache.github.io/) sytanx for template file
- directory name can be used as template
- json file, yaml file and input arguments as binding data
- comment for binding data
- may be more...

# prequire
```sh
go get github.com/cbroglie/mustache/...
```

# install
```sh
go get github.com/tomjamescn/gotmpl
```

# use
```sh
#using test_data including in gotmpl src directory
cd $GOPATH/github/tomjamescn/gotmpl
gotmpl -bindingDataPath `pwd`/test_data/test.json -templatePath `pwd`/test_data -outputPath /tmp/gotmpl/output/test
```


# thanks
- [Mustache](https://mustache.github.io/)
- [Mustache for Golang](https://github.com/cbroglie/mustache)

