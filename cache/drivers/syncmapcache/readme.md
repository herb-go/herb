# syncmap 缓存驱动

基于 sync库的Map 实现的缓存驱动

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    "Driver"="syncmapcache"
    "TTL"="1800"
    //空间占用大小，单位为byte，默认值50000000
    "Size"=50000000
    //清理过期数据间隔，单位为秒，默认值60
	"CleanupIntervalInSecond"=1800