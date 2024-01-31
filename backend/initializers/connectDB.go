package initializers

import (
	"fmt"

	"github.com/wpcodevo/google-github-oath2-golang/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	// dsn := "root:mypassword@tcp(db.fabkli.ch:80)/testdb"
	dsn := "root:secretpass@tcp(db.fabkli.ch:38487)/oauth?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("2")

	// mysqlConn := mysql.Open(dsn)
	fmt.Println("3")
	// DB, err := gorm.Open(mysqlConn, &gorm.Config{})

	// -----------
	var datetimePrecision = 2
	DB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,                // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}), &gorm.Config{})

	fmt.Println("4")
	// DB, err = gorm.Open(sqlite.Open("golang.db"), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to the Database")
		fmt.Println("error", err)
	}
	fmt.Println("5")

	// DB.Model(&models.User{}).
	DB.AutoMigrate(&models.User{})
	fmt.Println("🚀 Connected Successfully to the Database")
}
