# Formdata 用户提交数据接口
处理用户提交数据的接口，以及基于JOSN格式的实现

## 使用方式

    //创建自定义表单

    type ExampleForm struct{
        form.Form
        Field1 string
        Field2 *string
        //内部变量，不会被序列化，用于在InitWithRequest方法中储存当前用户信息
        uid string
    }

    //设置表单id，便于区分表单
    constExampleFormID = "exampleform"

    //表单初始化函数
    func NewExampleForm() *ExampleForm{
	    form:=&ExampleForm{}
	    form.SetModelID(ExampleFormID)
	    form.SetFieldLabels(ExampleFormFieldLabels)
	    return form
    }

    //表单验证函数
    func (f *ExampleForm) Validate() error {
        f.ValidateFieldf(f.Field1 != "", "Field1", "{{label}}必填") 
        f.ValidateFieldf(f.Field2 != nil, "Field2", "{{label}}必填") 
        if !f.HasError() {
            //此处添加需要没有其他错误才执行的代码
        }
        return nil
    }

    //表单实例通过http request初始化函数
    //一般用于获取请求对应的用户等信息，便于表单验证
    func (f *ExampleForm) InitWithRequest(r *http.Request) error {
        f.uid=MustGetUIDFromRequest(r)
    	return nil
    }    

    //标准的表单验证动作
    var formAction=func(w http.ResponseWriter, r *http.Request)
        //创建表单
        form := forms.NewExampleFormForm()
        //验证表单
        if formdata.MustValidateJSONRequest(r, form) {
            //替换为成功的业务逻辑
            render.MustJSON(w, form, 200)
        } else {
            //将表单错误渲染为状态码为422的JSON输出
            formdata.MustRenderErrorsJSON(w, form)
        }
