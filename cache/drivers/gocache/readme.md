# Gocache 缓存驱动

基于 github.com/patrickmn/go-cache 实现的本地缓存驱动

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    "Driver"="gocache"
    "TTL"="1800"
    //默认过期时间，单位为秒，默认值60
    "DefaultExpirationInSecond"=1800
    //清理过期数据间隔，单位为秒，默认值60
	"CleanupIntervalInSecond"=1800