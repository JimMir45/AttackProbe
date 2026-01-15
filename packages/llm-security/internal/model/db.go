package model

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	err = db.AutoMigrate(
		&Target{},
		&TestCase{},
		&Task{},
		&TaskTestCase{},
		&TaskResult{},
		&SysConfig{},
	)
	if err != nil {
		return err
	}

	// 初始化默认配置
	initDefaultConfig(db)

	// 初始化内置测试用例
	count := initBuiltinTestCases(db)
	if count > 0 {
		log.Printf("Initialized %d builtin test cases", count)
	}

	DB = db
	log.Println("Database initialized successfully")
	return nil
}

func initDefaultConfig(db *gorm.DB) {
	configs := []SysConfig{
		{ConfigKey: "executor.concurrency", ConfigValue: "5", Description: "并发执行数"},
		{ConfigKey: "executor.timeout", ConfigValue: "30000", Description: "单用例超时(ms)"},
		{ConfigKey: "system.version", ConfigValue: "1.0.0", Description: "系统版本"},
	}

	for _, cfg := range configs {
		db.Where("config_key = ?", cfg.ConfigKey).FirstOrCreate(&cfg)
	}
}

func GetDB() *gorm.DB {
	return DB
}

func initBuiltinTestCases(db *gorm.DB) int {
	count := 0
	for _, tc := range BuiltinTestCases {
		var existing TestCase
		err := db.Where("name = ? AND is_builtin = 1", tc.Name).First(&existing).Error
		if err != nil {
			newTC := tc
			newTC.IsBuiltin = 1
			newTC.Status = 1
			if db.Create(&newTC).Error == nil {
				count++
			}
		}
	}
	return count
}
