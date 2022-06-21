package models

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
