### 分布式ID生成器

#### 分布式ID的特点

```sh
全局唯一性: 不能出现有重复的ID标识,这是基本要求。
递增性: 确保生成ID对于用户或业务是递增的。
高可用性: 确保任何时候都能生成正确的ID。高性能性:在高幷发的环境下依然表现良好。
```

```sh
不仅仅是用于用户D,实际互联网中有很多场景需要能够生成类似MySαL自增D这样不断增大,同时又不会重复的d。以支持业务中的高并发场景。
```

```sh
比较典型的场景有:电商促销时短时间内会有大量的订单涌入到系统,比如每秒10w+;明星出轨时微博短时间內会产生大量的相关微博转发和评论消息。在这些业务场景下将数据插入数据库之前,我们需要给这些订单和消息先分配一个唯一1D,然后再保存到数据库中。对这个id的要求是希望其中能带有些时间信息,这样即使我们后端的系统对消息进行了分库分表,也能够以时间顺序对这些消息进行排序
```

![image-20220131151921469](D:\Golang常用库记录\基于雪花算法生成用户ID\基于雪花算法生成用户ID\image-20220131151921469.png)

![image-20220131152127790](D:\Golang常用库记录\基于雪花算法生成用户ID\基于雪花算法生成用户ID\image-20220131152127790.png)

![image-20220131152147454](D:\Golang常用库记录\基于雪花算法生成用户ID\基于雪花算法生成用户ID\image-20220131152147454.png)

```go
package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *sf.Node

func Init(startTime string, machindeID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		zap.L().Error("snowflake Error: ", zap.Error(err))
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machindeID)

	return
}

func GenID() int64 {
	// 还有其他的返回值类型
	return node.Generate().Int64()
}

func SGenID() string {
	// 还有其他的返回值类型
	return node.Generate().String()
}


package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowflake(t *testing.T) {
	err := Init("2022-01-31", 1)
	if err != nil {
		fmt.Println(err)
	}
	id := GenID()
	fmt.Println(id)
}

```

