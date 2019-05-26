# session 会话组件

用于网页会话数据储存的组件

## 功能
* 提供基于缓存或者客户端储存(cookie)机制的缓存驱动
* 提供访问刷新会话生命周期的功能
* 提供Cooke以及Header方式传递会话Token的功能
* 提供符合 github.com/herb/user 接口的登录功能
* 便于序列化的配置方式

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    
    #基本设置
	#驱动名，可选值如下
	#DriverName="cookie"  基于的客户端会话
	#DriverName="cache"  基于的缓存的服务器端会话
    DriverName="cache"
	#缓存内部使用的序列化器。默认值为msgpack,需要先行引入
	Marshaler="msgpack"

    #Token设置
	#基于小时的会话令牌有效时间
	TokenLifetimeInHour=0
	#基于天的会话令牌有效时间
	#当基于小时的挥发令牌有效事件非0时，本选项无效
	TokenLifetimeInDay =7
    #基于天的令牌最大有效时间
	#当UpdateActiveIntervalInSecond值大于0时，令牌在访问后会更新有效时间
	#这个值决定了有效事件的最大值
	TokenMaxLifetimeInDay=24
	#访问时更新token的间隔。默认值为60
	UpdateActiveIntervalInSecond=60
	#会话令牌在HTTP请求上下文中的名字。当同时使用多套上下文时需要指定不同的名字。默认值为"token"
	TokenContextName="token"
	#是否自动生成会话，默认值为false
	AutoGenerate=false


    #Cooke设置
	#储存session的cookie名，默认值为"herb-session"
	CookieName="herb-session"
	#Cookie的路径设置，默认值为"/"
	CookiePath="/"
	#Cookie的Secure,默认值为false
	CookieSecure=false

	#其他设置
	#默认会话的标志位信息.默认值为1
	DefaultSessionFlag=1


    #客户端会话设置
	#客户端会话密钥
    ClientStoreKey="key"

    #缓存会话设置
	#会话令牌前缀模式。可用值为
	#"empty":空
	#"raw":原始值
	#"md5":md5后的摘要值
	#默认值为raw
	TokenPrefixMode=""
	#令牌数据长度。
	#注意数据长度是原始数据长度。存入cookie时的长度还要经过base64转换
	#默认值为64
	TokenLength=64
	#缓存驱动
	"Cache.Driver"="syncmapcache"
	#缓存有限时间
    "Cache.TTL"=1800
	#具体缓存配置
    "Cache.Config.Size"=5000000
