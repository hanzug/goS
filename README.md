## 整体架构

![image-20240314123839722](https://raw.githubusercontent.com/hanzug/images/master/images/image-20240314123839722.png)





## Etcd + grpc 实现服务注册和发现

### 服务端注册

服务端启动时，需要将自己的地址注册到etcd中。这样，客户端就可以通过查询etcd来发现服务的地址。

1. **连接到etcd**：服务启动时，首先建立与etcd的连接。
2. **注册服务地址**：将服务的地址（通常是IP和端口）作为键值对存储到etcd。键可以是服务名称，值是服务地址。
3. **心跳机制**：定期更新etcd中的键值对，以表明服务仍然存活。如果服务停止更新，etcd可以通过过期机制自动移除该服务，从而实现服务的健康检查。

### 客户端服务发现

客户端使用自定义解析器来查询etcd，发现服务地址，并建立连接。

1. **使用自定义解析器**：客户端在初始化时指定使用的解析器和服务名称。这个解析器负责查询etcd，找到服务的实际地址。
2. **建立gRPC连接**：客户端通过解析器返回的地址，使用gRPC的`Dial`函数建立到服务的连接。



### 关于自定义解析器

在gRPC中，自定义解析器是用于实现服务发现的关键组件。它允许客户端动态地发现和连接到服务端的实例，特别是在使用像etcd这样的服务注册中心时。自定义解析器需要实现gRPC的`resolver.Resolver`接口，并通过gRPC的解析器注册机制注册使用。

**实现`resolver.Resolver`接口**

自定义解析器至少需要实现以下方法：

- `Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)`: 当gRPC客户端尝试连接到某个服务时，这个方法被调用。它负责创建解析器的实例，并开始服务发现过程。
- `ResolveNow(o resolver.ResolveNowOptions)`: 这个方法被设计为让解析器有机会对解析请求做出响应。在某些实现中，这个方法可能不会做任何事情。
- `Close()`: 当解析器不再需要时，这个方法被调用，用于执行清理工作，比如停止后台goroutine。

**注册自定义解析器**

在使用自定义解析器之前，需要将其注册到gRPC的解析器注册机制中。这通常在程序的初始化阶段完成。例如：

```go
import (
    "google.golang.org/grpc/resolver"
)

func init() {
    resolver.Register(&myResolverBuilder{})
}
```

这里，`myResolverBuilder`是实现了`resolver.Builder`接口的类型，它负责创建自定义解析器的实例。`resolver.Register`函数用于将解析器构建器注册到gRPC中，使得gRPC能够使用这个自定义解析器来解析服务名称。

使用自定义解析器

在客户端代码中，通过指定服务名称的前缀来使用自定义解析器。这个前缀与自定义解析器在注册时使用的scheme相匹配。例如，如果你的解析器使用了`"myScheme"`作为scheme，那么在创建gRPC客户端连接时，可以这样指定服务名称：

```go
conn, err := grpc.Dial("myScheme:///serviceName", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
```

这里，`"myScheme:///serviceName"`指定了要连接的服务名称，其中`myScheme`是自定义解析器的scheme，`serviceName`是要查询的服务名称。`grpc.Dial`函数会使用与`myScheme`相匹配的解析器来解析`serviceName`。



自定义解析器允许gRPC客户端通过服务发现机制（如etcd）动态地发现服务端实例。通过实现`resolver.Resolver`接口，并将解析器注册到gRPC中，可以实现服务的动态发现和连接。这对于构建可扩展的微服务架构非常有用。



## repository

### mysql

用户表

```
type User struct {
    UserID         int64  `gorm:"primarykey"`
    UserName       string `gorm:"unique"`
    NickName       string
    PasswordDigest string
}
```

输入数据表

```
type InputData struct {
    Id      int64  `gorm:"primarykey"`
    DocId   int64  `gorm:"index"`
    Title   string `gorm:"type:longtext"`
    Body    string `gorm:"type:longtext"`
    Url     string
    Score   float64
    Source  int
    IsIndex bool
}
```

收藏表

```
type Favorite struct {
    FavoriteID     int64             `gorm:"primarykey"` // 收藏夹id
    UserID         int64             `gorm:"index"`      // 用户id
    FavoriteName   string            `gorm:"unique"`     // 收藏夹名字
    FavoriteDetail []*FavoriteDetail `gorm:"many2many:f_to_fd;"`
}
```

收藏详细表

```
type FavoriteDetail struct {
    FavoriteDetailID int64       `gorm:"primarykey"`
    UserID           int64       // 用户id
    UrlID            int64       // url的id
    Url              string      // url地址
    Desc             string      // url的描述
    Favorite         []*Favorite `gorm:"many2many:f_to_fd;"`
}
```