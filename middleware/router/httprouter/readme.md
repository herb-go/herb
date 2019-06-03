# Httprouter httprouter 路由

基于 github.com/julienschmidt/httprouter 实现的路由

## 使用方式

    Router:=httprouter.New()

    //响应GET请求
    Router.GET("/getpath").
        Use(getmiddlewares...).
        HandlerFunc(getaction)

    //响应POST请求
    Router.POST("/postpath").
        Use(postmiddlewares...).
        HandlerFunc(poataction)

    //响应PUT请求
    Router.PUT("/putpath").
        Use(putmiddlewares...).
        HandlerFunc(putaction)

    //响应DELETE请求
    Router.DELETE("/deletepath").
        Use(deletemiddlewares...).
        HandlerFunc(deleteaction)

    //响应PATCH请求
    Router.PATCH("/patchpath").
        Use(patchmiddlewares...).
        HandlerFunc(patchaction)

    //响应OPTIONS请求
    Router.OPTIONS("/optionspath").
        Use(optionsmiddlewares...).
        HandlerFunc(optionsaction)

    //响应HEAD请求
    Router.HEAD("/deletepath").
        Use(headmiddlewares...).
        HandlerFunc(headaction)

    //响应所有请求
    Router.ALL("/allpath").
        Use(allmiddlewares...).
        HandlerFunc(allaction)

    //子路由
    subrouter:=httprouter.New()
    //被截取的参数会保存在路由参数的filepath里
    Router.StripPrefix("/sub).
         Handler(subrouter)

    //获取路径参数
    Router.GET("/get/:id").
        HandleFunc(func(w http.ResponseWriter, r *http.Request){
            params=router.GetParams(r)
            id:=params.Get("id)
        })