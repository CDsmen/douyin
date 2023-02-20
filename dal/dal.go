package dal

import (
	"fmt"

	"github.com/RaymondCode/simple-demo/controller"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const MySQLDefaultDSN = "root:12345678@tcp(localhost:3306)/lastproject?charset=utf8&parseTime=True&loc=Local"

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	var testuser controller.User

	DB.First(&testuser)

	fmt.Println(testuser)
}
