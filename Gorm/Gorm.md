## Goland之gorm

###  一、gorm常见函数

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
	Addr     string
	Bsd      time.Time
}

const (
	dbUser     = "root"
	dbPassWord = "123456"
	dbHost     = "192.168.1.208"
	dbPort     = 3306
	dbName     = "testgorm"
	dbDriver   = "mysql"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=PRC&parseTime=true",
		dbUser, dbPassWord, dbHost, dbPort, dbName)

	db, err := gorm.Open(dbDriver, dsn)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1) //有错误就退出
	}
    
    // db.LogMode(true) // 开启查询日志
    
	// 数据库自动迁移
	db.AutoMigrate(&User{}) // 根据模型user创建表

	// 判断表是否存在
	fmt.Println(db.HasTable(&User{}))

	// 删除表
	fmt.Println(db.DropTable(&User{}))

	// 创建表，如果表存在会报错
	fmt.Println(db.CreateTable(&User{}))

	// 修改某个字段类型 (少用)
	fmt.Println(db.Model(&User{}).ModifyColumn("bsd", "date"))

	// 删除字段，列
	db.Model(&User{}).DropColumn("bsd")

	// 添加索引
	db.Model(&User{}).AddIndex("name")

	// 设置唯一索引
	db.Model(&User{}).AddUniqueIndex("name")

	// 删除索引
	db.Model(&User{}).RemoveIndex("xxx")
}
```



### 二、gorm设置主键、索引等

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User1 struct {
	Id       int       `gorm:"primary_key"`                                   // 设置主键、int类型会自动增加
	Name     string    `gorm:"type:varchar(32); unique; not null;default:''"` // 设置长度、不能为空、默认''、唯一索引
	Password string    `gorm:"type:varchar(32)"`                              // 注意，设置属性 是key:value形式
	Addr     string    `gorm:"index:index_addr"`                              // 设置索引
	Bsd      time.Time `gorm:"column:bsdy"`                                   // 修改数据库中默认字段名(bsd---->bsdy)
}

// 修改数据库中的表名，默认是结构体的名字
func (u *User1) TableName() string {
	return "user2"
}

const (
	dbUser     = "root"
	dbPassWord = "123456"
	dbHost     = "192.168.1.208"
	dbPort     = 3306
	dbName     = "testgorm"
	dbDriver   = "mysql"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=PRC&parseTime=true",
		dbUser, dbPassWord, dbHost, dbPort, dbName)

	db, err := gorm.Open(dbDriver, dsn)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1) //有错误就退出
	}
	// 数据库自动迁移
	db.AutoMigrate(&User1{}) // 根据模型user创建表
}
```



### 三、gorm之CRUD

#### 3.1、插入数据

```go
// 初始化一个实例，insert到数据库中
	user := User1{
		Name:     "kk3",
		Password: "123456",
		Addr:     "上海",
		Bsd:      time.Date(1998, 11, 11, 0, 0, 0, 0, time.UTC),
	}
	// 判断数据是否存在
	fmt.Println(db.NewRecord(user))
	// 插入数据
	fmt.Println(db.Create(&user))

for i := 20; i < 30; i++ {
		// 初始化一个实例，insert到数据库中
		user := User1{
			Name:     fmt.Sprintf("kk_%d", i),
			Password: fmt.Sprintf("pwd_%d", i),
			Addr:     "上海",
			Bsd:      time.Date(1998, 11, i, 0, 0, 0, 0, time.UTC),
		}
		// 插入数据
		db.Create(&user)
	}
```

