# Router 路由接口
网页url地址解析的路由接口

## 可用路由实现

* [httprouter](httprouter)基于[httprouter](github.com/julienschmidt/httprouter)的高效率路由实现

## 路由参数

  //获取路由参数  
  params:=router.GetParams(r)

  //设置路由参数
  params.Set("paramname","value")

  //获取路由参数
  v=parans.Get("paramname")