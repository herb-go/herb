# cachegroup 缓存组缓存驱动

仅作理论验证实现，不推荐在生产环境使用

将一组缓存作为统一为一个缓存组进行使用，实现冷热分离


## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    #Caches []string数据，为配置中所有子缓存的主键列表，如
    Caches=["cache1","cache2","cache3"]
    #子缓存配置，以Cache中设置过的主键为前缀,如
    "cache1.Driver"="syncmapcache"
    "cache1.TTL"="1800"
    "cache1.Config.Size"=5000000
    "cache2.Driver"="freecache"
    "cache2.TTL"="1800"
    "cache2.Config.Size"=5000000
    "cache3.Driver"="gocache"
    "cache3.TTL"="1800"
    "cache3.Config.Size"=5000000