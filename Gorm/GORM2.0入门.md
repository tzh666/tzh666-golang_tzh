## GORM2.0从入门到精通

### 一、什么是GORM

- GORM是一个全能的、友好的、基于 golan的ORM库
- GORM倾向于约定,而不是配置。默认情况下,GORM使用D作为主键,使用结构体名的【蛇形复数】作为表名,字段名的【蛇形】作为列名,并使用 Createdat、 Updated字段追踪创建、更新时间



### 二、GORM2.0发布说明

查看官网：https://gorm.io/zh_CN/docs/v2_release_note.html

GORM1.X 参考：https://www.liwenzhou.com/posts/Go/gorm/  +  https://www.liwenzhou.com/posts/Go/gorm-crud/



### 三、GORM2.0安装

- GORM2.0 对 `go` 版本最低要求升级到1.16

```sh
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u gorm.io/driver/sqlite

// **注意** GORM `v2.0.0` 发布的 git tag 是 `v1.20.0
```



### 四、连接数据库

中文官网参考：https://gorm.io/zh_CN/docs/connecting_to_the_database.html#MySQL

#### 4.1、连接MySQL

- 想要正确的处理 `time.Time`，需要带上 `parseTime` 参数
- 要支持完整的 UTF-8 编码，需要将 `charset=utf8` 更改为 `charset=utf8mb4`

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

MySQL 驱动程序提供了 [一些高级配置](https://github.com/go-gorm/mysql) 可以在初始化过程中使用，例如：

```go
db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
  DefaultStringSize: 256,            // string 类型字段的默认长度
  DisableDatetimePrecision: true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
  DontSupportRenameIndex: true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
  DontSupportRenameColumn: true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
  SkipInitializeWithVersion: false,  // 根据当前 MySQL 版本自动配置
}), &gorm.Config{})
```

#### 4.2、连接PostgreSQL

```go
import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

#### 4.3、连接SQLite

```go
import (
  "gorm.io/driver/sqlite" // 基于 GGO 的 Sqlite 驱动
  // "github.com/glebarez/sqlite" // 纯 Go 实现的 SQLite 驱动, 详情参考： https://github.com/glebarez/sqlite
  "gorm.io/gorm"
)

// github.com/mattn/go-sqlite3
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
```

#### 4.4、连接Clickhouse

```go
import (
  "gorm.io/driver/clickhouse"
  "gorm.io/gorm"
)

func main() {
  dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
  db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})

  // Auto Migrate
  db.AutoMigrate(&User{})
  // Set table options
  db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&User{})

  // 插入
  db.Create(&user)

  // 查询
  db.Find(&user, "id = ?", 10)

  // 批量插入
  var users = []User{user1, user2, user3}
  db.Create(&users)
  // ...
}
```

#### 4.5、连接池

- GORM 使用 [database/sql](https://pkg.go.dev/database/sql) 维护连接池

```go
// 获取通用数据库对象 sql.DB，然后使用其提供的功能
sqlDB, err := db.DB()

// SetMaxIdleConns 设置空闲连接池中连接的最大数量
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns 设置打开数据库连接的最大数量。
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime 设置了连接可复用的最大时间。
sqlDB.SetConnMaxLifetime(time.Hour)

// Ping
sqlDB.Ping()

// Close
sqlDB.Close()

// 返回数据库统计信息
sqlDB.Stats()
```

#### 4.6、不支持的数据库

有些数据库可能兼容 `mysql`、`postgres` 的方言，在这种情况下，你可以直接使用这些数据库的方言。



### 五、GORM2.0操作MySQL基本示例

- 使用GORM连接MySQL

```go
// UserInfo 用户信息
type User struct {
	ID     int
	Name   string
	Gender string
	Hobby  string
}

//使用TableName()来修改默认的表名
func (User) TableName() string {
	return "user"
}

func main() {
	dsn := "root:123456@tcp(192.168.1.208:3306)/demo22?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return
	}

	// 自动迁移表,数据库会自动生成一张UserInfos表,当然也可以去掉s,使用TableName方法
	db.AutoMigrate(&User{})

	fmt.Println("连接数据库成功")
}
```

- 使用GORM连接MySQL进行创建、查询、更新、删除操作
- // 自动迁移表,数据库会自动生成一张UserInfos表,当然也可以去掉s,使用TableName方法
  	**db.AutoMigrate(&User{})**

```go
package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// UserInfo 用户信息
type User struct {
	ID     int
	Name   string
	Gender string
	Hobby  string
}

//使用TableName()来修改默认的表名
func (User) TableName() string {
	return "user"
}

func main() {
	dsn := "root:123456@tcp(192.168.1.208:3306)/demo22?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return
	}

	// 自动迁移表,数据库会自动生成一张UserInfos表,当然也可以去掉s,使用TableName方法
	db.AutoMigrate(&User{})

	// 新增一条数据
	db.Create(&User{
		ID:     4,
		Name:   "沉默",
		Gender: "女",
		Hobby:  "跳舞",
	})

	// 查找
	var user User
	db.First(&user, 4) // 根据整形主键查找, 等价于 db.Where("ID = ?", 1).First(&user)

	fmt.Println(user)

	// update 数据, 把名字是沉默的hobby字段值更改为爬山
	db.Model(&user).Update("hobby", "爬山").Where("Name = ?", "沉默")

	// 删除数据
	db.Delete(&user, 4)
}
```



### 六、GORM2.0模型定义

模型是标准的 struct，由 Go 的基本数据类型、实现了 [Scanner](https://pkg.go.dev/database/sql/?tab=doc#Scanner) 和 [Valuer](https://pkg.go.dev/database/sql/driver#Valuer) 接口的自定义类型及其指针或别名组成

例如：

```go
type User struct {
  ID           uint
  Name         string
  Email        *string
  Age          uint8
  Birthday     *time.Time
  MemberNumber sql.NullString
  ActivatedAt  sql.NullTime
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
```

#### 6.1、约定

GORM 倾向于约定优于配置 默认情况下，GORM 使用 `ID` 作为主键，使用结构体名的 `蛇形复数` 作为表名，字段名的 `蛇形` 作为列名，并使用 `CreatedAt`、`UpdatedAt` 字段追踪创建、更新时间

如果您遵循 GORM 的约定，您就可以少写的配置、代码。 如果约定不符合您的实际要求，[GORM 允许你配置它们](https://gorm.io/zh_CN/docs/conventions.html)

#### 6.2、gorm.Model

GORM 定义一个 `gorm.Model` 结构体，其包括字段 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt`

```go
// gorm.Model 的定义
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

您可以将它嵌入到您的结构体中，以包含这几个字段，详情请参考 [嵌入结构体](https://gorm.io/zh_CN/docs/models.html#embedded_struct)

#### 6.3、高级选项

##### 6.3.1、字段级权限控制

- 可导出的字段在使用 GORM 进行 CRUD 时拥有全部的权限，此外，GORM 允许您用标签控制字段级别的权限。
- 这样您就可以让一个字段的权限是只读、只写、只创建、只更新或者被忽略

**注意：** 使用 GORM Migrator 创建表时，不会创建被忽略的字段

```go
type User struct {
  Name string `gorm:"<-:create"` // 允许读和创建
  Name string `gorm:"<-:update"` // 允许读和更新
  Name string `gorm:"<-"`        // 允许读和写（创建和更新）
  Name string `gorm:"<-:false"`  // 允许读，禁止写
  Name string `gorm:"->"`        // 只读（除非有自定义配置，否则禁止写）
  Name string `gorm:"->;<-:create"` // 允许读和写
  Name string `gorm:"->:false;<-:create"`  // 仅创建（禁止从 db 读）
  Name string `gorm:"-"`                   // 通过 struct 读写会忽略该字段
  Name string `gorm:"-:all"`               // 通过 struct 读写、迁移会忽略该字段
  Name string `gorm:"-:migration"`         // 通过 struct 迁移会忽略该字段
}
```

##### 6.3.2、创建/更新时间追踪（纳秒、毫秒、秒、Time）

GORM 约定使用 `CreatedAt`、`UpdatedAt`、UpdatedAt 追踪创建/更新时间/更新记录的时间。如果您定义了这种字段，GORM 在创建、更新时会自动填充 [当前时间](https://gorm.io/zh_CN/docs/gorm_config.html#now_func)

要使用不同名称的字段，您可以配置 `autoCreateTime`、`autoUpdateTime` 标签

如果您想要保存 UNIX（毫/纳）秒时间戳，而不是 time，您只需简单地将 `time.Time` 修改为 `int` 即可

```go
type User struct {
  CreatedAt time.Time // 在创建时，如果该字段值为零值，则使用当前时间填充
  UpdatedAt int       // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
  Updated   int64 `gorm:"autoUpdateTime:nano"`  // 使用时间戳填纳秒数充更新时间
  Updated   int64 `gorm:"autoUpdateTime:milli"` // 使用时间戳毫秒数填充更新时间
  Created   int64 `gorm:"autoCreateTime"`       // 使用时间戳秒数填充创建时间
}
```

##### 6.3.3、字段标签

声明 model 时，tag 是可选的，GORM 支持以下 tag： tag 名大小写不敏感，但建议使用 `camelCase` 风格

| 标签名                 | 说明                                                         |
| :--------------------- | :----------------------------------------------------------- |
| column                 | 指定 db 列名                                                 |
| type                   | 列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：`not null`、`size`, `autoIncrement`… 像 `varbinary(8)` 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：`MEDIUMINT UNSIGNED not NULL AUTO_INCREMENT` |
| serializer             | 指定将数据序列化或反序列化到数据库中的序列化器, 例如: `serializer:json/gob/unixtime` |
| size                   | 定义列数据类型的大小或长度，例如 `size: 256`                 |
| primaryKey             | 将列定义为主键                                               |
| unique                 | 将列定义为唯一键                                             |
| default                | 定义列的默认值                                               |
| precision              | 指定列的精度                                                 |
| scale                  | 指定列大小                                                   |
| not null               | 指定列为 NOT NULL                                            |
| autoIncrement          | 指定列为自动增长                                             |
| autoIncrementIncrement | 自动步长，控制连续记录之间的间隔                             |
| embedded               | 嵌套字段                                                     |
| embeddedPrefix         | 嵌入字段的列名前缀                                           |
| autoCreateTime         | 创建时追踪当前时间，对于 `int` 字段，它会追踪时间戳秒数，您可以使用 `nano`/`milli` 来追踪纳秒、毫秒时间戳，例如：`autoCreateTime:nano` |
| autoUpdateTime         | 创建/更新时追踪当前时间，对于 `int` 字段，它会追踪时间戳秒数，您可以使用 `nano`/`milli` 来追踪纳秒、毫秒时间戳，例如：`autoUpdateTime:milli` |
| index                  | 根据参数创建索引，多个字段使用相同的名称则创建复合索引，查看 [索引](https://gorm.io/zh_CN/docs/indexes.html) 获取详情 |
| uniqueIndex            | 与 `index` 相同，但创建的是唯一索引                          |
| check                  | 创建检查约束，例如 `check:age > 13`，查看 [约束](https://gorm.io/zh_CN/docs/constraints.html) 获取详情 |
| <-                     | 设置字段写入的权限， `<-:create` 只创建、`<-:update` 只更新、`<-:false` 无写入权限、`<-` 创建和更新权限 |
| ->                     | 设置字段读的权限，`->:false` 无读权限                        |
| -                      | 忽略该字段，`-` 表示无读写，`-:migration` 表示无迁移权限，`-:all` 表示无读写迁移权限 |
| comment                | 迁移时为字段添加注释                                         |

##### 6.3.4、关联标签

GORM 允许通过标签为关联配置外键、约束、many2many 表，详情请参考 [关联部分](https://gorm.io/zh_CN/docs/associations.html#tags)