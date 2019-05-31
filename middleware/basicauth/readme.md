# Basicauth basic auth 认证组件
提供通过http basic auth 认证请求功能的组件


## 配置说明

单用户配置

    #TOML版本，其他版本可以根据对应格式配置
    #HTTP Realm值
    Realm="auth"
    #用户名
	Username="user"
    #密码
	Password="pass"

多用户配置

    #TOML版本，其他版本可以根据对应格式配置
    #HTTP Realm值
    Realm="auth"
    #用户部分
    [Users]
    #用户帐号密码，格式为用户名=密码
    "user"="pass"

## 使用说明

使用单用户中间件

    ba=&basiauth.SingleUser{}
    err=toml.Unmarshal(data,ba)
    app.Use(basiauth.Middleware(ba))

使用多用户中间件

    ba=&basiauth.Users{}
    err=toml.Unmarshal(data,ba)
    app.Use(basiauth.Middleware(ba))
