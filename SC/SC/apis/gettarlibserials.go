package apis

import (
	"UserInsert/dao/mysql"
	"UserInsert/models"
	"fmt"
)

func GetTarLib() {
	var library = make([]models.Info_Target_Library, 0)

	mysql.DB.Find(&library)

	for _, value := range library {
		fmt.Printf("人像库名字是: %s, 人像库Serial是: %s\n", value.Name, value.Serial)
		fmt.Println()
	}
}
