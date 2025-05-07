package patch

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/playwright-community/playwright-go"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InsertIntoDb() {
	err := playwright.Install()
	if err != nil {
		log.Fatalln("Playwright安装失败:", err)
	}
	pl := GetSqlPatch()
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db.AutoMigrate(&Patch{})
	for _, v := range pl {
		result := db.First(&v, "link = ?", v.Link)
		if result.RowsAffected > 0 {
			fmt.Printf("已存在记录：版本：%s | 类型：%-15s | 发布日期：%v | 链接：%s \n", v.Version, v.Name, v.ReleaseDate, v.Link)
		}
		if result.Error != nil {
			fmt.Printf("新创建记录：版本：%s | 类型：%-15s | 发布日期：%v | 链接：%s \n", v.Version, v.Name, v.ReleaseDate, v.Link)
			db.Create(&v)
		}
	}
}
