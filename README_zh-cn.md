# gotmpl
模板工具, 将软件开发中的模板模式应用于任意文件(夹), 如果将重复工作都模板化，将大大提高重复工作的生产效率. 名字来源于英文 GOlang TeMPLate.

Maybe [Mustache for Golang](https://github.com/cbroglie/mustache) is already a good solution for single template generation, BUT, it is not perfect because it does not support folder name including placeholder, more binding data input format, comment for binding data and etc.
[Mustache for Golang](https://github.com/cbroglie/mustache)已经是一个比较好的模板工具了，但是它只支持单个文件作为模板，远没有完美，完美的模板工具应该支持模板文件路径中包含占位符，更多的绑定数据格式，占位符注释，占位符概览等等

# 特性
- [x] 使用[Mustache](https://mustache.github.io/) 语法作为模板语言
- [x] 支持莫办文件路径包含占位符
- [x] 绑定数据支持json文件，yaml文件和命令行参数格式
- [x] 占位符概览
- [x] 占位符注释
- [x] git代码库作为模板路径, 支持https:// http:// git:// ssh:// 协议
- 更多...

# 安装
```sh
go get github.com/tomjamescn/gotmpl
```

# 使用
```sh
#使用本项目的test_data子目录作为模板路径
cd $GOPATH/github/tomjamescn/gotmpl

#使用json格式的绑定数据
gotmpl -b `pwd`/test_data/test.json -t `pwd`/test_data -o /tmp/gotmpl/output/test

#使用命令行参数的绑定数据
gotmpl -t `pwd`/test_data -o /tmp/gotmpl/output/test index=index

#显示占位符概览
gotmpl -t `pwd`/test_data -s 

#使用git代码库作为模板路径
#如果使用私有库, 请确保已使用ssh-copy-id命令将公钥加入到了目标机器中
gotmpl -t https://github.com/tomjamescn/gotmpl -s test_data -o /tmp/output index=TEST dir_name=TEST

```

# 感谢
- [Mustache](https://mustache.github.io/)
- [Mustache for Golang](https://github.com/cbroglie/mustache)

