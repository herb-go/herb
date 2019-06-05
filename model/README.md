# Model 数据模型模块

用于处理和用户有关的数据模型

## Message 信息

Message 是包含了一套翻译后的数据的集合。

使用方式:
    m:=models.NewMessages()
    m.
        SetMessage("value1","translated1 %[1]s").
        SetMessage("value2","translated1 %[2]s").
    
    