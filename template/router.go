package template

var RouterTmpl = `package api

import "github.com/gogf/gf/g/net/ghttp"

func ConfigRouter(router *ghttp.RouterGroup) {
    {{range .}}Config{{pluralize .}}Router(router)
    {{end}}
}
`
