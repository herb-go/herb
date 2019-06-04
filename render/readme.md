# Render Http响应渲染器

定义了常用的http响应输出，以及用于渲染的引擎

## 功能
* 常用的Html/json输出帮助函数
* 渲染引擎驱动接口
* 便于序列化为配置文件的视图管理

## 使用方式

### 直接输出

模块中提供了一些方表使用的标准输出工具

    //直接将对象序列化为JSON输出
    //Content-type将设为"application/json"
    render.MustJSON(w,v,200)

    //将数据作为HTML输出。
    //Content-type将设为"text/html"
    render.MustHTML(w,data,200)

    //将文件内容作为HTML输出
    //Content-type将设为"text/html"
    render.MustHTMLFile(w,"/tmp/file.html",200)

    //输出标准的HTML状态信息
    //Content-type将设为"text/plain"
    render.MustError(w,401)

### Renderer 渲染器

渲染器是渲染的主要工具

    renderer=render.New()

    //调用render库的直接输出
    renderer.MustJSON(w,v,200)
    renderer.MustHTML(w,data,200)
    renderer.MustHTMLFile(w,"/tmp/file.html",200)
    renderer.MustError(w,401)

    //渲染器配置
    oc:=render.NewOptionCommon()
    oc.Engine=jetengine
    oc.viewroot="/tmp"
    err=oc.ApplyTo(renderer)

### 视图
    
视图是渲染器的基本输出单元

    //创建视图配置
    viewconfig:=render.NewViewConfig("1.tmpl","2.tmpl")
    //视图可以先获取后初始化
    view:=renderer.GetView("viewname)
    //通过配置初始化视图
    view:=renderer.NewView("viewname",viewconfig)
    
    //直接渲染为html输出
    //Content-type将设为"text/html"
    view.MustRender(w,data)

    //渲染为指定状态码的html输出
    //Content-type将设为"text/html"
    view.MustRenderError(w,data,200)

    //将视图渲染为[]byte数据，再手动处理
    byteSlice=view.MustRenderBytes(data)

### 视图配置文件

    视图配置文件可以通过配置文件的方式批量设置视图

    #TOML版本，其他版本可以根据对应格式配置
    #全局开发模式开关。设为开发模式的视图每次都会重新渲染
    DevelopmentMode=false
    [Views.index]
    #开发者模式开关
    DevelopmentMode=true
    #视图列表
    Files=["views/index.tmpl"]
    [Views.news]
    #开发者模式开关
    DevelopmentMode=true
    #视图列表
    Files=["views/news.tmpl"

使用方式

    views:=&ViewsOptionCommon{}
    err=toml.Unmarshal(data,views)
    views.Init(renderer)

### 视图数据

视图数据为一个通用的传递给视图的数据结构

使用方式:

    data:=render.NewData()

    //设置数据
    data.Set("value1","123)
    //获取数据
    v:=data.Get("value1")
    //删除数据
    data.Del("value1")

    //合并数据
    data.Merge(data1)

### 可用渲染引擎

* [gotemplate](engines/gotemplate) 基于golang http template的渲染引擎
* [jet]((engines/jet) 基于 github.com/CloudyKit/jet 的渲染引擎