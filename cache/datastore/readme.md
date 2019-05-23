# Datastore 数据存储组件

将数据依据主键的方式进行存储，并能快速的通过创建函数(creator)和加载函数(loader)从缓存和数据源中进行批量加载。

一般用于作为用户或者其他数据模型的缓存

## 功能
* 提供store接口，便于根据实际数据类型自定义store实现，快速转换数据
* 使用提供的缓存进行数据的序列化与反序列化
* 预防缓存雪崩

## 缺点

* 只支持字符串类型的主键

## 使用方法

### Store 数据存储集
储存所有已经从数据源和缓存中取出的数据的接口类型。

需要实现5个方法

    Store.Load(key string) (value interface{}, ok bool)

从数据存储集中按给定主键取出对象值。

返回值value为具体的值

返回值ok为是否成功获取。

    Store.Store(key string, value interface{})

将给到的数据(value)以指定主键(key)存入数据存储集

	Store.Delete(key string)

在数据存储集中按主键删除数据。只有通过这个方法才能删除特定的数据，通过Store方法存入nil并不视为删除数据

	Store.Flush()

清除所有数据

	Store.LoadInterface(key string) interface{}

将数据强制以 interface{} 的形式取出。与Load方法的区别是并不关注数据是否存在。不存在的数据以nil返回

### MapStore与SyncMapStore

 MapStore和SyncMapStore的区别是具体实现的驱动

 MapStore的实现是map[string]interface{},效率理论较高，非线程安全。创建方式为

     s:=datastore.NewMapStore()

 SyncMapStore的实现是sync.Map，线程安全。创建方式为

     s:=datastore.NewSyncMapStore()

### Creator/Loader方法
这两个是负责数据初始化和加载的函数

Creator的形式为

    func() interface{}

返回一个初始化好的数据

Loader的形式为

    func(...string) (map[string]interface{}, error)

通过给定的主键列表加载返回加载成功的数据集，或者任何发生的错误。

###直接通过Store/Creator/Loader 加载数据

    err:=Load(store, datacache, loader, creator, key1,key2,key3,key4....)  

将主键列表中的数据加载到store中。同时在给定缓存中对未缓存的数据进行缓存

### 通过Datasource加载数据
Datasource是一个包含了数据Creator和Loader的结构。创建方式为

    datasource:=datastore.NewDataSource()
    datasource.Creator=creator
    datasource.SourceLoader=loader

通过datasouce加载数据有两种方式

1.传入缓存，直接加载

    err:=datasource.Load(store, c cache, key1,key2.key3...) error {
    
2.传入缓存，生成Loader,通过Loader加载

Loader是一个包含store,cache,datasource的结构

    loader:=datasource.NewSyncMapStoreLoader(cache)
    //如果需要使用MapStore的话,使用
    //loader:=datasource.NewMapStoreLoader(cache)
    err:=loader.Load(key1,key2,key3)
    store:=loader.Store

对loader可以进行Delete和Flush操作,会同时对对应的Store和缓存进行操作

    err:=loader.Delete(key)
    //同时删除对应的数据存储集和缓存里的对应主键的值
    err:=loader.Flust()
    //同时清空对应的数据存储集和缓存

### 通过实现BatchDataLoader 接口的结构加载
BatchDataLoader  是一个包含数据Creator和Loader的接口。Datasource就是一个BatchDataLoader 的实现。

通过BatchDataLoader加载数据有两种方式

1.传入缓存，直接加载

    err:=LoadWithBatchLoader(store, cache, batchdataloader, key1,key2,key3...)

2.传入缓存，生成Loader，通过Loader加载

    loader:=NewSyncMapLoader(cache,batchdataloader)
    //如果需要使用MapStore的话，使用
    //loader:=NewMapLoader(cache,batchdataloader)
    err:=loader.Load(key1,key2,key3)
    store:=loader.Store
