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
	UpdateActiveIntervalInSecond=60
	TokenContextName="token"
	TokenPrefixMode=""
	TokenLength=64

    #Cooke设置
	CookieName="herb-session"
	CookiePath="/"
	CookieSecure=false
	AutoGenerate=false

	DefaultSessionFlag=0


    #客户端会话设置
    ClientStoreKey="key"

    #缓存会话设置
	"Cache.Driver"="syncmapcache"
    "Cache.TTL"=1800
    "Cache.Config.Size"=5000000
