# roleservice 权限服务

## Service 角色服务

角色服务包含了用户识别器和基于用户的角色获取器，能够直接从 http 请求中获取当前用户的所有角色

    //创建新的角色服务
    service:=role.NewService(RoleProvider,Identifier)

    //获取当前用户的所有权限
    roles,err=service.RolesFromRequest(r)

    //传入规则生成器，生成响应权限验证中间件
    //第二个参数为空则失败时返回403状态
    authorizerMiddleware=service.AuthorizeMiddleware(RuleProvider,nil)
    app.Use(authorizerMiddleware)

    //根据传入的角色名，生成权限验证中间件
    authorizerMiddleware=service.RolesAuthorizeOrForbiddenMiddleware(role1,role2)
    app.Use(authorizerMiddleware)
