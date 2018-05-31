> [中文版README.md](https://github.com/tomjamescn/gotmpl/blob/master/README_zh-cn.md)

# gotmpl
Template tools stands for GOlang TeMPLate.

Maybe [Mustache for Golang](https://github.com/cbroglie/mustache) is already a good solution for single template, BUT, it is not perfect because it does not support folder name including placeholder, more binding data input format, comment for placeholder and etc.

# feature
- [x] using [Mustache](https://mustache.github.io/) sytanx for template file
- [x] support placeholder in folder name
- [x] json file, yaml file and input arguments as binding data
- [x] summary for all placeholder
- [x] comment support for placeholder
- [x] git repository as template path, support https:// http:// git:// ssh:// protocals
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

#using repository as template path
#if using private repository, ssh-copy-id must be executed before!
gotmpl -t https://github.com/tomjamescn/gotmpl -s test_data -o /tmp/output index=TEST dir_name=TEST

```

# thanks
- [Mustache](https://mustache.github.io/)
- [Mustache for Golang](https://github.com/cbroglie/mustache)

