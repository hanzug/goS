# gateway



1. 项目网关部分使用gin框架监听http请求。



2. 注册了etcd解析器到grpc框架中，并且使用watch监听etcd中服务的变化更新到本地中。



3. 使用grpc来进行负载均衡





```go
func NewRouter() *gin.Engine {
    zap.S().Info(logs.RunFuncName())
    r := gin.Default()
    r.Use(middleware.Cors(), middleware.ErrorMiddleware())
    //store := cookie.NewStore([]byte("something-very-secret"))
    //r.Use(sessions.Sessions("mysession", store))
    v1 := r.Group("/api/v1")
    {
       v1.GET("ping", func(context *gin.Context) {
          context.JSON(200, "success")
       })
       // 用户服务
       v1.POST("/user/register", http.UserRegister)
       v1.POST("/user/login", http.UserLogin)
       // 索引平台
       IndexPlatformRegisterHandlers(v1)
       // 搜索平台
       SearchRegisterHandlers(v1)
       // 需要登录保护
       authed := v1.Group("/")
       authed.Use(middleware.AuthMiddleware())
       {
          // 收藏夹模块
          FavoriteRegisterHandlers(authed)
       }
    }
    return r
}
```