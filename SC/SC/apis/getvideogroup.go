package apis

import (
	"UserInsert/dao/mysql"
	"UserInsert/models"
	"fmt"
)

func GetVGroup() {
	var vGroup = make([]models.Info_Video_Resource_Group, 0)

	mysql.DB.Find(&vGroup)

	for _, value := range vGroup {
		fmt.Printf("GroupID is: %d, GroupName is: %s, GroupSerial is: %s\n", value.GroupID, value.Name, value.Serial)
		fmt.Println()
	}
}
