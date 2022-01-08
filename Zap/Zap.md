### 一、什么是Zap

```sh
是非常快的、结构化的，分日志级别的Go日志库
```

```sh
安装：
	go get -u go.uber.org/zap
```



### 二、入门Zap Logger

Zap提供了两种类型的日志记录器—`Sugared Logger`和`Logger`

- SugaredLogger： 
  - 在性能很好但不是很关键的上下文中
  - 比其他结构化日志记录包快4-10倍
  - 支持结构化和printf风格的日志记录
- Logger：
  - 在每一微秒和每一次内存分配都很重要的上下文中
  - 比`SugaredLogger`更快,内存分配次数也更少
  - 只支持强类型的结构化日志记录

#### 2.1、Logger记录日志记录器案例

- 通过调用`zap.NewProduction()`/`zap.NewDevelopment()`或者`zap.Example()`创建一个Logger
- 上面的每一个函数都将创建一个logger。唯一的区别在于它将记录的信息不同。例如production logger默认记录调用函数信息、日期和时间等
- 通过Logger调用Info/Error等
- 默认情况下日志都会打印到应用程序的console界面

```go
// 定义一个全局的logger
var logger *zap.Logger

// 初始化logger
func InitLogger() {
	logger, _ = zap.NewProduction()
}

func main() {
	InitLogger()
	// 延迟推出,让程序关闭之前把缓冲区的日志都刷到文件里面(或者控制台吧)
	defer logger.Sync()
	simpleHttpGet("www.baidu.com")
}

// 模拟一个请求,不管成功还是失败都记录日志
func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		logger.Info(
			"url is ok",
			zap.String("statusCode", resp.Status),
			zap.String("url", url),
		)
		resp.Body.Close()
	}
}
```

日志记录器方法的语法是这样的：

```go
func (log *Logger) MethodXXX(msg string, fields ...Field) 
```

其中`MethodXXX`是一个可变参数函数，可以是Info / Error/ Debug / Panic等。每个方法都接受一个消息字符串和任意数量的`zapcore.Field`场参数。

**每个`zapcore.Field`其实就是一组键值对参数。**

我们执行上面的代码会得到如下输出结果：

```go
{"level":"info","ts":1572159219.1227388,"caller":"zap_demo/temp.go:30","msg":"Success..","statusCode":"200 OK","url":"http://www.sogo.com"}
```



#### 2.2、SugaredLogger日志记录器案例

- 大部分的实现基本都相同
- 惟一的区别是，我们通过调用主logger的`. Sugar()`方法来获取一个`SugaredLogger`
- 然后使用`SugaredLogger`以`printf`格式记录语句

```go
var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}

func InitLogger() {
	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
```



### 三、自定义Zap Logger

**Encoder**

- **编码器(指定写入文件的格式,如json、)**
- 使用开箱即用的`NewJSONEncoder()`，并使用预先设置的`ProductionEncoderConfig()`

```go
zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
```

**WriterSyncer**

- **指定日志将写到哪里去**
- 使用`zapcore.AddSync()`函数并且将打开的文件句柄传进去

```go
file, _ := os.Create("./test.log")
return zapcore.AddSync(file)
```

**Log Level**

- 哪种级别的日志将被写入
- 有七八种级别吧

```go
zapcore.DebugLevel
```

#### 3.1、自定义日志示例代码

```go
var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func InitLogger() {
	// 指定日志将写到哪里去
	writeSyncer := getLogWriter()
	// 编码器(日志格式)
	encoder := getEncoder()
	// 自定义core,也就是自定义日志格式  zapcore.DebugLevel 日志级别
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 打印代码行数
	logger = zap.New(core, zap.AddCaller())
	logger = zap.New(core)
	sugarLogger = logger.Sugar()
}

// 编码器(日志格式),也可以自己实现他的结构体
func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// 另一种日志格式
	// return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

// 指定日志将写到哪里去(追加的形式)
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	return zapcore.AddSync(file)
}

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("http://www.google.com")
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
```

#### 3.1、日志切割

**讲getLogWriter函数的代码改成如下即可**

```go
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
```

Lumberjack Logger采用以下属性作为输入:

- Filename: 日志文件的位置
- MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
- MaxBackups：保留旧文件的最大个数
- MaxAges：保留旧文件的最大天数
- Compress：是否压缩/归档旧文件

### 四、最终版本

```go
package main

import (
	"net/http"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 自定义给一个全局变量
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.sogo.com")
	simpleHttpGet("http://www.sogo.com")
}

// Zap库初始化信息
func InitLogger() {
	// 指定日志将写到哪里去
	writeSyncer := getLogWriter()
	// 编码器(日志格式)
	encoder := getEncoder()
	// 自定义core,也就是自定义日志格式  zapcore.DebugLevel 日志级别
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 打印代码行数
	logger := zap.New(core, zap.AddCaller())
	//
	sugarLogger = logger.Sugar()
}

// 编码器(日志格式),也可以自己实现他的结构体
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 指定日志将写到哪里去(追加的形式)
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 测试函数
func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
```



### 五、Gin框架使用Zap库