# redispool Redis pool 模块

便于序列化为配置文件的redis pooal结构

基于 github.com/garyburd/redigo/redis 库

##  配置文件

    #TOML版本，其他版本可以根据对应格式配置
    #网络链接方式，一般为tcp
	Network="tcp"
	#redis服务器地址
	Address="127.0.0.1:6379"
    #redis密码
	Password="123456"
	#按秒计算的连接超时属性。默认值为 defaultIdleTimeout(60)
	ConnectTimeoutInSecond=30
	#按秒计算的读取数据的超时
	ReadTimeoutInSecond=30
	#按秒计算的写入数据的超时
	WriteTimeoutInSecond=30
	#redis的数据库编号
	Db=1
	#redis连接池的最大空闲连接数，默认值为defaultMaxIdle(200)
	MaxIdle=100
	#redis连接池的最大有效连接数，默认值为defaultMaxAlive(200)
	MaxAlive=100
	#按秒计算的空闲连接超时上线，默认值为defaultIdleTimeout(60)
	IdleTimeoutInSecond=30

## 使用方式

    pool:=redispool.New()
    config:=redispoll.NewConfig()
    err:=toml.Unmarshal(data,config)
    pool.Open()
    <- quitchan
    //退出时需要关闭连接池
    pool.Close()

    func useredisconn(){
        conn:=poll.Get()
        //使用链接后需要关闭
        defer conn.Close()
        //使用连接接，参考 github.com/garyburd/redigo/redis 的Conn 对象
        conn.Do(command)
    }