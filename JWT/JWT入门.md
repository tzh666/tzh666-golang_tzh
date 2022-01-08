### 一、何为JWT

JWT全称JSON Web Token是一种跨域认证解决方案，属于一个开放的标准，它规定了一种Token实现方式，目前多用于前后端分离项目和OAuth2.0业务场景下。



### 二、为什么需要JWT

在之前的一些web项目中，我们通常使用的是`Cookie-Session`模式实现用户认证。相关流程大致如下：

1. 用户在浏览器端填写用户名和密码，并发送给服务端
2. 服务端对用户名和密码校验通过后会生成一份保存当前用户相关信息的session数据和一个与之对应的标识（通常称为session_id）
3. 服务端返回响应时将上一步的session_id写入用户浏览器的Cookie
4. 后续用户来自该浏览器的每次请求都会自动携带包含session_id的Cookie
5. 服务端通过请求中的session_id就能找到之前保存的该用户那份session数据，从而获取该用户的相关信息。

这种方案依赖于客户端（浏览器）保存Cookie，并且需要在服务端存储用户的session数据。

在移动互联网时代，我们的用户可能使用浏览器也可能使用APP来访问我们的服务，我们的web应用可能是前后端分开部署在不同的端口，**有时候我们还需要支持第三方登录**，这下`Cookie-Session`的模式就有些力不从心了。

**JWT就是一种基于Token的轻量级认证模式**，服务端认证通过后，会生成一个JSON对象，经过签名后得到一个Token（令牌）再发回给用户，用户后续请求只需要带上这个Token，服务端解密之后就能获取该用户的相关信息了。



### 三、生成JWT

使用`dgrijalva/jwt-goo`这个库来实现我们生成JWT和解析JWT的功能

#### 3.1、定义需求

`定制自己的需求来决定JWT中保存哪些数据，比如我们规定在JWT中要存储username信息，那么我们就定义一个MyClaims结构体如下：`

```go
// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
```

#### 3.2、自定义JWT过期时间跟Secret

```go
const TokenExpireDuration = time.Hour * 1

var MySecret = []byte("CXVBNMYHUJKI")
```

#### 3.3、生成JWT

```go
func main() {
	c := MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "tzh666",                                   // 签发人
		},
		"tzh666",
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	my, err := token.SignedString(MySecret)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(my)
}
```

#### 3.4、函数封装

```go
type MyClaims struct {
	jwt.StandardClaims        // jwt包自带的jwt.StandardClaims只包含了官方字段
	UserName           string `json:"username"` // 自定义要存储的信息,可有多个
}

// 设置过期时间
const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("CXVBNMYHUJKI")

func GenToken(username string) (string, error) {
	c := MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "tzh666",                                   // 签发人
		},
		username,
	}
	/*
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
		token.SignedString(MySecret)
	*/
	// 使用指定的签名方法创建签名对象,此处指定用S256,再调用SignedString函数使用指定的secret签名并获得完整的编码后的字符串token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(MySecret)
}
```



### 四、解析JWT

#### 4.1、解析token

```go
func main() {
	Token, err := GenToken("tzh666")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 解析token
	token, err := jwt.ParseWithClaims(Token, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// .(*MyClaims)断言 token校验
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Println(claims)
	}

}
```

#### 4.2、函数封装（token解析）

```go
func ParseToken(Token string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(Token, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	// 校验token
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
```

