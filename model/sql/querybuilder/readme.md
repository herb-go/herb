# querybuilder SQL 语句拼接器

提供拼接 sql 语句的接口，支持设置 sql 类型生成不同的语句

## 使用方法

使用 QueryBuilder 需要创建 builder 或者使用实现 Table 接口的对象
注意，第一次使用 builder 创建语句后，就无法再修改驱动类型。

### 直接创建

    //创建builder
    builder:=querybuilder.New()
    builder.Driver="mysql"

### 使用 Table 对象

Table 对象是使用实现了 DBTable 接口的对象操作实现的操作类。
github.com/herb-go/herb/model/sql/db 的 Table 对象是一个标准的可用于创建 Table 的对象

    //使用 Table 对象
    table:=NewTable(dbtable)
    builder:=table.QueryBuilder()

### 执行语句

    builder创建的语句可以通过Exec ,QueryRow ,QueryRows 三个方法及其衍生方式执行。

    执行时需要传入一个实现了DB接口的对象

    github.com/herb-go/herb/model/sql/db 的DB,Table,TxDB都是有效的DB对象
    Table对象也通过继承实现了DB对象的接口

    //通过builder来执行
    sqlResult,err:=builder.Exec(db,query)
    sqlRow:=builder.QueryRow(db,query)
    sqlRows,err:=builder.QueryRows(db,query)

    //通过查询语句来执行
    sqlResult,err:=query.Exec(db)
    sqlRow:=query.QueryRow(db)
    sqlRows,err:=query.QueryRows(db)

## 支持的语句

### Field 数据字段

Field 与 Fields 对象是为了方便进行数据库数据与程序数据结构进行映射儿创建的对象

### SELECT 语句

用户查询语句。使用方式:

    selectquery:=builder.NewSelect()
    selectquery.Select.AddFields(fields)
    selectquery.OrderBy.Add("id", false)
    selectquery.Limit.SetLimit(1)
    selectquery.Limit.SetOffset(1)
    selectquery.Join.LeftJoin().On(builder.New("field1=field2")).Alias("t2", "table2")
    row := selectquery.QueryRow(table1)
    err := selectquery.Result().BindFields(fields).ScanFrom(row)

#### SELECT 子语句

#### FROM 子语句

#### JOIN 子语句

#### WHERE 子语句

#### GROUPBY 子语句

#### ORDERBY 子语句

#### OTHER 额外部分子语句

### INSERT 语句

### UPDATE 语句

### DELETE 语句

### MISC 杂项

## 支持的数据库及类型

### mysql

驱动:github.com/go-sql-driver/mysql

### sqlite

驱动:github.com/mattn/go-sqlite3

### posgresql

驱动:github.com/lib/pq

### mssql

驱动:github.com/denisenkom/go-mssqldb
