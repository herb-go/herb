# Model 数据模型模块

用于处理和用户有关的数据模型

## Message 信息

Message 是包含了一套翻译后的数据的集合。

使用方式:

    m:=models.NewMessages()
    //设置信息
    m.
        SetMessage("value1","translated1 %[1]s").
        SetMessage("value2","translated1 %[2]s").

    //获取信息
    //如果没有对应的翻译，会将传入的字符串原样返回
    label:=m.GetMessage("value1")

    //获取信息和是否翻译
    //如果没有对应的翻译，会将传入的字符串原样返回，并返回false
    label,TrueOrFalse=m.LoadMessage("value1")

## MessagesChain 信息链

MessagesChain 是一个把多个 Messages 对象按顺序存储，获取翻译时依次获取值的的对象

    //创建MessagesChain
    m:=models.NewMessagesChain(messages1,messages2)

    //将更多messages加入MessagesChain
    m.Use(messages3,messages4).Use(message5)

    //获取信息
    //如果没有对应的翻译，会将传入的字符串原样返回
    //将按顺序从已经加入的Messages中查找
    label:=m.GetMessage("value1")

    //获取信息和是否翻译
    //如果没有对应的翻译，会将传入的字符串原样返回，并返回false
    //将按顺序从已经加入的Messages中查找
    label,TrueOrFalse=m.LoadMessage("value1")
