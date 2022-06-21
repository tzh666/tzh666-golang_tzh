package utils

import (
	"strconv"
	"strings"
)

// 从返回值中获取 UserID
func GetUserID(str string) (userID int64, err error) {

	// 获取字符串类型的UserID
	strUserID := strings.Split(strings.Split(str, ":")[1], ",")[0]

	if strUserID == "null" || strUserID == "false" {
		return -1, err
	}
	// 返回int64类型
	return strconv.ParseInt(strUserID, 10, 64)
}
