package apis

import (
	"UserInsert/dao/mysql"
	"UserInsert/models"
	"fmt"
)

func GetRole() {
	var role = make([]models.Info_Role, 0)

	mysql.DB.Find(&role)
	for _, value := range role {
		fmt.Printf("%s Role_ID: %d \n", value.Name, value.RoleID)
		fmt.Println()
	}
}
