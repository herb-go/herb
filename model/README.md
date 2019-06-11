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

## DefaultMessagesChain 默认信息链

DefaultMessagesChain 是一个空的 MessagesChain。是默认情况下的翻译来源。

    //直接链入Messages
    model.Use(messages1,message2)

    //直接获取翻译
    label:=model.GetMessage("label1")

## Model 模型对象

模型对象是一个用于集成的用户数据验证基础结构

提供了初始化，绑定 Http request,验证并存储用户输出数据错误的接口

    //继承对象
    type Form stuct{
        model.Model
        httprequest *http.Request
    }

    //实现InitWithRequest方法
    func (model *Form) InitWithRequest(r *http.Request) error {
        model.httprequest=r
        return nil
    }

    //设置字段名称，可以为Model设置错误提示时每个字段的名称。需要每次在实例化的时候添加
    //没有设置过的字段名称会以属性名表示
    var FormFieldLabels=map[string]strint{
        "field1":"Field 1",
        "field2:"Field 2"
    }

    //表单对应的ID
    const FormID="formid"

    //创建新Form对象
    func New() *Form{
        form:=&Form{}
        form.SetFieldLabels(FormFieldLabels)

        //设置表单使用的翻译集，不设置或者为空的话使用model.DefaultMessageChain
        form.SetMessages(nil)

        //设置表单ID,便于在需要的时候快速处理和创建表单
        form.SetModelID(FormID)
    }

使用 Model 对象

    //创建新表单
    form:=NewForm()

    //验证字段。第一个参数为false的话，会为model添加名称为第二参数，值为第三参数的错误
    form.ValidateField(form.field1==1,"field1","field1 must be 1")
    //验证字段。第一个参数为false的话，会为model添加名称为第二参数，值为第三参数的错误
    //传入的字符串有两个特殊的占位符
    //%[1]s代表原始的Field名
    //%[2]s代表通过SetFieldLabels设置的字段名
    form.ValidateFieldf(form.field1==1,"field1","%[2]s must be 1")

    //添加不转换的错误信息
    form.AddPlainError("field", "error msg")
    //添加只转换字段信息的错误信息
    form.AddError("field", "error msg")
    //添加转换字段信息和错误信息
    //传入的字符串有两个特殊的占位符
    //%[1]s代表原始的Field名
    //%[2]s代表通过SetFieldLabels设置的字段名
    form.AddErrorf("field", " %[1]s error msg")

    //判断model对象是否有错误
    TrueOrFalse:=form.HasError()
    //返回model的所有错误
    errors:=form.Errors()
