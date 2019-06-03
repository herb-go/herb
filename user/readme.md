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