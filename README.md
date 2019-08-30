## gf-Framework
github.com/ziyeziye/gf-framework
是一个Golang-Go Frame的脚手架

## gf-gen
github.com/ziyeziye/gf-gen
是一个可以通过数据库生成对应gf-Framework的model,struct以及相应的restful api的工具。

## 使用工具时注意代码被覆盖

## Usage

```BASH
-v --version The version number of framework-gen is
--connstr database connection string
--database Database to for connection
--table Table to build struct from
--prefix Table prefix
--package name to set for package,default:framework

go run main.go --connstr "root:pass@tcp(127.0.0.1:3306)/dbname?&parseTime=True" --package github.com/ziyeziye/gf-framework --prefix gf_ --json --guregu --rest

gf-gen --connstr "root:pass@tcp(127.0.0.1:3306)/dbname?&parseTime=True" --prefix gf_ --json --guregu --rest
```
## TODO
1.计划支持获取数据表字段注释/表注释
2.写入接口文档注释，支持swagger接口文档