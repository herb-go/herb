# Events 简单事件模型

提供异步触发和处理事件的接口，用于代码解耦

##使用方法

1.定义事件类型

    EventTypeUserLogin=events.Type("userlogin")

2.创建处理程序，监听事件

    EventHandler=func(e *eventsEvent){
       group,ok:=e.Data.(string)
       if ok==false{
           group="unknown group"
       }
       logger.Log(" user "+e.Target+"@" +group+" login ")
    }
    //监听默认事件系统
    eventes.On(EventTypeUserLogin,EventHandler)

3.在必要的地方触发事件

    event:=EventTypeUserLogin.NewEvent().
        WithTarget(uid).
        WithData(usergroup)
    //Emit方法是异步调用的
     //返回值为是否注册过相应的事件
    handled:=eventes.Emit(event)

4.快速包装监听和触发事件的帮助函数 

    EmitUserLogin:=WrapEmit(EventTypeUserLogin)
    OnUserLogn:=WrapOn(EventTypeUserLogin)

    //监听事件
    OnUserlogin(EventHandler)

    //触发事件,参数为空则创建新事件
    EmitUserLogin(nil)
