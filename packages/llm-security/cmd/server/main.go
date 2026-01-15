package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"llm-security-bas/internal/model"
	"llm-security-bas/internal/router"
	"llm-security-bas/internal/service/executor"

	"gorm.io/gorm"
)

var (
	port   = flag.Int("port", 8080, "服务端口")
	dbPath = flag.String("db", "./data/llm-security.db", "数据库路径")
)

func main() {
	flag.Parse()

	// 打印启动信息
	fmt.Println("========================================")
	fmt.Println("  LLM Security BAS v1.0.0")
	fmt.Println("  大模型安全有效性验证平台")
	fmt.Println("========================================")

	// 初始化数据库
	log.Printf("初始化数据库: %s", *dbPath)
	if err := model.InitDB(*dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化执行器
	initExecutor(model.GetDB())

	// 设置路由
	r := router.SetupRouter(model.GetDB())

	// 打印路由信息
	log.Println("注册的API路由:")
	for _, route := range r.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	// 启动服务
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("服务启动: http://localhost%s", addr)
	log.Printf("健康检查: http://localhost%s/api/health", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// initExecutor 初始化任务执行器
func initExecutor(db *gorm.DB) {
	// 从配置表读取执行器参数
	concurrency := 5
	timeout := 30 * time.Second

	var cfgConcurrency model.SysConfig
	if db.Where("config_key = ?", "executor.concurrency").First(&cfgConcurrency).Error == nil {
		if v, err := strconv.Atoi(cfgConcurrency.ConfigValue); err == nil && v > 0 {
			concurrency = v
		}
	}

	var cfgTimeout model.SysConfig
	if db.Where("config_key = ?", "executor.timeout").First(&cfgTimeout).Error == nil {
		if v, err := strconv.Atoi(cfgTimeout.ConfigValue); err == nil && v > 0 {
			timeout = time.Duration(v) * time.Millisecond
		}
	}

	executor.Init(db, concurrency, timeout)
	log.Printf("执行器初始化完成: 并发数=%d, 超时=%v", concurrency, timeout)
}
