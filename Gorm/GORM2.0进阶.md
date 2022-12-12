## GORM2.0进阶CRUD

### 一、什么是CRUD

CRUD通常指数据库的增删改查操作，本文详细介绍了如何使用GORM实现创建、查询、更新和删除操作。

本文中的`db`变量为`*gorm.DB`对象，例如：

```go
package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
func main() {
	dsn := "root:123456@tcp(192.168.1.208:3306)/demo22?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return
	}
}

```



### 二、创建数据

```go
package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	Id         int       `gorm:"column:id;primaryKey"`
	Name       string    `gorm:"column:name"`
	Province   string    `gorm:"column:province`
	City       string    `gorm:"column:city"`
	Address    string    `gorm:"column:addr"`
	Score      float32   `gorm:"column:score"`
	Enrollment time.Time `gorm:"column:enrollment;type:date"`
}

//使用TableName()来修改默认的表名
func (Student) TableName() string {
	return "student"
}

func main() {
	dsn := "root:123456@tcp(192.168.1.208:3306)/demo22?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败: ", err)
		return
	}

	// 调用函数
	insert(db)

	// 自动迁移表,数据库会自动生成一张UserInfos表,当然也可以去掉s,使用TableName方法
	db.AutoMigrate(&Student{})
}

func insert(db *gorm.DB) {
	// 插入单条记录 <Id是自增主键,传值的时候可以忽略,如果其他的不传应该是数据类型对应的默认值
	student := Student{Name: "光绪", Province: "北京", City: "北京", Score: 80, Enrollment: time.Now()}
	student2 := Student{Name: "慈禧", Province: "北京", Enrollment: time.Now()}
	db.Create(&student)
	db.Create(&student2)

	// 一次性插入多条记录
	students := []Student{
		{Name: "无极", Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()},
		{Name: "小王", Province: "上海", City: "上海", Score: 12, Enrollment: time.Now()},
		{Name: "小亮", Province: "北京", City: "北京", Score: 20, Enrollment: time.Now()},
	}
	db.Create(students)

	// 数据量太大时分批插入
	students = []Student{
		{Name: "大壮", Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()},
		{Name: "刘二", Province: "上海", City: "上海", Score: 12, Enrollment: time.Now()},
		{Name: "文明", Province: "北京", City: "北京", Score: 20, Enrollment: time.Now()},
	}
	db = db.CreateInBatches(students, 2) //一次插入2条
	fmt.Printf("insert %d rows\n", db.RowsAffected)
	fmt.Println("=============insert end=============")
}
```



### 三、查找数据

- 基于创建数据模型操作数据，重复代码不再贴

#### 3.1、一般查询

```go
func search(db *gorm.DB) {
	var student Student
	// 返回一条记录
	db.First(&student)
	fmt.Println(student)

	// 返回多条记录
	var students []Student
	db.Find(&students) // select * from student

	// 随机获取一条记录
	var student1 Student
	db.Take(&student1)

	// 根据主键查询最后一条记录
	var student2 Student
	db.Last(&student2)

	// 查询指定的某条记录(仅当主键为整型时可用)
	var student3 Student
	db.First(&student3, 12) // SELECT * FROM student WHERE id = 12;

}
```

#### 3.2、条件查询

```go
func search(db *gorm.DB) {
	// 条件查询一条记录
	var student Student
	db.Where("name=?", "李四").First(&student) // select name from Student where name = "";

	// 条件查询多条记录
	var students []Student
	db.Where("id > ?", 7).Find(&students) // select * from Student where id > 7;

	// <>
	db.Where("name <> ?", "李四").Find(&students)

	// IN
	db.Where("name IN ?", []string{"张三", "李四"}).Find(&students)

	// LIKE
	db.Where("name LIKE ?", "%张%").Find(&students)

	// AND
	db.Where("name = ? AND score >= ?", "张三", 80).Find(&students)

	// Time
	db.Where("updated_at > ?", time.Now()).Find(&students)

	// BETWEEN
	db.Where("created_at BETWEEN ? AND ?", time.Now(), time.Now()).Find(&students)
	//  SELECT * FROM student WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
}
```

#### 3.3、Struct & Map查询

```go
func search(db *gorm.DB) {
	// 条件查询一条记录
	var student Student
	// 条件查询多条记录
	var students []Student

	// struct 查询
	db.Where(&Student{Name: "张三", Score: 88}).First(&student)
	// SELECT * FROM student WHERE name = "张三" AND score = 88 LIMIT 1;

	// Map查询
	db.Where(map[string]interface{}{"name": "张三", "city": "深圳"}).Find(&students)
	// SELECT * FROM student WHERE name = "张三" AND city = "深圳";

	// 主键的切片
	db.Where([]int{1, 2, 3}).Find(&students)
	// SELECT * FROM student WHERE id IN (1, 2, 3);
}
```

**提示：**当通过结构体进行查询时，**GORM将会只通过非零值字段查询**，这意味着如果你的字段值为`0`，`''`，`false`或者其他`零值`时，将不会被用于构建查询条件，例如：

```go
db.Where(&User{Name: "张三", score: 0}).Find(&students)
//// SELECT * FROM users WHERE name = "jinzhu";
```

- 可以使用指针或实现 Scanner/Valuer 接口来避免这个问题

```go
// 使用指针
type User struct {
  gorm.Model
  Name string
  Age  *int
}

// 使用 Scanner/Valuer
type User struct {
  gorm.Model
  Name string
  Age  sql.NullInt64  // sql.NullInt64 实现了 Scanner/Valuer 接口
}
```

#### 3.4、Not 条件

- 有多个添加的切片就not in，否则就是单纯的not

```go
db.Not("name", "jinzhu").First(&user)
// SELECT * FROM users WHERE name <> "jinzhu" LIMIT 1;

// Not In
db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

// Not In slice of primary keys
db.Not([]int64{1,2,3}).First(&user)
// SELECT * FROM users WHERE id NOT IN (1,2,3);

db.Not([]int64{}).First(&user)
// SELECT * FROM users;

// Plain SQL
db.Not("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE NOT(name = "jinzhu");

// Struct
db.Not(User{Name: "jinzhu"}).First(&user)
// SELECT * FROM users WHERE name <> "jinzhu";
```

#### 3.5、Or条件

```go
db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

// Struct
db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2';

// Map
db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2';
```

#### 3.6、内联条件

- 作用与`Where`查询类似，当内联条件与多个[立即执行方法](https://www.liwenzhou.com/posts/Go/gorm-crud/#autoid-1-3-1)一起使用时, 内联条件不会传递给后面的立即执行方法

```go
// 根据主键获取记录 (只适用于整形主键)
db.First(&user, 23)
//  SELECT * FROM users WHERE id = 23 LIMIT 1;
// 根据主键获取记录, 如果它是一个非整形主键
db.First(&user, "id = ?", "string_primary_key")
//  SELECT * FROM users WHERE id = 'string_primary_key' LIMIT 1;

// Plain SQL
db.Find(&user, "name = ?", "jinzhu")
// SELECT * FROM users WHERE name = "jinzhu";

db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
//  SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;

// Struct
db.Find(&users, User{Age: 20})
//  SELECT * FROM users WHERE age = 20;

// Map
db.Find(&users, map[string]interface{}{"age": 20})
//  SELECT * FROM users WHERE age = 20;
```

#### 3.7、额外查询选项

```go
// 为查询 SQL 添加额外的 SQL 操作
db.Set("gorm:query_option", "FOR UPDATE").First(&user, 10)
// SELECT * FROM users WHERE id = 10 FOR UPDATE;
```

#### 3.8、FirstOrInit

- 获取匹配的第一条记录，否则根据给定的条件初始化一个新的对象 (仅支持 struct 和 map 条件)

```go
// 未找到
db.FirstOrInit(&user, User{Name: "non_existing"})
//  user -> User{Name: "non_existing"}

// 找到
db.Where(User{Name: "Jinzhu"}).FirstOrInit(&user)
//  user -> User{Id: 111, Name: "Jinzhu", Age: 20}

db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})
//  user -> User{Id: 111, Name: "Jinzhu", Age: 20}
```

#### 3.9、Attrs

- 如果记录未找到，将使用参数初始化 struct.

```go
// 未找到
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)
// SELECT * FROM USERS WHERE name = 'non_existing';
// user -> User{Name: "non_existing", Age: 20}

db.Where(User{Name: "non_existing"}).Attrs("age", 20).FirstOrInit(&user)
// SELECT * FROM USERS WHERE name = 'non_existing';
// user -> User{Name: "non_existing", Age: 20}

// 找到
db.Where(User{Name: "Jinzhu"}).Attrs(User{Age: 30}).FirstOrInit(&user)
// SELECT * FROM USERS WHERE name = jinzhu';
// user -> User{Id: 111, Name: "Jinzhu", Age: 20}
```

#### 3.10、Assign

- 不管记录是否找到，都将参数赋值给 struct.

```go
// 未找到
db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrInit(&user)
// user -> User{Name: "non_existing", Age: 20}

// 找到
db.Where(User{Name: "Jinzhu"}).Assign(User{Age: 30}).FirstOrInit(&user)
// SELECT * FROM USERS WHERE name = jinzhu';
// user -> User{Id: 111, Name: "Jinzhu", Age: 30}
```

#### 3.11、高级查询

##### 3.11.1、子查询

- 基于 `*gorm.expr` 的子查询

```go
// SubQuery() 方法貌似gorm2.0噶了
db.Where("id > ?", db.Table("student").Select("id").Where("id = ?", 5).SubQuery()).Find(students)
// SELECT * FROM student WHERE id > (SELECT id from student WHERE id = 5);
```

##### 3.11.2、选择字段

- Select，指定你想从数据库中检索出的字段，默认会选择全部字段

```go
db.Select("name, age").Find(&users)
// SELECT name, age FROM users;

db.Select([]string{"name", "age"}).Find(&users)
// SELECT name, age FROM users;

db.Table("users").Select("COALESCE(age,?)", 42).Rows()
// SELECT COALESCE(age,'42') FROM users;
```

##### 3.11.3、排序

- Order，指定从数据库中检索出记录的顺序。设置第二个参数 reorder 为 `true` ，可以覆盖前面定义的排序条件。

```go
db.Order("age desc, name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

// 多字段排序
db.Order("age desc").Order("name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

// 覆盖排序
db.Order("age desc").Find(&users1).Order("age", true).Find(&users2)
// SELECT * FROM users ORDER BY age desc; (users1)
// SELECT * FROM users ORDER BY age; (users2)
```

##### 3.11.4、数量

- -1 取消 Limit 条件

```go
db.Limit(3).Find(&users)
// SELECT * FROM users LIMIT 3;

// -1 取消 Limit 条件
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
// SELECT * FROM users LIMIT 10; (users1)
// SELECT * FROM users; (users2)
```

##### 3.11.5、偏移

- Offset，指定开始返回记录前要跳过的记录数
- -1 取消 Offset 条件

```go
db.Offset(3).Find(&users)
// SELECT * FROM users OFFSET 3;

// -1 取消 Offset 条件
db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
// SELECT * FROM users OFFSET 10; (users1)
// SELECT * FROM users; (users2)
```

##### 3.11.6、总数

**注意** `Count` 必须是链式查询的最后一个操作 ，因为它会覆盖前面的 `SELECT`，但如果里面使用了 `count` 时不会覆盖

```go
db.Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Find(&users).Count(&count)
// SELECT * from USERS WHERE name = 'jinzhu' OR name = 'jinzhu 2'; (users)
// SELECT count(*) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'; (count)

db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)
// SELECT count(*) FROM users WHERE name = 'jinzhu'; (count)

db.Table("deleted_users").Count(&count)
// SELECT count(*) FROM deleted_users;

db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
// SELECT count( distinct(name) ) FROM deleted_users; (count)
```

##### 3.11.7、Group & Having

```go
rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
for rows.Next() {
  ...
}

// 使用Scan将多条结果扫描进事先准备好的结构体切片中
type Result struct {
	Date time.Time
	Total int
}
var rets []Result
db.Table("users").Select("date(created_at) as date, sum(age) as total").Group("date(created_at)").Scan(&rets)

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
for rows.Next() {
  ...
}

type Result struct {
  Date  time.Time
  Total int64
}
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

##### 3.11.8、Joins

```go
rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
for rows.Next() {
  ...
}

db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

// 多连接及参数
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```

#### 3.12、链式操作相关

##### 3.12.1、链式操作

在调用立即执行方法前不会生成`Query`语句，借助这个特性你可以创建一个函数来处理一些通用逻辑

Method Chaining，Gorm 实现了链式操作接口，所以你可以把代码写成这样：

```go
// 创建一个查询
tx := db.Where("name = ?", "jinzhu")

// 添加更多条件
if someCondition {
  tx = tx.Where("age = ?", 20)
} else {
  tx = tx.Where("age = ?", 30)
}

if yetAnotherCondition {
  tx = tx.Where("active = ?", 1)
}
```

##### 3.12.2、范围Scope

`Scopes`，Scope是建立在链式操作的基础之上的。

基于它，你可以抽取一些通用逻辑，写出更多可重用的函数库。

```go
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
  return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
  return func (db *gorm.DB) *gorm.DB {
    return db.Scopes(AmountGreaterThan1000).Where("status IN (?)", status)
  }
}

db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
// 查找所有金额大于 1000 的信用卡订单

db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)
// 查找所有金额大于 1000 的 COD 订单

db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
// 查找所有金额大于 1000 且已付款或者已发货的订单
```



### 四、更新数据

- 基于创建数据模型操作数据，重复代码不再贴

#### 4.1、更新修改字段

- 如果你只希望更新指定字段，可以使用`Update`或者`Updates`

```go
func update(db *gorm.DB) {
	// 根据where更新一列, 注意是一列不是一行 Update
	db.Model(&Student{}).Where("city = ?", "北京").Update("score", 60)
	// update  student set score = 60 where city = "北京"

	// 更新多列 Updates
	db.Model(&Student{}).Where("city = ?", "北京").Updates(map[string]interface{}{"score": 10, "addr": "中原区"})
	// update  student set score = 60,addr = "中原区" where city = "北京"
}
```

#### 4.2、更新选定字段

- 如果你想更新或忽略某些字段，你可以使用 `Select`，`Omit`

```go
// 更新某些字段
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10'

// 忽略某些字段
db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10'
```

#### 4.3、无Hooks更新

- 上面的更新操作会自动运行 model 的 `BeforeUpdate`, `AfterUpdate` 方法，更新 `UpdatedAt` 时间戳, 在更新时保存其 `Associations`, 如果你不想调用这些方法，你可以使用 `UpdateColumn`， `UpdateColumns`

```go
// 更新单个属性，类似于 `Update`
db.Model(&user).UpdateColumn("name", "hello")
// UPDATE users SET name='hello' WHERE id = 111;

// 更新多个属性，类似于 `Updates`
db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE id = 111;
```

#### 4.4、批量更新

- 批量更新时`Hooks（钩子函数）`不会运行

```go
db.Table("users").Where("id IN (?)", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})
// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);

// 使用 struct 更新时，只会更新非零值字段，若想更新所有字段，请使用map[string]interface{}
db.Model(User{}).Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18;

// 使用 `RowsAffected` 获取更新记录总数
db.Model(User{}).Updates(User{Name: "hello", Age: 18}).RowsAffect
```



### 五、删除数据

- 基于创建数据模型操作数据，重复代码不再贴

#### 5.1、软删除

**警告** 删除记录时，请确保主键字段有值，GORM 会通过主键去删除记录，如果主键为空，GORM 会删除该 model 的所有记录

如果一个 model 有 `DeletedAt` 字段，他将自动获得软删除的功能！ 当调用 `Delete` 方法时， 记录不会真正的从数据库中被删除， 只会将`DeletedAt` 字段的值会被设置为当前时间

```go
func delete(db *gorm.DB) {
	// 使用主键删除
	db = db.Delete(&Student{}, 12)
	// 打印删除的数量
	fmt.Println(db.RowsAffected)

	// 用where删除
	db = db.Where("city in ?", []string{"北京", "深圳"}).Delete(&Student{})
	fmt.Println(db.RowsAffected)

	// 使用主键批量删除
	db = db.Delete(&Student{}, []int{1, 2, 3})
	// 打印删除的数量
	fmt.Println(db.RowsAffected)
}
```

#### 5.1、物理删除

```go
// Unscoped 方法可以物理删除记录
db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;
```



### 六、事务

- https://www.cnblogs.com/infodriven/p/16351565.html

- gorm事务默认是开启的。为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）

```go
// 全局禁用
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})
```

- 把多个操作放一个事务里，可以减少对数据操作的压力

```go
func transaction(db *gorm.DB) {
    // 开启一个事务
	tx := db.Begin()
    
	defer tx.Rollback() // 调用commit后也可以调用rollback，相当于rollback没起任何作用
	for i := 0; i < 10; i++ {
		student := Student{Name: "学生" + strconv.Itoa(i), Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
		if err := tx.Create(&student).Error; err != nil { //注意是tx.Create，不是db.Create
			return //函数返回
		}
	}
    
    // 提交事务
	tx.Commit()
	fmt.Println("=============transaction end=============")
}
------------------------------------------------------------------------------------------------------------
func transaction2(db *gorm.DB) {
	tx := db.Begin()
	defer tx.Rollback()
	for i := 0; i < 10; i++ {
		student := Student{Name: "学生" + strconv.Itoa(i), Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
		if err := tx.Create(&student).Error; err != nil {
			return
		}
		if i == 5 {
            // 模拟失败就回滚
			tx.Rollback() // Rollback之后事务就被关闭了，不能再调用tx.Create()了
		}
	}
	tx.Commit()
	fmt.Println("=============transaction end=============")
}
```

