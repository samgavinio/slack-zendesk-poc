package operations

import (
	"github.com/zendesk/slack-poc/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"

)

var DB *gorm.DB

func InitDatabase() {
	cfg := config.GetConfig()
	DB = connect(cfg)
}

func connect(cfg config.Config) *gorm.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DatabaseUsername,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName)
	db, dbError := gorm.Open("mysql", connectionString)
	db.LogMode(true)
	db.SingularTable(true)

	if dbError != nil {
		panic(dbError)
	}

	return db
}
