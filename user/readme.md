# User 用户模块
提供一系列与网站用户相关的接口与操作函数

## 接口
* http请求身份识别接口 Identifier
* 登录与登出接口 LoginProvider/LogoutProvider
* 用户帐号接口 Account
* 用户鉴权接口 Authorizer

## 使用方式

### 帐号管理

    //Account对象
    account:=user.NewAccount()
    account.Keyword="keyword"
    account.Account="account"

    TrueOrFalse:=account.Equal(account2)

    //Accounts对象
    accounts:=user.NewAccounts()
    //绑定帐号。如果帐号已经存在，会返回错误 user.ErrAccountBindingExists
    acccounts.Bind(account1)
    //解绑帐号。如果帐号不存在，会返回错误 user.ErrAccountUnbindingNotExists
    acccounts.Unbind(account1)
    //判断帐号是否存在
    TrueOrFalse:=accounts.Exist(account1)
    
    //区分大小写的帐号创建器，返回的帐号为"AaBbCc"
    account=CaseSensitiveAcountProvider.NewAccount("keyword","AaBbCc")

    //不区分大小写的帐号创建器，返回的帐号为"aabbcc"
    account=CaseInsensitiveAcountProvider.NewAccount("keyword","AaBbCc")    

### 授权管理

  授权管理定义了授权管理器接口Authorizer，以及对应使用的中间件

  [role](role)实现了一个根据用户角色进行权限检测的Authorizer

    //对用户进行授权的检测.第二个参数为授权失败后的动作，如传入nil则失败返回403状态
     app.Use(user.AuthorizeMiddleware(Authorizer,nil))

    //授权检测，失败则返回403状态
     app.Use(user.AuthorizeOrForbiddenMiddleware(Authorizer))

### 用户识别

用户识别模块定义了基于http.Request的用户识别器 Identifier,以及响应的操作工具

Identifier 用户识别接口，识别给到的http请求对应的用户id

    //登录限制中间件，第二个参数为没有用户信息时执行的动作，为nil的话返回401错误
    loginrequired=user.LoginRequiredMiddleware(identifier,nil)
    app.Use(loginrequired)

    //登录跳转器，用于在用户需要登录时进行跳转，并记录原始信息
    //第一个参数为登录Url，第二个参数为储存原始链接的cookie名
    redirector=user.NewLoginRedirector("/login","login-returnurl")
    //可以通过设置跳转器的Cookie属性进行cookie的调整
    redirector.Cookie.Path="/app"

    app.Use(redirector.Middleware(identifier))

    //登录成功后获取原始信息的方式
    loginsuccess=func(w http.ResponseWriter, r *http.Request){
      //获取原始登录地址
      url:=redirector.MustClearSource(w,r)
      //判断是否为空
      if url==""{
        url="/home"
      }
      	http.Redirect(w, r, url, 302)
    }

### 登录/登出服务

用户模块定义了基本的用户登录，登出接口，及相关的快捷操作

* LoginProvider用户登录接口，将用户登录到指定的请求上
* LogoutProvider 用户登出接口，将用户从指定的请求上登出

登出中间件:

    app.
        Use(user.LogoutMiddleware(logoutProvider)).
        HandleFunc(func(w http.ResponseWriter, r *http.Request){
            w.Write([]byte("成功登出"))
        })

### 用户限制列表中间件

    //限制只有指定的用户才能访问。用户Id不在指定的列表里的话返回403错误
    app.Use(user.MiddlewareForbiddenExceptForUsers(Identifier,[]string["admin1","admin2"]))

### 重定向器

重定向器与登录跳转器功能类似。但可以用于更详细判断用户跳转的标准。
一般用于强制用户在访问的指定动作钱必须进行额外的验证或者资料补充。

    //跳转条件,返回是否需要
    condition=func(w http.ResponseWriter, req *http.Request) bool{
      return getUserConfrimed(w,req)
    }

    redirector:=user.NewRedirector("/bindemail","bindemail-returnurl",condition)

    app.Use(redirector)

    //操作成功后获取原始信息的方式
    bindemailSuccess=func(w http.ResponseWriter, r *http.Request){
      //获取原始地址
      url:=redirector.MustClearSource(w,r)
      //判断是否为空
      if url==""{
        url="/home"
      }
      http.Redirect(w, r, url, 302)
    }