# Store 文件储存接口

定义了用于网站上传下载的文件传输接口，以及本地仓库的实现

## 可用驱动
* assest 本地文件仓库

## 配置说明

    #TOML版本，其他版本可以根据对应格式配置
    #驱动名，具体值取决于需要使用的驱动
    Driver="assest"
    #驱动配置部分
    [Config]

## 本地文件仓库驱动配置

    #TOML版本，其他版本可以根据对应格式配置
    #驱动名
    Driver="assest"
    #驱动配置部分
    [Config
    #前台路径访问域名部分
    URLHost="http://www.test.com"
    #前台路径目录部分
	URLPrefix="/upload"
    #本地文件根目录
	Root="/tmp"
    #路径为绝对路径还是相对路径
	Absolute=true
    #本地文件子目录
	Location="/savedfiles"

## 使用方法

### 创建Store

    s:=store.New()
    config=&store.ConfigMap{}
    err=toml.Unmarshal(data,config)
    err=config.ApplyTo(store)

### 上传和载入文件

    //上传，传入文件名和reader
    //reader会被驱动关闭
    //返回文件id,长度和错误
    
    id, length, err := s.Save("filename",reader)
    
    //下载,传入文件id,返回reader和错误
    //用户需要自行关闭reader
    reader, err := s.Load(id)
    defer reader.Close()

    //删除文件
	err = s.Remove(id)
    //通过id换取文件url
    url, err = s.URL(f.ID)

### 创建和使用File对象
    
    //根据id创建file对象
    f:=s.File(id)
    
    //获取file对象的url
	url, err =f.URL()

### 错误列表

文件错误(*store.ErrorType)是一类文件相关的错误

可以通过把错误换转为文件错误来获取响应的文件信息

    if storeerr=err.(*store.ErrorType){
        fmt.Println("文件名",storeerr.File)
        fmt.Println("错误类型",storeerr.Type)
        
    }

文件错误类型列表为

* store.ErrorTypeNotExists 文件不存在
* store.ErrorTypeUnavailableID 文件ID无效