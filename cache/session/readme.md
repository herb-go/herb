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
    DriverName=""
	Marshaler="msgpack"

    //Token色纸
	TokenLifetimeInHour=24
	TokenLifetimeInDay =7
	TokenMaxLifetimeInDay=24
	TokenContextName="token"
	TokenPrefixMode=""
	TokenLength=64

    #Cooke设置
	CookieName="herb-session"
	CookiePath="/"
	CookieSecure=false
	AutoGenerate=false
	UpdateActiveIntervalInSecond=60
	DefaultSessionFlag=0


    #客户端会话设置
    ClientStoreKey="key"

    #缓存会话设置
	"Cache.Driver"="syncmapcache"
    "Cache.TTL"=1800
    "Cache.Config.Size"=5000000
