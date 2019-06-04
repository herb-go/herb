# Go Template Golang Template渲染引擎

基于golang 自带 template模块的渲染引擎

## 使用方式

    engine:=gotemplate.Engine
    //注册模板中使用的函数
    engine.RegisterFunc("fn",func(s string)string{
        return s
    })

	oc := NewOptionCommon()
	oc.Engine = engine
    //设置视图根路径
    oc.ViewRoot="/tmp/views"

    err=oc.ApplyTo(renderer)
