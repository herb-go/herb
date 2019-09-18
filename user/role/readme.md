# Role 用户角色模块
提供用户角色管理功能以实现 github.com/herb-go/user 的Authorizer接口

## Role 角色对象

角色对象是用户权限的最基础的元素，包括两部分内容，角色名和数据

    r:=role.New("rolename")
    //添加权限数据
    r.AddData("fieldname","data1","data2")

    //将role作为用户已经拥有的权限
    ownedRoles=[]Role{role1,role2}

    //将role作为待验证的权限规则
    requiredRole=roleRequired

    //判断用户是否拥有给定的权限。
    //当用户有匹配的权限(权限名一致，所有的roledata都覆盖)时，返回true，否则返回false
    TrueOrFalse,err=requiredRole.Execute(ownedRoles...)

## Roles 角色列表对象

角色列表一般代表用户拥有的所有角色信息
    //通过角色名列表创建角色列表对象。这种创建方式不能添加角色数据
    roles:=role.NewRoles("role1","role2","role3")

    //添加角色
    roles.Add(role4)

    //判断用户是否拥有给定的权限。
    //当roles里所有的规则都被用户角色匹配时，返回true,否则返回false
    TrueOrFalse,err=roles.Execute(ownedRoles...)

## Rule 用户规则及相关操作

用户规则是用于角色权限控制的核心接口，代表被访问的资源需要的权限。

    //多个规则进行and操作
    rulesAnd=role.And(rule1,rule2,rule3)
    //多个规则进行or操作
    rulesOr=role.Or(rule2,rule2,rule3)
    //not操作
    ruleNot=role.Not(rule1)

    //使用RuleSet
    rs=role.NewRuleSet().
        And(rule1,rule2.rule3).
        Or(rule4,rule5).
        Not()

