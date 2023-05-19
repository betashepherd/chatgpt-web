package bootstrap

import (
	"chatgpt-web/config"
	"chatgpt-web/pkg/logger"
	"chatgpt-web/pkg/model"
	"chatgpt-web/pkg/model/user"
	"gorm.io/gorm"
	"os"
)

// SetupDB 启动数据库
func SetupDB() {
	// create db file dir
	dir := "data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
	}

	//建立连接池
	db := model.ConnectDB()

	migration(db)

	//insertAdmin()
}

// migration 迁移
func migration(db *gorm.DB) {
	err := db.AutoMigrate(&user.User{})
	if err != nil {
		logger.Danger("migration model error:", err)
	}
}

func insertAdmin() {
	cf := config.LoadConfig()
	if cf.AuthUser != "" {
		_, err := user.GetByName(cf.AuthUser)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Danger("insert admin error:", err)
		}
		if err == gorm.ErrRecordNotFound {
			_, err = user.CreateUser(cf.AuthUser, cf.AuthPassword, "admin")
			if err != nil {
				logger.Danger("create admin error:", err)
			}
		}
	}
}
