# db 数据库模块
便于序列化为配置的SQL数据库结构


##  配置文件

    #TOML版本，其他版本可以根据对应格式配置
    #数据库驱动名
    Driver="mysql"
	#数据库类型名。一般用于odbc等可以对应多套数据库的驱动。默认和Driver一致
	Type="mysql"
	#数据库链接字符串，具体值取决于驱动
	DataSource=""
	#数据库表前缀
	Prefix string
    #最大连接数量
    MaxOpenConns=60
	#最大空闲链接数量
	MaxIdleConns=30
    #按秒计算的连接最长生命周期。默认值为30秒
	ConnMaxLifetimeInSecond=60
	
## 使用方式

	//创建数据库对象
	db:=New()
    config:=NewConfig()
	err:=toml.Unmarshal(data,config)
	err=db.Init(config)

	//创建数据库表对象
	table:=db.Table("tablename")

	//创建数据库事务对象。使用后请自行提交或者回退。
	txdb:=NewTxDb(db)
	defer txdb.RollBack()
	err=txdb.Commit()

	//通过数据库对象操作数据库
	sqlresult,err=db.Exec(query, arg1,arg2) 
	sqlrows,err=db.Query(query, arg1,arg2) 
	sqlrow,err=db.QueryRow(query, arg1,arg2)

	//通过数据库表对象操作数据库
	sqlresult,err=table.Exec(query, arg1,arg2) 
	sqlrows,err=table.Query(query, arg1,arg2) 
	sqlrow,err=table.QueryRow(query, arg1,arg2)

	//通过事务对象操作数据库
	sqlresult,err=table.Exec(query, arg1,arg2) 
	sqlrows,err=table.Query(query, arg1,arg2) 
	sqlrow,err=table.QueryRow(query, arg1,arg2)
	err=txdb.Commit()