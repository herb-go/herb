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

## 缓存复用

缓存复用指将创建好的缓存划分成多个cacheable的组件，便于在不同的莫快中进行使用。目前支持的复用组件为Collection和Node

### Collection
Collection可以通过cacheable.Collection(Name)的方式创建。
Collection支持flush数据，不支持永久储存。通过利用两个储存当前实际主键和实际数据的字段来实现。对于访问速度和内存占用有较大影响。请仅在需要对一系列数据进行flush操作时使用。

### Node
Node可以通过cacheable.Node(Name)的方式创建。
Node支持永久储存，不支持flush数据。通过给主键加上固定的前缀实现，对于访问速度和内存占用影响较小，推荐一般情况下使用。