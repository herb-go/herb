# Cache  缓存组件

将数据存储于各种缓存驱动中的组件

## 关于缓存的约定

* 过期的数据一定不可被访问
* 未过期的数据不保证在所有情况下的可用性
* 支持永不过期的数据，但不代表数据永不丢失，也不建议使用永不过期的数据
* 所有的数据通过[]byte保存，通用可序列化结构的数据可以直接序列化保存
* 所有缓存取出的原始[]byte数据不应该被修改。需要修改请自行复制或者反序列化
* 缓存的Flush方法和用不过期数据是否支持取决与具体驱动的实现

## 特性

* 通用的缓存驱动接口，便于引入更多驱动
* 支持并发
* 引入Load方法，一定程度避免缓存雪崩
* 提供cacheable接口，并提供collection和node两个额外实现，方便缓存的复用
* 支持自定义数据的序列化方式

## 可用驱动

* dummpycache:空缓存。不缓存任何数据
* [syncmapcache](drivers/syncmapcache): 基于sync.Map的本地缓存驱动
* [freecache](drivers/freecache): 基于 github.com/coocood/freecache 的本地缓存驱动
* [gocache](drivers/gocache) 基于 github.com/patrickmn/go-cache 实现的本地缓存驱动
* [sqlcache](https://github.com/herb-go/providers/tree/master/sql/sqlcache) 基于sql的缓存
* [rediscache](https://github.com/herb-go/providers/tree/master/redis/rediscache) 基于redis的缓存，需要独占db
* [redisluacache](https://github.com/herb-go/providers/tree/master/redis/redisluacache) 基于redis和lua的缓存。多个缓存可以共用db
* [versioncache](drivers/versioncache) 利用本地、远程缓存和版本控制，兼顾访问效率和多机可用的缓存接口
  
## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    #驱动名，具体值取决于需要使用的驱动
    Driver="syncmapcache"
    #缓存默认有效时间，单位为秒。
    TTL=60   
    #Config部分为具体驱动设置，参考各个驱动的文档
    [Config]
    Size=50000000
    
## 使用缓存 

### 创建缓存对象

    c:=cache.New()
    config:=&cache.OptionConfig{}
    err:=config.ApplyTo(c)

### 操作二进制数据([]byte)

     //根据主键获取数据
    bs,err:=c.GetBytesValue("name")

    //根据数据设置数据.第三个参数为缓存最大有效时间
    err=c.SetBytesValue("name",bs,60*time.Second)

    //根据数据更新数据.只有本身缓存中就有数据才会更新。不然不储存数据。也不会返回错误。
    err=c.UpdateBytesValue("name",bs,60*time.Second)

    //删除缓存
    err=c.Del("name")

    //重设缓存过期时间
    err=c.Expire("name",60*time.Second)

    //批量获取缓存数据。返回的结果为map[string][]byte形式
	data,err=c.MGetBytesValue(keys ...string) (map[string][]byte, error)

    //批量设置map[string][]byte形式的数据
	err=c.MSetBytesValue(data,60*time.Second) 

### 使用预设的序列化器直接存取结构
    //根据主键获取缓存值.必须传入指针
    var v string
    err=c.Get("name",&v)

    //根据主键设置缓存
    err=c.Set("name","value",60*time.Second)

    //根据主键更新缓存
    err=c.Set("name","value",60*time.Second)

### 通过Load方法和loader函数加载数据

### 使用计数器

同名的计数器和二进制/结构数据是独立额，互相不影响

    //计数器递增值
    value,err=c.IncrCounter("name", 1, 60 * time.Second)

    //设置计数器的值
	err=SetCounter("name", 10, 60 * time.Second)

	//获取计数器的值
	v,err=GetCounter("name")

	//删除计数器的值
	err=DelCounter("name")

    //刷新计数器的过期时间
	err=ExpireCounter("name", 10*time.Second) error

### 其他杂项操作

    //清除所有数据。不是所有驱动都能支持
    err=c.Flush()

    //关闭缓存
    err=c.Close()

    //获取绝对主键
    key,err=c.FinalKey(string)

   //序列化数据
   bs,err=c.Marshal(v)

   //反序列化数据
   var v string
   err=c.Unmarshal(bs,&v)

## 缓存操作错误
* ErrNotFound: 指定主键的数据未找到
* ErrNotCacheable :数据无法储存
* ErrEntryTooLarge:数据过大
* ErrKeyTooLarge:主键过长
* ErrKeyUnavailable:主键无效
* ErrFeatureNotSupported:驱动不支持该功能
* ErrPermanentCacheNotSupport:不支持永不过期的缓存


## 缓存复用

缓存复用指将创建好的缓存划分成多个cacheable的组件，便于在不同的莫快中进行使用。目前支持的复用组件为Collection和Node

### Collection
Collection可以通过cacheable.Collection(Name)的方式创建。
Collection支持flush数据，不支持永久储存。通过利用两个储存当前实际主键和实际数据的字段来实现。对于访问速度和内存占用有较大影响。请仅在需要对一系列数据进行flush操作时使用。

### Node
Node可以通过cacheable.Node(Name)的方式创建。
Node支持永久储存，不支持flush数据。通过给主键加上固定的前缀实现，对于访问速度和内存占用影响较小，推荐一般情况下使用。