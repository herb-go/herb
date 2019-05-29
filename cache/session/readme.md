# Session 会话组件

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
	#当使用自动模式安装中间件时使用的模式
	#"header"使用header模式
	#"cookie"或其他值使用cookie模式
	Mode="cookie"
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

## 使用方法

### 创建会话，进行配置

    store:=session.New()
	config:=&session.StoreConfig{}
	err=toml.Unmarshal(data,&config)
	err=config.ApplyTo(store)

### 安装会话中间件

安装会话中间件的方式包括cookie模式,header模式，自动模式

1.cookie模式安装，自动通过配置中CookieName的cookie值作为session的token

    app.Use(store.CookieMiddleware)

2.header模式安装，通过指定的请求头的值做为session token。

客户端需要自行维护token

    app.Use(store.HeaderMiddleware("headername"))

3.自动模式安装。

通过配置文件中的Mode值决定使用cookie模式还是session模式安装

如果Mode值为header,使用配置中的CookieName为请求头作为session token值，由客户端自行维护token

其他情况下同cookie模式安装

    app.Use(store.InstallMiddleware())

4.使用注销中间件

将对应请的session清除，一般用于注销

    app.Use(store.DestoryMiddleware()).HandleFunc(logoutaction)

### 在动作中设置与获取Session值

使用Store.Get，Store.Set和Store.Del维护session

注意，正常情况下session值会在程序正常运行，返回至会话中间件时才进行更新和cookie变更。之间如果程序panic的话，之前的设置会失效

    func(w http.ResponseWriter, r *http.Request) {
		err=store.Set(r,"sessionfieldname","new value")
		var v string
		//Get时需要传入指针
		err=store.Get(r,"sessionfieldname",&v)
		err=store.Del(r,"sessionfieldname")
	}

### Session对象

Session是一个存放了所有会话数据的可序列化的结构。

#### 获取Session

获取和维护Session主要有两个方向

1.通过Session Store从http请求中获取

    //获取session
	s,err:=store.GetRequestSession(r)

	//将请求中的Session进行保存
	err=s.SaveRequestSession(r)

2.通过Session Store直接创建/维护Session

    s:=session.NewSession(store,"token")

#### 使用Session对象

	err=s.Set("sessionfieldname","value)
	var v string
    err=s.Get("sessionfieldname",&v)
	err=s.Del("sessionfieldname")
	//获取session的token值

注意，session值的修改需要保存后才能存入store

#### 维护Session对象
 
     //重置session除token外的所有数据
	 s.Regenerate()
	 //从Store中通过token加载Session
	 err=s.Load()
	 //保存Session
	 err=s.Save()
	 //删除SEssion
	 err=s.DeleteAndSave()

#### Sesssion的token操作

    //获取token
	token,err:=s.Token()
	token=s.MustToken()
	//设置token
	s.SetToken(token)
	//重新生成token,需要设置Prefix，比如用户Id
	err=s.RegenerateToken("prefix)

### Field对象

Field对象是指指定字段名的session字段。

他的用途主要有两个，一个是为使用session的程序提供一个储存特定数据的接口，另一个是可以直接用来的进行用户的登录和登出

    //创建field
	field:=store.Field("fieldname)
	//使用field数据
	err=field.Set("value")
	var v string
	err=field.Get(&v)
	//删除数据
	err=field.Flush()

	//实现用户接口
	uid,err:=field.IdentifyRequest(r)
	//实现用户登录接口
	err=field.Login(w,r,"userid)
	//实现用户登出接口
	err=field.Logout(w,r)