# Misc 常用小中间组件
定义了以系列的常用小中间件

## 组件列表

* ElapsedTime 请求执行时间中间件。在返回头里加入请求的执行时间
* Headers 响应头中间件。在返回头内加入固定的响应头
* 逻辑判断中间件，根据给定的参数或者判断函数

## 使用方式

### ElapsedTime

使用该组件后，程序执行时间会加在响应头的 "Elapsed-Time" 中
    
    app.Use(ElapsedTime)

### Headers

配置方式:

    #TOML版本，其他版本可以根据对应格式配置
    "Header1"=["key1"]
    "Header2"=["key2","key3"]

使用:
    m:=&misc.Headers{}
    err=toml.Unmarshal(data,m)
    app.Use(m)

### 逻辑判断中间件

静态判断

    app.Use(
        //启动时条件为真执行动作(http.HandlerFunc)
        misc.If(true,action1),
        //启动时条件为真返回指定状态码的错误
        misc.ErrorIf(true,404),
        //启动时条件为真加入指定的中间件
        misc.MiddlewareIf(true,middleware1)
    )

动态判断

    //创建条件判断函数
    condition=func() (bool, error){
        return time.Now().Before(starttime),nil
    }

    app.Use(
        //运行时条件为真执行动作(http.HandlerFunc)
        misc.When(condition,action1),
        //运行时条件为真返回指定状态码的错误
        misc.ErrorWhen(condition,404),
        //运行时条件为真加入指定的中间件
        misc.MiddlewareWhen(condition,middleware1)
    )