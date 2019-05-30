# Simple http server 基于http标准库封装的 静态文件服务

通过http标准库的ServeFile来提供文件和目录的访问服务

## 使用方法

### 服务单文件

    app.HandleFunc(simplehttpserver.ServeFile("filerealpath"))

### 服务目录

目录或子目录下如果有index.html.会当成目录正文，否则返回403

    app.StripPrefix("'public",simplehttpserver.ServeFolder("folderrealpath))