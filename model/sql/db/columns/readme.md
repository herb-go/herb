# Columns 数据库列信息组件
根据数据库配置获取数据库列信息的组件

## 数据库列结构

Column是数据库的每一列的信息

    type Column struct {
        Field      string
        ColumnType string
        AutoValue  bool
        PrimayKey  bool
        NotNull    bool
    }

* Field 列名
* ColumnType 列对应的变量类型
* AutoValue 列是否是自动赋值的
* PrimayKey 是否是主键
* NotNull 列是否为非空

## 使用方式

    columns,err:=culumns.Driver("drivername")

## 可用驱动

* [Mysql](mysqlcolumns/readme.md)
* [Sqlite](sqlitecolumns/readme.md)