# CSRF 预防跨站请求伪造组件
一个简单的基于Cookie的预防跨站请求伪造组件

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置

    #回传csrf token的Cookie名。默认值为herb-csrf-token
    CookieName="herb-csrf-token".
    #回传csrf token的Cookie名路径。默认值为/
	CookiePath="/"
    #接受并验证csrf token的头字段。默认为"X-CSRF-TOKEN"
	HeaderName="X-CSRF-TOKEN"
    #接受并验证csrf token的表单字段。默认为"X-CSRF-TOKEN"
	FormField="X-CSRF-TOKEN"
    #失败时返回的响应状态码。默认值为400.
	FailStatus=400
    #中间件是否起效。默认值为false
	Enabled=true
    #验证失败时添加的响应头
	FailHeader="csrffail"
    #验证失败时添加的响应头的值
	FailValue="true"
