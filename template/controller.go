package template

var ControllerTmpl = `package api

import (
	"framework/library/request"
	"framework/library/response"
	"net/http"

	"github.com/gogf/gf/g/frame/gmvc"
	"github.com/gogf/gf/g/net/ghttp"

	"framework/app/model"
)

type {{.StructName}}Api struct {
	gmvc.Controller
}

func Config{{pluralize .StructName}}Router(router *ghttp.RouterGroup) {
	controller := {{.StructName}}Api{}
	router.GET("/{{pluralize .StructName | toLower}}", controller.GetAll{{pluralize .StructName}})
	router.POST("/{{pluralize .StructName | toLower}}", controller.Add{{.StructName}})
	router.GET("/{{pluralize .StructName | toLower}}/:id", controller.Get{{.StructName}})
	router.PUT("/{{pluralize .StructName | toLower}}/:id", controller.Update{{.StructName}})
	router.DELETE("/{{pluralize .StructName | toLower}}/:id", controller.Delete{{.StructName}})
}

func (c *{{.StructName}}Api) GetAll{{pluralize .StructName}}(r *ghttp.Request) {
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	var total int
	total, _ = model.Get{{.StructName}}Total(maps)
	data["total"] = total

	maps = request.GetPage(r, maps, false)
	respJson := response.Json(r)
	if {{pluralize .StructName | toLower}}, err := model.Get{{pluralize .StructName}}(maps); err != nil {
		respJson.SetState(false).SetMsg("error")
	} else {
		data["list"] = {{pluralize .StructName | toLower}}
		respJson.SetData(data)
	}

	respJson.Return()
}

func (c *{{.StructName}}Api) Get{{.StructName}}(r *ghttp.Request) {
	id := r.GetInt("id")
	{{.StructName | toLower}}, err := model.Get{{.StructName}}(id)

	respJson := response.Json(r)
	if err == nil && {{.StructName | toLower}}.ID > 0 {
		respJson.SetData({{.StructName | toLower}})
	} else {
		respJson.SetState(false).SetCode(response.ERROR_NOT_EXIST)
	}
	respJson.Return()
}

func (c *{{.StructName}}Api) Add{{.StructName}}(r *ghttp.Request) {
	//maps := make(map[string]interface{})

	{{.StructName | toLower}} := model.{{.StructName}}{}

	respJson := response.Json(r)
	if err := model.Add{{.StructName}}(&{{.StructName | toLower}}); err != nil {
		respJson.Set(http.StatusInternalServerError, "新增失败", false, {{.StructName | toLower}})
	} else {
		respJson.SetData({{.StructName | toLower}})
	}
	respJson.Return()
}

func (c *{{.StructName}}Api) Update{{.StructName}}(r *ghttp.Request) {
	id := r.GetInt("id")
	maps := make(map[string]interface{})

	{{.StructName | toLower}}, err := model.Get{{.StructName}}(id)

	respJson := response.Json(r)
	if err == nil && {{.StructName | toLower}}.ID > 0 {
		if err := model.Update{{.StructName}}({{.StructName | toLower}}, maps); err != nil {
			respJson.Set(http.StatusInternalServerError, "修改失败", false, {{.StructName | toLower}})
		} else {
			respJson.SetData({{.StructName | toLower}})
		}
	} else {
		respJson.SetState(false).SetCode(response.ERROR_NOT_EXIST)
	}
	respJson.Return()
}

func (c *{{.StructName}}Api) Delete{{.StructName}}(r *ghttp.Request) {
	id := r.GetInt("id")
	{{.StructName | toLower}}, err := model.Get{{.StructName}}(id)

	respJson := response.Json(r)
	if err != nil {
		respJson.SetState(false).SetCode(response.ERROR_NOT_EXIST)
	}

	if err := model.Delete{{.StructName}}({{.StructName | toLower}}); err != nil {
		respJson.Set(http.StatusInternalServerError, "删除失败", false, nil)
	}
	respJson.Return()
}

`
