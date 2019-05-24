# VersionCache 基于版本的远程缓存驱动
基于版本的本地、远程缓存驱动，通过将大体积数据储存在本地提升效率的缓存结构。一般在需要多进程/多机使用缓存时提升效率

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    #本地缓存配置
    "Local.Driver"="syncmapcache"
    "Local.TTL"="1800"
    "Local.Config.Size"=5000000
    #远程缓存配置
    "Remote.Driver"="freecache"
    "Remote.TTL"="1800"
    "Remote.Config.Size"=5000000
