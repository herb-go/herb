# PageCache 页面缓存组件

提供将整个http响应进行缓存功能的组件

## 功能
* 缓存后的页面能够继续添加或者维护Header
* 通过缓存的Loader功能防止缓存雪崩
* 能状态码进行判断页面是否需要需要缓存

## 使用方式

### 通过传入缓存创建PageCache对象，再通过PageCache对象创建相应的中间件

1.创建pagecache

    pc:=pagecache.New(cache)
    
2.创建合适的主键生成器。主键生成器根据传入的Http请求生成缓存的主键。如果返回的主键为空则该页面不进行缓存。

    kg:=func(r *http.Request) string(
        uid:=r.Header.Get("uid")
        if uid!="" {
            return ""
        }
        return "news-"+r.URL.Query().Get("id")
    )

3.创建中间件并使用

    newspagecache:=pagecache.Middleware(kg,time.Hour)
    app.Use(newspagecache)

###直接通过FieldGenerator创建中间件

1.创建FieldGenerator。FieldGeneratorg根据传入的http请求生成响应的缓存字段。如果缓存字段为空则该页面不进行缓存。

    
    fg=func(r *http.Request) *cache.Field{
        uid:=r.Header.Get("uid")
        if uid!="" {
            return nil
        }
        return cache.Field("news-"+r.URL.Query().Get("id"))
    }

2.通过FieldGenerator创建中间件

    newspagecache:=FieldMiddleware(fg,time.Hour,nil)
    app.Use(newspagecache)

### 请求的状态与缓存的关系

默认情况下，页面缓存组件对且只对StatusCode<500的请求进行缓存处理

如果需要对状态码进行更信息的操作，可以通过自定义状态验证器来进行处理

    statusvalidator:=func(status int) bool {
	    return status == 200 || status ==301 || status == 302
    }
    //调整pagecache对象属性
    newspagecache.StatusValidator=statusvalidator

    //通过FieldGenerator
    newspagecache:=FieldMiddleware(fg,time.Hour,statusvalidator)