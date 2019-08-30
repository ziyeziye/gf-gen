package template

var ModelTmpl = `package {{.PackageName}}

import (
	"database/sql"
)

var {{.StructName}}Model = Model({{.StructName}}{})

type {{.StructName}} struct {
    {{range .Fields}}{{.}}
    {{end}}
}

// TableName sets the insert table name for this struct type
func ({{.ShortStructName}} {{.StructName}}) TableName() string {
	return GetPrefix() + "{{.TableName}}"
}

func Get{{.StructName}}(id int) ({{.StructName | toLower}} *{{.StructName}}, err error) {
	err = {{.StructName}}Model.Where("id", id).Struct(&{{.StructName | toLower}})
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}

func Get{{pluralize .StructName}}(maps map[string]interface{}) ({{pluralize .StructName | toLower}} []*{{.StructName}}, err error) {
	query, err := ModelSearch({{.StructName}}{}, maps)
	if err != nil {
		return
	}
	err = query.Structs(&{{pluralize .StructName | toLower}})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return
}

func Get{{.StructName}}Total(maps map[string]interface{}) (count int, err error) {
	cond, values, err := whereBuild(maps)
	if err != nil {
		return 0, err
	}
	count, err = {{.StructName}}Model.Filter().Where(cond, values...).Count()

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return
}

func Delete{{.StructName}}({{.StructName | toLower}} *{{.StructName}}) (err error) {
	_, err = {{.StructName}}Model.Where({{.StructName | toLower}}).Delete()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}

func Update{{.StructName}}({{.StructName | toLower}} *{{.StructName}}, data map[string]interface{}) (err error) {
	_, err = {{.StructName}}Model.Where({{.StructName | toLower}}).Data(data).Update()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}

func Add{{.StructName}}({{.StructName | toLower}} *{{.StructName}}) (err error) {
	_, err = {{.StructName}}Model.Data({{.StructName | toLower}}).Insert()
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return
}

`
