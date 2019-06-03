# Forwarded 转发请求原始信息接口
一个获取被转发请求的原始信息的接口

## 功能
* 获取请求原始IP地址
* 获取请求原始协议
* 获取请求原始域名信息
* 通过设置token验证请求是否有效

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    #是否启用
	Enabled=true

    #原始ip头，为空则该功能不起效
	ForwardedForHeader=""

    #原始Host头，为空则该功能不起效
	ForwardedHostHeader=""

    #原始协议头，为空则该功能不起效
	ForwardedProtoHeader=""

    #转发信息认证头，为空该功能不起效
	ForwardedTokenHeader="forwardedtoken"
	#转发信息认真值，转发信息认证头起效时，传入的请求的对应请求头的值必须和该值相等，不然是无效请求
	ForwardedTokenValue="1234567890"

    #失败状态码，默认为400
	FailStatusCode=400

    #degbu模式，启用后会将客户ip加载响应头的"X-Remote-Addr"字段内
	Debug=false

## 使用方式

    m:=&forwarded.Middleware{}
    err=toml.Unmarshal(data,m)
    app.Use(m)