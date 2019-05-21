# querybuilder SQL语句拼接器

提供拼接sql语句的接口，支持设置sql类型生成不同的语句

## 支持的语法 

### SELECT

用户查询语句。使用方式:

```
selectquery:=builder.NewSelect()
selectquery.Select.AddFields(fields)
selectquery.OrderBy.Add("id", false)
selectquery.Limit.SetLimit(1)
selectquery.Limit.SetOffset(1)
selectquery.Join.LeftJoin().On(builder.New("field1=field2")).Alias("t2", "table2")
row := selectquery.QueryRow(table1)
err := selectquery.Result().BindFields(fields).ScanFrom(row)

```
#### SELECT子语句

#### FROM子语句

#### JOIN子语句

#### WHERE 子语句

#### GROUPBY子语句

#### ORDERBY子语句

#### OTHER 部分


### INSERT

### UPDATE

### DELETE

### MISC

## 支持的数据库及类型
### mysql
驱动:github.com/go-sql-driver/mysql

### sqlite
驱动:github.com/mattn/go-sqlite3

### posgresql
驱动:github.com/lib/pq

### mssql
驱动:github.com/denisenkom/go-mssqldb
