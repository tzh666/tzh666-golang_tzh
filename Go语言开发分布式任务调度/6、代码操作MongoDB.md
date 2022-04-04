### 一、MongoDB常见命令

```sql
-- tab按键可提示补全
> show databases;
> use my_db
> db.help
> db.createCollection("my_collection")
> show collections;
> db.my_collection.help()  -- 查看表操作帮助

> db.my_collection.insertOne({uid:1000,name:"zhangsan"})
{
        "acknowledged" : true,
        "insertedId" : ObjectId("6235f0bbec6a4141df0dc711")
}

> db.my_collection.find()
{ "_id" : ObjectId("6235f0bbec6a4141df0dc711"), "uid" : 1000, "name" : "zhangsan" }

> db.my_collection.createIndex({uid:1})
{
        "createdCollectionAutomatically" : false,
        "numIndexesBefore" : 1,
        "numIndexesAfter" : 2,
        "ok" : 1
}
```



### 二、代码操作

```sh
官网：https://pkg.go.dev/github.com/mongodb/mongo-go-driver

go get go.mongodb.org/mongo-driver/mongo
```

#### 2.1、连接

```go
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database   *mongo.Database
	collection *mongo.Collection
)

func main() {
	// mongodb 连接信息
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.210:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 选择数据库
	database = client.Database("my_db")

	// 选择表
	collection = database.Collection("my_collection")
}
```

#### 2.2、insert

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // 任务名字
	Command   string    `bson:"command"`   // shell 命令
	Err       string    `bson:"err"`       // 脚本错误
	Content   string    `bson:"content"`   // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` // 执行时间点信息
}

var (
	database   *mongo.Database
	collection *mongo.Collection
	record     *LogRecord
	result     *mongo.InsertOneResult
	objectID   primitive.ObjectID
	err        error
)

func main() {
	// mongodb 连接信息
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.210:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to MongoDB!")

	// 选择数据库
	database = client.Database("cron")

	// 选择表
	collection = database.Collection("log")

	//  定义数据类型
	record = &LogRecord{
		"job20",
		"echo hello",
		"",
		"hello",
		TimePoint{
			time.Now().Unix(),
			time.Now().Unix() + 10,
		},
	}

	// 插入数据, 插入的数据是bson格式的
	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		log.Fatal(err)
		return
	}

	// InsertedID是一个全局唯一ID，ObjectID
	// fmt.Printf("%v %T", result.InsertedID, result.InsertedID)
	/*
		> db.log.find()
		{ "_id" : ObjectId("6236e0677323e251655a7304"), "jobName" : "job20", "command" : "echo hello", "err" : "", "content" : "hello", "timePoint" : { "startTime" : NumberLong(1647763559), "endTime" : NumberLong(1647763569) } }
	*/
	objectID = result.InsertedID.(primitive.ObjectID)
	fmt.Println(objectID.Hex())
} 
```

#### 2.2、查询

```go
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // 任务名字
	Command   string    `bson:"command"`   // shell 命令
	Err       string    `bson:"err"`       // 脚本错误
	Content   string    `bson:"content"`   // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` // 执行时间点信息
}

// jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"` // 任务名字
}

var (
	database   *mongo.Database
	collection *mongo.Collection
	cond       *FindByJobName
	record     *LogRecord
	cursor     *mongo.Cursor
	err        error
)

func main() {
	// mongodb 连接信息
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.210:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to MongoDB!")

	// 选择数据库
	database = client.Database("cron")

	// 选择表
	collection = database.Collection("log")

	// 条件查询,按照JobName字段过滤,想找出jobName=job20的,找3条
	cond = &FindByJobName{JobName: "job20"} // 自动被序列化为 {"jobName":"job10"},bson标签

	// 其他查询功能--FindOneAndDelete--FindOneAndReplace--FindOneAndUpdate
	// cursor是结果集
	if cursor, err = collection.Find(context.TODO(), cond); err != nil {
		log.Fatal(err)
		return
	}
	// 遍历结果集
	for cursor.Next(context.TODO()) {
		// 定义日志对象
		record = &LogRecord{}

		// 序列化bson对象
		if err = cursor.Decode(record); err != nil {
			log.Fatal(err)
			return
		}
		// 打印数据
		fmt.Println(record)
	}

	// 释放游标
	defer cursor.Close(context.TODO())
}
```

#### 2.3、删除

```GO
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // 任务名字
	Command   string    `bson:"command"`   // shell 命令
	Err       string    `bson:"err"`       // 脚本错误
	Content   string    `bson:"content"`   // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` // 执行时间点信息
}

// jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"` // 任务名字
}

// startTime小于某时间,这是删除的条件
// 这结构体要序列为 {"$It":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$It"`
}

// {"TimePoint.startTime":{"$It":timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"TimePoint.startTime"`
}

var (
	database   *mongo.Database
	collection *mongo.Collection
	cond       *FindByJobName
	record     *LogRecord
	cursor     *mongo.Cursor
	deleteCond *DeleteCond
	delResult  *mongo.DeleteResult
	err        error
)

func main() {
	// mongodb 连接信息
	clientOptions := options.Client().ApplyURI("mongodb://192.168.1.210:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to MongoDB!")

	// 选择数据库
	database = client.Database("cron")

	// 选择表
	collection = database.Collection("log")

	// 删除开始时间早于当前时间的所有日志; {"$It":当前时间}
	// 删除条件deleteCond
	deleteCond = &DeleteCond{
		beforeCond: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}

	fmt.Printf("删除条件是: %v, %T\n", deleteCond, deleteCond)

	// 执行删除
	if delResult, err = collection.DeleteMany(context.TODO(), deleteCond); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("删除的行数: %v\n", delResult.DeletedCount)
}
```

#### 2.4、update

```go

```

