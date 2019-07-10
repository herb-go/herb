# Sqlite columns sqlite列信息驱动

从mysql连接中获取mysql表列信息的驱动

## 使用方式
    import _ github.com/herb/model/sql/db/columns
    import _ github.com/herb/model/sql/db/columns/sqlitecolumns
    
    columns,err:=columns.Driver("sqlite3")

### 支持的字段类型

* BOOL =>byte
* SMALLINT =>int
* MEDIUMINT =>int
* INT =>int
* INTEGER  =>int
* TINYINT  =>int
* INT2 =>int
* INT8 =>int
* BIGINT =>int64
* FLOAT =>float32
* DOUBLE =>float64
* DOUBLE PRECISION  =>float64
* REAL =>float64
* DATETIME =>time.TIme
* CHAR =>string
* VARCHAR =>string
* CHARACTER =>string
* NVARCHAR =>string
* NCHAR =>string
* TEXT =>string
* BLOB =>byte

