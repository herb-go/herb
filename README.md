# Herb 网络工具库
一套完整的网络开发的解决方案库。

状态：稳定

## 特点
* 基于golang标准库的http库
* 定义基本功能并采用自带/第三方驱动的方式实现，注重标准化，可替换性
* 大模块直接相互独立。可以分别独自使用

## 缺点
* 不将性能优化作为第一重点
## 功能

###  缓存组件
提供一系列的缓存驱动/缓存/缓存组件的接口，确保程序的高效有序执行
>  [Cache 缓存组件](cache/readme.md)
>  > [Blocker 拦截器组件](cache/blocker/readme.md)
> >
>  > [DataStore 数据存储组件](cache/datastore/readme.md)
> >
> >Drivers 驱动
>  >
> > >  [Freecache 驱动](cache/drivers/freecache/readme.md)
> > > 
> > > [Gocache 驱动](cache/drivers/gocache/readme.md)
> > >
> > >[Syncmap 驱动](cache/drivers/syncmapcache/readme.md)
> > >
> > > [VersionCache 驱动](cache/drivers/versioncache/readme.md)
> >
> >Marshalers 序列化器
> > > [Msgpack 序列化器](cache/marshalers/msgpackmarshaler/readme.md)
> >
> > [Pagecache 页面缓存](cache/pagecache/readme.md)
> >
> > [Session 会话组件](cache/session/readme.md)
> >
> >> [Captcha 验证码组件](cache/session/captcha/readme.md)
### 事件支持
### 文件组件
### 中间件工具
### 数据模型工具
### 渲染接口
### 用户权限接口




## 约定
* 中间件:采用 func(w http.Writer,r *http.Request,Next http.Hanlderfunc) 作为中间件的形式
* 通用可序列化结构：以Golang默认的可序列化结构(首字母大写,仅支持string主键的Map，不带注解)作为通用的数据传输形式，通过JSON/TOML/MSGPACK协议进行压缩、配置、缓存
  
