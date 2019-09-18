# User 用户模块
提供一系列与网站用户相关的接口与操作函数

## 接口
* 用户帐号接口 Account

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

