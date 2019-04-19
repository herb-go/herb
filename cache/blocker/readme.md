# Blocker 请求拦截器
根据指定时间段内，服务器返回的http响应指定状态码出现的频率，对请求进行拦截

## 使用方法
    b:=blocker.New(cache)
    //每分钟请求不能超过20次
    b.Block(blocker.StatusAny, 20, 1*time.Minute)
    //每分钟错误请求(status code >400)不能超过10次
    b.Block(blocker.StatusAnyError, 10, 1*time.Minute)
    //每分钟404错误不能超过5次
	b.Block(404, 5, 1*time.Minute)
    //每分钟403错误不能超过5次
	b.Block(403, 5, 1*time.Minute)

    App.Use(blocker)


## 设置
### 任意状态于任意错误状态
blocker.StatusAny (0) 代表任何状态码

blocker.StatusAnyError(-1) 代表任何大于等于400的状态码

### 自定义超过限制时的错误码

设置拦截器的 StatusCodeBlocked属性可以指定请求被拦截后的状态码

默认值为429

    b:=blocker.New(cache)
    b.StatusCodeBlocked=400

### 自定义请求的标识函数

设置拦截器的 Identifier方法可以设置用什么函数去标识请求.
默认值为,通过http请求的remoteAddr中的ip地址来进行标识

    b:=blocker.New(cache)
    b.Identifier=func(r *http.Request) (string, error) {
    	return r.Header.Get("name"), nil
    }

