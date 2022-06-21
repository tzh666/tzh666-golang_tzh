package models

// 角色信息info_role
type Info_Role struct {
	RoleID int64
	Name   string
	// State  int64
}

// 库信息
type Info_Target_Library struct {
	Serial string
	Name   string
}

// 分组信息
type Info_Video_Resource_Group struct {
	GroupID int64
	Serial  string
	Name    string
}

// Token信息
type TokenUser struct {
	UserName  string
	PassWord  string
	GrantType string
}
