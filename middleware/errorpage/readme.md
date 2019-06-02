# ErrorPage 自定义错误页组件
根据http响应的状态码，显示自定义的页面

## 功能
* 根据不同的状态码显示不同的页面
* 提供默认的错误页，在非正常状态码(>=400)的情况下默认显示
* 提供禁用组件的功能，使得能在使用了本组件的路由子组件下能显示原始的反馈

## 使用方法
    //创建新的组件
    em:=errorpage.New()
    
    
    em.
        //默认错误页，当状态码>399并<600时使用
        OnError(func(w http.ResponseWriter, r *http.Request, status int){
            http.Error(w,http.StatusText(status),status)
        }).
        //指定状态码的错误页
        OnStatus(404,func(w http.ResponseWriter, r *http.Request, status int){
            http.Error(w,"页面未找到",404)
        }).
        //跳过执行状态
        IgnoreStatus(422)

        app.Use(em)
        
        //强制关闭自定义错误页
        em2.Use(em.MiddlewareDisable())