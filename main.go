package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/droundy/goopt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jimsmart/schema"
	"github.com/jinzhu/inflection"
	_ "github.com/lib/pq"
	"github.com/serenize/snaker"
	gtmpl "github.com/ziyeziye/gf-gen/template"
	"github.com/ziyeziye/gf-gen/util"
)

var (
	sqlType     = goopt.String([]string{"--sqltype"}, "mysql", "sql database type such as mysql, postgres, etc.")
	sqlConnStr  = goopt.String([]string{"-c", "--connstr"}, "nil", "database connection string")
	sqlDatabase = goopt.String([]string{"-d", "--database"}, "nil", "Database to for connection")
	sqlTable    = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
	prefix      = goopt.String([]string{"--prefix"}, "", "Table prefix")

	packageName = goopt.String([]string{"--package"}, "", "name to set for package")

	jsonAnnotation = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
	gormAnnotation = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
	gureguTypes    = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")

	rest = goopt.Flag([]string{"--rest"}, []string{}, "Enable generating RESTful api", "")

	version = goopt.Flag([]string{"-v", "--version"}, []string{}, "Enable version output", "")
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "ORM and RESTful API generator for Mysql"
	}
	goopt.Version = "v1.2"
	goopt.Summary = `gf-gen [-v] --connstr "user:password@/dbname" --package pkgName --database databaseName --table tableName [--json] [--guregu]`
	//Parse options
	goopt.Parse(nil)
}

func main() {
	if *version == true {
		fmt.Println("The version number of framework-gen is " + goopt.Version)
		return
	}

	// Username is required
	if sqlConnStr == nil || *sqlConnStr == "" {
		fmt.Println("sql connection string is required! Add it with --connstr=s")
		return
	}

	if sqlDatabase == nil || *sqlDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}

	var db, err = sql.Open(*sqlType, *sqlConnStr)
	if err != nil {
		fmt.Println("Error in open database: " + err.Error())
		return
	}
	defer db.Close()

	// parse or read tables
	var tables []string
	if *sqlTable != "" {
		tables = strings.Split(*sqlTable, ",")
	} else {
		tables, err = schema.TableNames(db)
		if err != nil {
			fmt.Println("Error in fetching tables information from mysql information schema")
			return
		}
	}
	// if packageName is not set we need to default it
	if packageName == nil || *packageName == "" {
		*packageName = "framework"
	}
	os.Mkdir("model", 0777)

	apiName := "api"
	if *rest {
		os.Mkdir(apiName, 0777)
	}

	t, err := getTemplate(gtmpl.ModelTmpl)
	if err != nil {
		fmt.Println("Error in loading models template: " + err.Error())
		return
	}

	ct, err := getTemplate(gtmpl.ControllerTmpl)
	if err != nil {
		fmt.Println("Error in loading controller template: " + err.Error())
		return
	}

	var structNames []string

	// generate go files for each table
	for _, tableName := range tables {
		tableName := strings.Replace(tableName, *prefix, "", 1)
		structName := util.FmtFieldName(tableName)
		structName = inflection.Singular(structName)
		structNames = append(structNames, structName)

		modelInfo := util.GenerateStruct(db, *prefix, tableName, structName, "model", *jsonAnnotation, *gormAnnotation, *gureguTypes)
		//fmt.Printf("%+v",tableName)
		//os.Exit(0)
		var buf bytes.Buffer
		err = t.Execute(&buf, modelInfo)
		if err != nil {
			fmt.Println("Error in rendering models: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join("model", inflection.Singular(tableName)+".mod.go"), data, 0777)

		if *rest {
			//write api
			buf.Reset()
			err = ct.Execute(&buf, map[string]string{
				"Package":     *packageName,
				"PackageName": *packageName + "/model",
				"StructName":  structName,
				"TableName":   tableName,
			})
			if err != nil {
				fmt.Println("Error in rendering controller: " + err.Error())
				return
			}
			data, err = format.Source(buf.Bytes())
			if err != nil {
				fmt.Println("Error in formating source: " + err.Error())
				return
			}
			ioutil.WriteFile(filepath.Join(apiName, inflection.Singular(tableName)+".api.go"), data, 0777)
		}
	}

	if *rest {
		rt, err := getTemplate(gtmpl.RouterTmpl)
		if err != nil {
			fmt.Println("Error in lading router template")
			return
		}
		var buf bytes.Buffer
		err = rt.Execute(&buf, structNames)
		if err != nil {
			fmt.Println("Error in rendering router: " + err.Error())
			return
		}
		data, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println("Error in formating source: " + err.Error())
			return
		}
		ioutil.WriteFile(filepath.Join(apiName, "router.go"), data, 0777)
	}
}

func getTemplate(t string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"pluralize":        inflection.Plural,
		"title":            strings.Title,
		"toLower":          strings.ToLower,
		"toLowerCamelCase": camelToLowerCamel,
		"toLowerUnderline": inflection.Singular,
		"toSnakeCase":      snaker.CamelToSnake,
	}

	tmpl, err := template.New("model").Funcs(funcMap).Parse(t)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func camelToLowerCamel(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToLower(ss[0])

	return strings.Join(ss, "")
}
