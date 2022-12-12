package models

// 添加用户完整的数据,后面分别传递给三个接口
type FileUserInfo struct {
	UserName            string   `josn:"userName"`
	RealName            string   `josn:"realName"`
	Serial              string   `josn:"serial"`
	Phone               string   `josn:"phone"`
	RoleIds             []int    `josn:"roleIds"`
	OrgId               int      `josn:"orgId"`
	ImageUrl            string   `josn:"imageUrl"`
	AddTarLibSerials    []string `json:"addTarLibSerials"`
	RemoveTarLibSerials []string `json:"removeTarLibSerials"`
	GroupSerials        []string `json:"groupSerials"`
}

// 用户信息
type UserInfo struct {
	UserName string `josn:"userName"`
	RealName string `josn:"realName"`
	Serial   string `josn:"serial"`
	Phone    string `josn:"phone"`
	RoleIds  []int  `josn:"roleIds"`
	OrgId    int    `josn:"orgId"`
	ImageUrl string `josn:"imageUrl"`
}

// 人像库接口参数
type PortraitLib struct {
	UserID              int64
	RealName            string
	AddTarLibSerials    []string `json:"addTarLibSerials"`
	RemoveTarLibSerials []string `json:"removeTarLibSerials"`
}

// 视频源分组参数
type CameraGroup struct {
	UserID       int64
	GroupSerials []string
}
