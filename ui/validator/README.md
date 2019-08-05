# Model 数据模型模块

用于处理和用户有关的数据模型


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
        form.SetComponentI(ui.MapLabels(FormFieldLabels))
        //设置表单ID,便于在需要的时候快速处理和创建表单
        form.SetComponentID(FormID)
    }

使用 Model 对象

    //创建新表单
    form:=NewForm()

    //验证字段。第一个参数为false的话，会为model添加名称为第二参数，值为第三参数的错误
    form.ValidateField(form.field1==1,"field1","field1 must be 1")
    //验证字段。第一个参数为false的话，会为model添加名称为第二参数，值为第三参数的错误
    //传入的字符串有两个特殊的占位符
    //{{field}}代表原始的Field名
    //{{label}}代表通过SetFieldLabels设置的字段名
    form.ValidateFieldf(form.field1==1,"field1","{{label}} must be 1")

    //添加不转换的错误信息
    form.AddPlainError("field", "error msg")
    //添加只转换字段信息的错误信息
    form.AddError("field", "error msg")
    //添加转换字段信息和错误信息
    //传入的字符串有两个特殊的占位符
    //{{label}}代表原始的Field名
    //{{field}}代表通过SetFieldLabels设置的字段名
    form.AddErrorf("field", " {{label}} error msg")

    //判断model对象是否有错误
    TrueOrFalse:=form.HasError()
    //返回model的所有错误
    errors:=form.Errors()
