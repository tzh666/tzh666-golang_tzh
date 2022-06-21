package models

// type TokenUser struct {
// 	UserName  string `json:"username"`
// 	PassWord  string `json:"password"`
// 	GrantType string `json:"grant_type"`
// }

// 获取接口返回的人像库列表数据
type PageX struct {
	HasNext  bool  `json:"hasNext"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}

type ListX struct {
	TargetName    string      `json:"targetName"`    // 小图名字
	Key           string      `json:"key"`           // 身份ID
	Gender        string      `json:"gender"`        // 性别
	CoverImageUrl string      `json:"coverImageUrl"` // 人脸小图
	ImageList     interface{} `json:"imageUrl"`      // 人像图片列表
}

type DataX struct {
	List []ListX `json:"list"`
	Page PageX   `json:"page"`
}

type TortraitList struct {
	Data      DataX  `json:"data"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Success   bool   `json:"success"`
}

type Portraits struct {
	Page          int64    `json:"page"`
	SenseType     int64    `json:"sensetype"`
	SageSize      int64    `json:"pagesize"`
	TarLibSerials []string `json:"tarlibserials"`
}
