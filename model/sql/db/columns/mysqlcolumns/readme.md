# Mysql columns Mysql列信息驱动

从mysql连接中获取mysql表列信息的驱动

## 使用方式
    import _ github.com/herb/model/sql/db/columns
    import _ github.com/herb/model/sql/db/columns/mysqlcolumns
    
    columns,err:=columns.Driver("mysql")

### 支持的字段类型

* TINYINT  =>byte
* BIT =>byte
* BOOL =>byte
* SMALLINT =>int
* MEDIUMINT =>int
* INT  =>int
* INTEGER =>int
* BIGINT =>int64
* FLOAT =>float32
* DOUBLE =>float64
* DOUBLE PRECISION  =>float64
* DATETIME =>time.Time
* TIMESTAMP =>time.Time
* CHAR =>string
* VARCHAR =>string
* TINYTEXT =>string
* TEXT =>string
* MEDIUMTEXT =>string
* LONGTEXT =>string
* BINARY =>[]byte
* VARBINARY =>[]byte
* TINYBLOB =>[]byte
* BLOB =>[]byte
* MEDIUMBLOB =>[]byte
* LONGBLOB =>[]byte
