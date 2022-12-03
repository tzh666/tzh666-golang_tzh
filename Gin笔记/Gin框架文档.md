## Gin框架入门 

### 一、Gin安装

```go
// 安装
go get -u github.com/gin-gonic/gin

// 使用
import "github.com/gin-gonic/gin"
```

#### 1.1、简单例子

- 测试访问：curl [127.0.0.1:8080](http://127.0.0.1:8080/)

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、绑定路由规则,执行的函数
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}
```



### 二、Gin路由

#### 2.1、基本路由

- gin 框架中采用的路由库是基于httprouter做的
- 地址为：https://github.com/julienschmidt/httprouter

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、绑定路由规则,执行的函数
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}
```

#### 2.2、路由分组

- 路由分组是为了管理一些相同的URL
- 测试方法：
  - curl  "http://127.0.0.1:8080/v1/login?name=xiaoli"
  - curl -X POST 'http://127.0.0.1:8080/v2/submit?name=xiaoli' -d '{}'

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由, 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 2、路由组
	v1 := r.Group("/v1")
	// {} 是书写规范
	{
		v1.GET("/login", login)
		v1.GET("/submit", submit)
	}

	v2 := r.Group("/v2")
	{
		v2.POST("/login", login)
		v2.POST("/submit", submit)
	}

	// 3、设置gin启动模式
	// gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.DebugMode) // 发布模式, 屏蔽debug日志

	// 4、启动监听,默认在8080
	r.Run(":8080")
}

func login(c *gin.Context) {
	// 获取参数,默认jack
	name := c.DefaultQuery("name", "jack")
	c.String(http.StatusOK, name)
}

func submit(c *gin.Context) {
	// 获取参数,默认libai
	name := c.DefaultQuery("name", "libai")
	c.String(http.StatusOK, name)
}
```

#### 2.3、路由拆分与注册

- 首先看下基本路由：适用于路由条目比较少的简单项目或者项目demo
- 然后我们基于这个demo进行拆分

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello Gin!",
	})
}

func main() {
	r := gin.Default()

	r.GET("/", helloHandler)

	r.Run(":8080")
}
```

##### 2.3.1、路由拆分成单独文件或包

- 当项目的规模增大后就不太适合继续在项目的main.go文件中去实现路由注册相关逻辑了，我们会倾向于把路由部分的代码都拆分出来，形成一个单独的文件或包：

- 首先创建一个routers文件夹，然后在文件夹中创建routers.go

```go
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello Gin!",
	})
}

// 定义一个函数,返回值是*gin.Engine【这样是一种设计模式】
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", helloHandler)
	return r
}
```

- 再看main.go中如何调用路由的代码

```go
package main

import (
	"demo16/routers"
	"fmt"
)

func main() {
	// 调用routers.go定义好的setupRouter函数
	r := routers.SetupRouter()

	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
```

- 此时的文件目录层级

```sh
demo16
├── go.mod
├── go.sum
├── main.go
└── routers
    └── routers.go
```

##### 2.3.2、路由拆分成多个文件

当我们的业务规模继续膨胀，单独的一个routers文件或包已经满足不了我们的需求

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.GET("/demo01", helloHandler)
    r.GET("/demo02", xxHandler1)
    ...
    ...
    r.GET("/demo50", xxHandler30)
    return r
}
```

因为我们把所有的路由注册都写在一个SetupRouter函数中的话就会太复杂了。

- 我们可以分开定义多个路由文件，例如：

```sh
demo16
├── go.mod
├── go.sum
├── main.go
└── routers
    ├── blog.go
    └── shop.go
```

- routers/shop.go中添加一个LoadShop的函数，将shop相关的路由注册到指定的路由器：

```go
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Handler",
	})
}

func checkoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "checkoutHandler",
	})
}

func goodsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "goodsHandler",
	})
}

// 路由注册
func LoadShop(c *gin.Engine) {
	c.GET("/hello", helloHandler)
	c.GET("/goods", goodsHandler)
	c.GET("/checkout", checkoutHandler)
}
```

- routers/blog.go中添加一个LoadBlog的函数，将blog相关的路由注册到指定的路由器：

```go
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func postHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "postHandler",
	})
}

func commentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "commentHandler",
	})
}
func LoadBlog(c *gin.Engine) {
	c.GET("/post", postHandler)
	c.GET("/comment", commentHandler)
}
```

- **在main函数中实现最终的路由注册逻辑如下：**

```go
package main

import (
	"demo16/routers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 在main函数中实现最终的注册逻辑如下, 路由注册
	r := gin.Default()
	routers.LoadBlog(r)
	routers.LoadShop(r)

	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
```

##### 2.3.3、路由拆分到不同的APP

有时候项目规模实在太大，那么我们就更倾向于把业务拆分的更详细一些，例如把不同的业务代码拆分成不同的APP。

因此我们在项目目录下单独定义一个app目录，用来存放我们不同业务线的代码文件，这样就很容易进行横向扩展。大致目录结构如下：

```sh
demo17
├── app
│   ├── blog
│   │   ├── handler.go
│   │   └── router.go
│   └── shop
│       ├── handler.go
│       └── router.go
├── go.mod
├── go.sum
├── main.go
└── routers
    └── routers.go
```

其中app/blog/router.go用来定义post相关路由信息，具体内容如下：

```go
package blog

import (
	"github.com/gin-gonic/gin"
)
// postHandler、commentHandler届时放在controller里面即可
// 路由规则
func Routers(e *gin.Engine) {
	e.GET("/post", postHandler)
	e.GET("/comment", commentHandler)
}

```

app/shop/router.go用来定义shop相关路由信息，具体内容如下：

```go
package shop

import (
	"github.com/gin-gonic/gin"
)

// 路由规则
// goodsHandler、checkoutHandler届时放在controller里面即可
func Routers(e *gin.Engine) {
	e.GET("/goods", goodsHandler)
	e.GET("/checkout", checkoutHandler)
}

```

routers/routers.go中根据需要定义Include函数用来注册子app中定义的路由，Init函数用来进行路由的初始化操作：

```go
package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

// 函数类型的组数
var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	// 把所有路由都放到这个组数中
	options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
	// 初始化没有中间件方法
	r := gin.New()
	// 遍历这个包含了全部路由的数组,opt遍历出来是函数
	for _, opt := range options {
		// 路由注册
		opt(r)
		fmt.Println(opt)
	}
	return r
}
```

main.go中按如下方式先注册子app中的路由，然后再进行路由的初始化：

```go
package main

import (
	"demo17/app/blog"
	"demo17/app/shop"
	"demo17/routers"
	"fmt"
)

func main() {
	// 加载多个APP的路由配置
	routers.Include(shop.Routers, blog.Routers)
	// 初始化路由
	r := routers.Init()
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}

	// 启动
	r.Run(":8080")
}
```



### 三、Gin参数绑定

#### 3.1、Restful风格的API

- gin支持Restful风格的API
- 即Representational State Transfer的缩写。直接翻译的意思是"表现层状态转化"，是一种互联网应用程序的API设计理念：URL定位资源，用HTTP描述操作

```sh
1.获取文章 /blog/getXxx Get blog/Xxx

2.添加 /blog/addXxx POST blog/Xxx

3.修改 /blog/updateXxx PUT blog/Xxx

4.删除 /blog/delXxxx DELETE blog/Xxx
```

#### 3.2、Gin参数获取

##### 3.2.1、API参数获取

- 可以通过Context的Param方法来获取API参数

```go
package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、绑定路由规则,执行的函数
	r.GET("/user/:name/*action", func(c *gin.Context) {
		// 获取参数
		name := c.Param("name")
		action := c.Param("action")
		// 截取 / 符号, 如果不截取获取到的第二个参数带有 /
		action = strings.Trim(action, "/")
		c.String(http.StatusOK, name+" is "+action)
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}

# 演示
$ curl http://127.0.0.1:8080/user/zg/sz
zg is sz
```

##### 3.2.2、URL参数获取

- **URL参数可以通过DefaultQuery()或Query()方法获取** 
- DefaultQuery()若参数不存在，返回默认值，Query()若不存在，返回空串
- API ? name=sz，name 就是我们url上设置参数名,传递的时候通过?name=sz传递

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、绑定路由规则,执行的函数
	r.GET("/user", func(c *gin.Context) {
		// name 就是我们url上设置参数名,传递的时候通过?name=sz传递
		name := c.DefaultQuery("name", "默认值")
		c.String(http.StatusOK, fmt.Sprintf("传递的参数是: %s", name))
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}

// 演示
$ curl http://127.0.0.1:8080/user
传递的参数是: 默认值

$ curl http://127.0.0.1:8080/user?name=sz
传递的参数是: sz
```

##### 3.2.3、表单参数获取

- **表单传输为post请求，http常见的传输格式为四种：**
  - application/json
  - application/x-www-form-urlencoded
  - application/xml
  - multipart/form-data
- 表单参数可以通过PostForm()方法获取，该方法默认解析的是x-www-form-urlencoded或from-data格式的参数

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、绑定路由规则,执行的函数
	r.POST("form", func(c *gin.Context) {
		// 调用接口POST方式调用
		types := c.DefaultPostForm("type", "post")
		// 获取参数
		username := c.PostForm("username")
		password := c.PostForm("password")
		// 返回值
		c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}

// 演示
$ curl -X POST 'http://127.0.0.1:8080/form?username=zg&password=sz'
username:,password:,type:post
```

#### 3.3、Gin文件上传

##### 3.1、单文件上传

- multipart/form-data格式用于文件上传
- gin文件上传与原生的net/http方法类似，不同在于gin把原生的request封装到c.Request中

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、限制最大上传尺寸, 最大8M
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"文件上传错误: ": err,
			})
		}
		// 保存文件
		c.SaveUploadedFile(file, "./file/"+file.Filename)
		// 保存成功给个返回值
		c.String(http.StatusOK, file.Filename)
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}

// 演示
curl -X POST 'http://127.0.0.1:8080/upload' --form 'file=@".\壁纸.jpg"' 
```

##### 3.2、多文件上传

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1、创建路由
	r := gin.Default()

	// 2、限制最大上传尺寸, 最大8M
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		// 获取所有图片
		files := form.File["files"]

		// 遍历全部上传的文件
		for _, file := range files {
			// 逐个存储
			fmt.Println("./file/" + file.Filename)
			if err := c.SaveUploadedFile(file, "./file/"+file.Filename); err != nil {
				// 接口返回值
				c.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}
		}
		c.String(200, fmt.Sprintf("upload ok %d files", len(files)))
	})

	// 3、启动监听,默认在8080
	r.Run(":8080")
}

// 演示
curl -X POST 'http://127.0.0.1:8080/upload' --form 'files=@".\壁纸.jpg"'  --form 'files=@".\壁纸1.jpg"' 
```



### 四、Gin数据解析和绑定

#### 4.1、Json 数据解析和绑定

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义接收数据的结构体
type Login struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	// 1.创建路由,默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 2、Json绑定
	r.POST("/login", func(c *gin.Context) {
		// 声明接收前端传入参数的变量
		var user Login
		// 将request的body中的数据，自动按照json格式解析到结构体
		if err := c.ShouldBind(&user); err != nil {
			// 返回错误信息
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			// 失败则不执行下面的代码
			return
		}
		// 假设判断用户密码是否正确
		fmt.Println(user)
		if user.User != "root" || user.Password != "admin" {
			// 密码错误返回信息
			c.JSON(http.StatusBadRequest, gin.H{
				"status": 401,
			})
			return
		}
		// 密码正确则返回200
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "登录成功",
		})
	})

	// 3、启动
	r.Run(":8080")
}

// 演示
$ curl --location --request POST 'http://127.0.0.1:8080/login' --header 'Content-Type: application/json' --data-raw '{
    "user": "root",
    "password": "admin"
}'
{"message":"登录成功","status":200}
```

#### 4.2、单数据解析和绑定

- Bind()默认解析并绑定form格式
- 根据请求头中content-type自动推断

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义接收数据的结构体
type Login struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	// 1、创建路由,默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 2、
	r.POST("/loginForm", func(c *gin.Context) {
		// 声明接收的变量
		var form Login
		// Bind()默认解析并绑定form格式、根据请求头中content-type自动推断
		if err := c.Bind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(form)
		// 判断用户名密码是否正确
		if form.User != "root" || form.Password != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "401",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "200",
			"message": "登陆成功",
		})
	})
	// 3、启动
	r.Run(":8080")
}

// 演示
$ curl --location --request POST 'http://127.0.0.1:8080/loginForm?user=root&password=admin'
{"message":"登陆成功","status":"200"}
```

#### 4.3、URI数据解析和绑定

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义接收数据的结构体
type Login struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	User     string `uri:"user" binding:"required"`
	Password string `uri:"password" binding:"required"`
}

func main() {
	// 1、创建路由,默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 2、
	r.GET("/:user/:password", func(c *gin.Context) {
		// 声明接收的变量
		var login Login
		// Bind()默认解析并绑定form格式
		// 根据请求头中content-type自动推断
		if err := c.ShouldBindUri(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 判断用户名密码是否正确
		if login.User != "root" || login.Password != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "200"})
	})
	// 3、启动
	r.Run(":8080")
}


// 演示
curl --location --request GET 'http://127.0.0.1:8080/root/admin'
```

#### 4.4、xml数据解析和绑定

```go
--- 暂无
```



### 五、Gin数据渲染

#### 5.1、各种数据格式的响应

- **后端接口相应给前端的数据**

- json、结构体、XML、YAML

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
)

func main() {
	// 1、创建路由,默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 2、数据响应
	//    2.1、JSON
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "相应JSON格式数据",
			"status":  200,
		})
	})
	//    2.2、结构体响应
	r.GET("/struct", func(c *gin.Context) {
		type Msg struct {
			Name     string
			Password string
			Number   int
		}
		msg := Msg{"root", "123456", 1}

		c.JSON(200, msg)
	})
	//    2.3、XML响应
	r.GET("/xml", func(c *gin.Context) {
		c.XML(200, gin.H{
			"message": "abc",
		})
	})
	//    2.4、YAML响应
	r.GET("/yaml", func(c *gin.Context) {
		c.YAML(
			200, gin.H{
				"name": "zhangsan",
			})
	})
	//    2.5、ProtoBuf 格式响应,谷歌开发的高效存储读取的工具
	// 数组？切片？如果自己构建一个传输格式，应该是什么格式？
	r.GET("/ProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		// 定义数据
		label := "label"
		// 传ProtoBuf格式数据
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		c.ProtoBuf(200, data)
	})

	// 3、启动
	r.Run(":8080")
}
```

#### 5.2、HTML模板渲染

- gin支持加载HTML模板, 然后根据模板参数进行配置并返回相应的数据，本质上就是字符串替换
- LoadHTMLGlob()方法可以加载模板文件
- 如果你需要引入静态文件需要定义一个静态文件目录 【r.Static("/assets", "./assets")】

```go
func main() {
	r := gin.Default()
	// 加载模板文件
	r.LoadHTMLGlob("tem/*")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
	})

	r.Run(":8080")
}
```

```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.title}}</title>
</head>

<body>
    {{.ce}}
</body>

</html>
```

#### 5.3、重定向

- 重定向：Redirect

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/index", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "http://www.5lmh.com")
    })
    r.Run()
}
```

#### 5.4、 同步异步

- goroutine机制可以方便地实现异步处理，**异步执行前端不会有感知，逻辑会 在后台处理**
- 另外，在启动新的goroutine时，不应该使用原始上下文，必须使用它的只读副本
- 不使用goroutine就是默认同步处理

```go
package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.创建路由、默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()
	// 1.1、异步
	r.GET("/long_async", func(c *gin.Context) {
		// 需要搞一个副本
		copyContext := c.Copy()
		// 异步处理
		go func() {
			time.Sleep(3 * time.Second)
			log.Println("异步执行：" + copyContext.Request.URL.Path)
		}()
	})
	// 1.2、同步
	r.GET("/long_sync", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		log.Println("同步执行：" + c.Request.URL.Path)
	})
	r.Run(":8080")
}
```



### 六、Gin中间件

#### 6.1、全局中间件

- 所有请求都经过此中间件

```go
package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义中间件
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前时间
		t := time.Now()
		fmt.Println("中间件开始执行了")

		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "中间件")
		status := c.Writer.Status()
		fmt.Println("中间件执行完毕", status)

		// time.Since计算执行时间
		t2 := time.Since(t)
		fmt.Println("time:", t2)
	}
}

func main() {
	// 1、创建路由、默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 2、注册中间件
	r.Use(MiddleWare())

	r.GET("/", func(c *gin.Context) {
		// 取自定义环境变量的值
		req, _ := c.Get("request")
		// 页面接收
		c.JSON(200, gin.H{"request": req})
	})

	r.Run(":8080")
}
```

#### 6.2、Next()方法

- Next 函数会挂起当前所在的函数，然后调用后面的中间件，待后面中间件执行完毕后，再接着执行当前函数

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

// 定义中间件1
func MiddleWare1() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("调用中间件1")

        // 调用 Next，开始执行后续的中间件
        c.Next()

        fmt.Println("中间件继续执行")
    }
}

// 定义中间件2
func MiddleWare2() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("调用中间件2")
    }
}

// 定义中间件3
func MiddleWare3() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("调用中间件3")
    }
}

func main() {
    // 创建路由
    engine := gin.Default()

    // 注册中间件
    engine.Use(MiddleWare1(), MiddleWare2(), MiddleWare3())

    // 路由规则
    engine.GET("/", func(c *gin.Context) {
        fmt.Println("调用路由处理函数")
        // 页面接收
        c.JSON(200, gin.H{"request": "编程宝库 gin框架"})
    })
    engine.Run()
}
```

运行程序，并在浏览器输入：http://localhost:8080，控制台日志会输出：

```
[GIN-debug] GET    /   --> main.main.func1 (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080
调用中间件1
调用中间件2
调用中间件3
调用路由处理函数
中间件继续执行
[GIN] 2021/05/31 - 12:03:13 | 200 |  193.22µs | ::1 | GET "/"
```

以上输出显示：中间件1执行 Next 后，开始执行其他中间件，再执行页面处理函数。它们执行完毕后，又开始执行自己后续的代码。

#### 6.3、Abort 方法

- Abort 函数在被调用的函数中阻止后续中间件的执行。例如，你有一个验证当前的请求是否是认证过的 Authorization 中间件。如果验证失败(例如，密码不匹配)，调用 Abort 以确保这个请求的其他函数不会被调用。
- 有个细节需要注意，调用 Abort 函数不会停止当前的函数的执行，除非后面跟着 return

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

// 定义中间件1
func MiddleWare1() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("调用中间件1")

        // 调用 Abort，终止执行后续的中间件
        c.Abort()
    }
}

// 定义中间件2
func MiddleWare2() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("调用中间件2")
    }
}

// 定义中间件3
func MiddleWare3() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("调用中间件3")
    }
}

func main() {
    // 创建路由
    engine := gin.Default()

    // 注册中间件
    engine.Use(MiddleWare1(), MiddleWare2(), MiddleWare3())

    // 路由规则
    engine.GET("/", func(c *gin.Context) {
        fmt.Println("调用路由处理函数")
        // 页面接收
        c.JSON(200, gin.H{"request": "编程宝库 gin框架"})
    })

    engine.Run()
}
```

运行程序，并在浏览器输入：http://localhost:8080，控制台日志会输出：

```
[GIN-debug] GET    /   --> main.main.func1 (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080
调用中间件1
[GIN] 2021/05/31 - 12:03:13 | 200 |  193.22µs | ::1 | GET "/"
```

说明只有中间件1被调用，其余的中间件被取消执行，而且页面处理函数也被取消执行。

#### 6.4、局部中间件

```go
package main

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

// 定义中间
func MiddleWare() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        fmt.Println("中间件开始执行了")
        // 设置变量到Context的key中，可以通过Get()取
        c.Set("request", "中间件")
        // 执行函数
        c.Next()
        // 中间件执行完后续的一些事情
        status := c.Writer.Status()
        fmt.Println("中间件执行完毕", status)
        t2 := time.Since(t)
        fmt.Println("time:", t2)
    }
}

func main() {
    // 1.创建路由
    // 默认使用了2个中间件Logger(), Recovery()
    r := gin.Default()
    //局部中间键使用
    r.GET("/ce", MiddleWare(), func(c *gin.Context) {
        // 取值
        req, _ := c.Get("request")
        fmt.Println("request:", req)
        // 页面接收
        c.JSON(200, gin.H{"request": req})
    })
    r.Run()
}
```

#### 6.5、中间件推荐

**谷歌翻译欢迎查看原文 https://github.com/gin-gonic/contrib/blob/master/README.md**

- [RestGate](https://github.com/pjebs/restgate) - REST API端点的安全身份验证
- [staticbin](https://github.com/olebedev/staticbin) - 用于从二进制数据提供静态文件的中间件/处理程序
- [gin-cors](https://github.com/gin-contrib/cors) - CORS杜松子酒的官方中间件
- [gin-csrf](https://github.com/utrack/gin-csrf) - CSRF保护
- [gin-health](https://github.com/utrack/gin-health) - 通过[gocraft/health](https://github.com/gocraft/health)报告的中间件
- [gin-merry](https://github.com/utrack/gin-merry) - 带有上下文的漂亮 [打印](https://github.com/ansel1/merry) 错误的中间件
- [gin-revision](https://github.com/appleboy/gin-revision-middleware) - 用于Gin框架的修订中间件
- [gin-jwt](https://github.com/appleboy/gin-jwt) - 用于Gin框架的JWT中间件
- [gin-sessions](https://github.com/kimiazhu/ginweb-contrib/tree/master/sessions) - 基于mongodb和mysql的会话中间件
- [gin-location](https://github.com/drone/gin-location) - 用于公开服务器的主机名和方案的中间件
- [gin-nice-recovery](https://github.com/ekyoung/gin-nice-recovery) - 紧急恢复中间件，可让您构建更好的用户体验
- [gin-limit](https://github.com/aviddiviner/gin-limit) - 限制同时请求；可以帮助增加交通流量
- [gin-limit-by-key](https://github.com/yangxikun/gin-limit-by-key) - 一种内存中的中间件，用于通过自定义键和速率限制访问速率。
- [ez-gin-template](https://github.com/michelloworld/ez-gin-template) - gin简单模板包装
- [gin-hydra](https://github.com/janekolszak/gin-hydra) - gin中间件[Hydra](https://github.com/ory-am/hydra)
- [gin-glog](https://github.com/zalando/gin-glog) - 旨在替代Gin的默认日志
- [gin-gomonitor](https://github.com/zalando/gin-gomonitor) - 用于通过Go-Monitor公开指标
- [gin-oauth2](https://github.com/zalando/gin-oauth2) - 用于OAuth2
- [static](https://github.com/hyperboloide/static) gin框架的替代静态资产处理程序。
- [xss-mw](https://github.com/dvwright/xss-mw) - XssMw是一种中间件，旨在从用户提交的输入中“自动删除XSS”
- [gin-helmet](https://github.com/danielkov/gin-helmet) - 简单的安全中间件集合。
- [gin-jwt-session](https://github.com/ScottHuangZL/gin-jwt-session) - 提供JWT / Session / Flash的中间件，易于使用，同时还提供必要的调整选项。也提供样品。
- [gin-template](https://github.com/foolin/gin-template) - 用于gin框架的html / template易于使用。
- [gin-redis-ip-limiter](https://github.com/Salvatore-Giordano/gin-redis-ip-limiter) - 基于IP地址的请求限制器。它可以与redis和滑动窗口机制一起使用。
- [gin-method-override](https://github.com/bu/gin-method-override) - _method受Ruby的同名机架启发而被POST形式参数覆盖的方法
- [gin-access-limit](https://github.com/bu/gin-access-limit) - limit-通过指定允许的源CIDR表示法的访问控制中间件。
- [gin-session](https://github.com/go-session/gin-session) - 用于Gin的Session中间件
- [gin-stats](https://github.com/semihalev/gin-stats) - 轻量级和有用的请求指标中间件
- [gin-statsd](https://github.com/amalfra/gin-statsd) - 向statsd守护进程报告的Gin中间件
- [gin-health-check](https://github.com/RaMin0/gin-health-check) - check-用于Gin的健康检查中间件
- [gin-session-middleware](https://github.com/go-session/gin-session) - 一个有效，安全且易于使用的Go Session库。
- [ginception](https://github.com/kubastick/ginception) - 漂亮的例外页面
- [gin-inspector](https://github.com/fatihkahveci/gin-inspector) - 用于调查http请求的Gin中间件。
- [gin-dump](https://github.com/tpkeeper/gin-dump) - Gin中间件/处理程序，用于转储请求和响应的标头/正文。对调试应用程序非常有帮助。
- [go-gin-prometheus](https://github.com/zsais/go-gin-prometheus) - Gin Prometheus metrics exporter
- [ginprom](https://github.com/chenjiandongx/ginprom) - Gin的Prometheus指标导出器
- [gin-go-metrics](https://github.com/bmc-toolbox/gin-go-metrics) - Gin middleware to gather and store metrics using [rcrowley/go-metrics](https://github.com/rcrowley/go-metrics)
- [ginrpc](https://github.com/xxjwxc/ginrpc) - Gin 中间件/处理器自动绑定工具。通过像beego这样的注释路线来支持对象注册



### 七、会话控制Cookie、Sessions

#### 7.1、Cookie介绍

- HTTP是无状态协议，服务器不能记录浏览器的访问状态，也就是说服务器不能区分两次请求是否由同一个客户端发出
- Cookie就是解决HTTP协议无状态的方案之一，中文是小甜饼的意思
- Cookie实际上就是服务器保存在浏览器上的一段信息。浏览器有了Cookie之后，每次向服务器发送请求时都会同时将该信息发送给服务器，服务器收到请求后，就可以根据该信息处理请求
- Cookie由服务器创建，并发送给浏览器，最终由浏览器保存

##### 7.1.1、Cookie的用途

- 测试服务端发送cookie给客户端，客户端请求时携带cookie

##### 7.1.2、Cookie的缺点

- 不安全，明文
- 增加带宽消耗
- 可以被禁用
- cookie有上限

#### 7.2、Cookie模拟实现权限验证中间件

```go
package main

import (
   "github.com/gin-gonic/gin"
   "net/http"
)

func AuthMiddleWare() gin.HandlerFunc {
   return func(c *gin.Context) {
      // 获取客户端cookie并校验
      if cookie, err := c.Cookie("abc"); err == nil {
         if cookie == "123" {
            c.Next()
            return
         }
      }
      // 返回错误
      c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
      // 若验证不通过，不再调用后续的函数处理
      c.Abort()
      return
   }
}

func main() {
   // 1.创建路由
   r := gin.Default()
   r.GET("/login", func(c *gin.Context) {
      // 设置cookie
      c.SetCookie("abc", "123", 60, "/",
         "localhost", false, true)
      // 返回信息
      c.String(200, "Login success!")
   })
   r.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
      c.JSON(200, gin.H{"data": "home"})
   })
   r.Run(":8000")
}
```

#### 7.3、Sessions

gorilla/sessions为自定义session后端提供cookie和文件系统session以及基础结构。

主要功能是：

- 简单的API：将其用作设置签名（以及可选的加密）cookie的简便方法。
- 内置的后端可将session存储在cookie或文件系统中。
- Flash消息：一直持续读取的session值。
- 切换session持久性（又称“记住我”）和设置其他属性的便捷方法。
- 旋转身份验证和加密密钥的机制。
- 每个请求有多个session，即使使用不同的后端也是如此。
- 自定义session后端的接口和基础结构：可以使用通用API检索并批量保存来自不同商店的session。

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/sessions"
)

// 初始化一个cookie存储对象
// something-very-secret应该是一个你自己的密匙，只要不被别人知道就行
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
    http.HandleFunc("/save", SaveSession)
    http.HandleFunc("/get", GetSession)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("HTTP server failed,err:", err)
        return
    }
}

func SaveSession(w http.ResponseWriter, r *http.Request) {
    // Get a session. We're ignoring the error resulted from decoding an
    // existing session: Get() always returns a session, even if empty.

    //　获取一个session对象，session-name是session的名字
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 在session中存储值
    session.Values["foo"] = "bar"
    session.Values[42] = 43
    // 保存更改
    session.Save(r, w)
}
func GetSession(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    foo := session.Values["foo"]
    fmt.Println(foo)
}
```

- 删除session的值：

```
    // 删除
    // 将session的最大存储时间设置为小于零的数即为删除
    session.Options.MaxAge = -1
    session.Save(r, w)
```



### 八、参数验证

#### 8.1、结构体验证

- 用gin框架的数据验证，可以不用解析数据，减少if else，会简洁许多。

```go
package main

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

//Person ..
type Person struct {
    //不能为空并且大于10
    Age      int       `form:"age" binding:"required,gt=10"`
    Name     string    `form:"name" binding:"required"`
    Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
    r := gin.Default()
    r.GET("/5lmh", func(c *gin.Context) {
        var person Person
        if err := c.ShouldBind(&person); err != nil {
            c.String(500, fmt.Sprint(err))
            return
        }
        c.String(200, fmt.Sprintf("%#v", person))
    })
    r.Run()
}
```

#### 8.2、自定义验证

```go
package main

import (
    "net/http"
    "reflect"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "gopkg.in/go-playground/validator.v8"
)

/*
    对绑定解析到结构体上的参数，自定义验证功能
    比如我们要对 name 字段做校验，要不能为空，并且不等于 admin ，类似这种需求，就无法 binding 现成的方法
    需要我们自己验证方法才能实现 官网示例（https://godoc.org/gopkg.in/go-playground/validator.v8#hdr-Custom_Functions）
    这里需要下载引入下 gopkg.in/go-playground/validator.v8
*/
type Person struct {
    Age int `form:"age" binding:"required,gt=10"`
    // 2、在参数 binding 上使用自定义的校验方法函数注册时候的名称
    Name    string `form:"name" binding:"NotNullAndAdmin"`
    Address string `form:"address" binding:"required"`
}

// 1、自定义的校验方法
func nameNotNullAndAdmin(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {

    if value, ok := field.Interface().(string); ok {
        // 字段不能为空，并且不等于  admin
        return value != "" && !("5lmh" == value)
    }

    return true
}

func main() {
    r := gin.Default()

    // 3、将我们自定义的校验方法注册到 validator中
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // 这里的 key 和 fn 可以不一样最终在 struct 使用的是 key
        v.RegisterValidation("NotNullAndAdmin", nameNotNullAndAdmin)
    }

    /*
        curl -X GET "http://127.0.0.1:8080/testing?name=&age=12&address=beijing"
        curl -X GET "http://127.0.0.1:8080/testing?name=lmh&age=12&address=beijing"
        curl -X GET "http://127.0.0.1:8080/testing?name=adz&age=12&address=beijing"
    */
    r.GET("/5lmh", func(c *gin.Context) {
        var person Person
        if e := c.ShouldBind(&person); e == nil {
            c.String(http.StatusOK, "%v", person)
        } else {
            c.String(http.StatusOK, "person bind err:%v", e.Error())
        }
    })
    r.Run()
}
```

```go
package main

import (
    "net/http"
    "reflect"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "gopkg.in/go-playground/validator.v8"
)

// Booking contains binded and validated data.
type Booking struct {
    //定义一个预约的时间大于今天的时间
    CheckIn time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
    //gtfield=CheckIn退出的时间大于预约的时间
    CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(
    v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
    field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
    //field.Interface().(time.Time)获取参数值并且转换为时间格式
    if date, ok := field.Interface().(time.Time); ok {
        today := time.Now()
        if today.Unix() > date.Unix() {
            return false
        }
    }
    return true
}

func main() {
    route := gin.Default()
    //注册验证
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        //绑定第一个参数是验证的函数第二个参数是自定义的验证函数
        v.RegisterValidation("bookabledate", bookableDate)
    }

    route.GET("/5lmh", getBookable)
    route.Run()
}

func getBookable(c *gin.Context) {
    var b Booking
    if err := c.ShouldBindWith(&b, binding.Query); err == nil {
        c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
}

// curl -X GET "http://localhost:8080/5lmh?check_in=2019-11-07&check_out=2019-11-20"
// curl -X GET "http://localhost:8080/5lmh?check_in=2019-09-07&check_out=2019-11-20"
// curl -X GET "http://localhost:8080/5lmh?check_in=2019-11-07&check_out=2019-11-01"
```

#### 8.3、多语言翻译验证

当业务系统对验证信息有特殊需求时，例如：返回信息需要自定义，手机端返回的信息需要是中文而pc端发挥返回的信息需要时英文，如何做到请求一个接口满足上述三种情况。

```go
package main

import (
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/locales/en"
    "github.com/go-playground/locales/zh"
    "github.com/go-playground/locales/zh_Hant_TW"
    ut "github.com/go-playground/universal-translator"
    "gopkg.in/go-playground/validator.v9"
    en_translations "gopkg.in/go-playground/validator.v9/translations/en"
    zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
    zh_tw_translations "gopkg.in/go-playground/validator.v9/translations/zh_tw"
)

var (
    Uni      *ut.UniversalTranslator
    Validate *validator.Validate
)

type User struct {
    Username string `form:"user_name" validate:"required"`
    Tagline  string `form:"tag_line" validate:"required,lt=10"`
    Tagline2 string `form:"tag_line2" validate:"required,gt=1"`
}

func main() {
    en := en.New()
    zh := zh.New()
    zh_tw := zh_Hant_TW.New()
    Uni = ut.New(en, zh, zh_tw)
    Validate = validator.New()

    route := gin.Default()
    route.GET("/5lmh", startPage)
    route.POST("/5lmh", startPage)
    route.Run(":8080")
}

func startPage(c *gin.Context) {
    //这部分应放到中间件中
    locale := c.DefaultQuery("locale", "zh")
    trans, _ := Uni.GetTranslator(locale)
    switch locale {
    case "zh":
        zh_translations.RegisterDefaultTranslations(Validate, trans)
        break
    case "en":
        en_translations.RegisterDefaultTranslations(Validate, trans)
        break
    case "zh_tw":
        zh_tw_translations.RegisterDefaultTranslations(Validate, trans)
        break
    default:
        zh_translations.RegisterDefaultTranslations(Validate, trans)
        break
    }

    //自定义错误内容
    Validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
        return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
    }, func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("required", fe.Field())
        return t
    })

    //这块应该放到公共验证方法中
    user := User{}
    c.ShouldBind(&user)
    fmt.Println(user)
    err := Validate.Struct(user)
    if err != nil {
        errs := err.(validator.ValidationErrors)
        sliceErrs := []string{}
        for _, e := range errs {
            sliceErrs = append(sliceErrs, e.Translate(trans))
        }
        c.String(200, fmt.Sprintf("%#v", sliceErrs))
    }
    c.String(200, fmt.Sprintf("%#v", "user"))
}
```

正确的链接：http://localhost:8080/testing?user_name=枯藤&tag_line=9&tag_line2=33&locale=zh

http://localhost:8080/testing?user_name=枯藤&tag_line=9&tag_line2=3&locale=en 返回英文的验证信息

http://localhost:8080/testing?user_name=枯藤&tag_line=9&tag_line2=3&locale=zh 返回中文的验证信息

查看更多的功能可以查看官网 gopkg.in/go-playground/validator.v9