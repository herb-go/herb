# Middleware 中间件组件
提供一系列中间件操作工具和常用中间件，包括:

### 洋葱模型

中间件实现了洋葱模型。

所有的请求会一次执行next前的部分，然后执行后响应代码会反向执行next之后的部分。

    app.
        Use(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
            //执行顺序1
            DoSomethingBeforeHanlderA()

            //代码执行挂起，运行下一个中间件
             next(w,r)

            //执行顺序5
            DoSomethingAfterHandlerE()
        }).
        Use(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
            //执行顺序2
            DoSomethingBeforeHanlderB()

            //代码执行挂起，运行下一个中间件
             next(w,r)

            //执行顺序4
            DoSomethingAfterHandlerD()

        }).HanldeFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){
            //执行代码处理部分
            //执行顺序3
            DoSomethingC()

            //执行完毕，开始反向执行中间件next之后部分
        })
如果代码中有Panic,会直接跳过所有之后的请求部分和响应部分的代码

## 使用方法

### 使用app对象
本组件提供了一个用于进行中间件操作的类 App,可以用来组合一系列的中间件以供使用

    //创建空app
    app=middlewre.New()
    //根据现有中间件创建app
    app2=middlewre.New(middleware1,middleware2,middleware3)
    app.
        //按顺序引入中间件
        Use(middleware1,middleware2,middleware3)
        //按顺序引入实现了HandlerSlice的类
        .Chain(hanlder1,handler2)
        //按顺序引入其他app
        .UseApp(app2,app3,app4)
        //最终执行http.HanlderFunc
        .HanldeFunc(httphanlderfunc)    
    
    //最终执行http.Hanlder
    app2.Hanlde(httphanler)
    
### 转换中间件

    //将http.Hanlder转换为中间件
    //将执行http.Hanlder的ServerHTTP，再调用下一中间件
    middlewareHttpHanlder:=Wrap(httphanler)

    //将http.HanlderFunc转换为中间件
    //将先执行http.HanlderFunc，再调用下一个中间件
    
    //将app转化为中间件
    miedlewareApp:=AppToMiddlewares(app,app2,app3)
    
## 子组件列表

* 用于快速管理中间件的 Middleware与App对象
* [Basicauth](basicauth) HTTPBasic验证
* [Csrf](csrf) Csrf验证组件
* [Errorpage](errorpage) 自定义错误页中间件
* [Forwarded](forwarded) 被请求转发的信息管理组件
* [Misc](misc) 杂项中间件
* [Router](router) 路由
