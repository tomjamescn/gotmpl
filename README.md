# gotmpl
Template generation tools stands for GOlang TeMPLate.

Maybe [Mustache for Golang](https://github.com/cbroglie/mustache) is already a good solution for single template generation, BUT, it is not perfect because it does not support folder name including placeholder, more binding data input format, comment for binding data and etc.

# feature
- [x] using [Mustache](https://mustache.github.io/) sytanx for template file
- [x] support placeholder in folder name
- [x] json file, yaml file and input arguments as binding data
- [ ] summary for all placeholder
- [ ] comment support for placeholder
- [ ] git repository as template path
- may be more...

# install
```sh
go get github.com/tomjamescn/gotmpl
```

# usage
```sh
#using test_data including in gotmpl src directory
cd $GOPATH/github/tomjamescn/gotmpl
gotmpl -b `pwd`/test_data/test.json -t `pwd`/test_data -o /tmp/gotmpl/output/test

#using input arguments as binding data
gotmpl -t `pwd`/test_data -o /tmp/gotmpl/output/test index=index

#show summary for all placeholder
gotmpl -t `pwd`/test_data -s 

```

# thanks
- [Mustache](https://mustache.github.io/)
- [Mustache for Golang](https://github.com/cbroglie/mustache)

