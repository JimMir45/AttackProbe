package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"llm-security-bas/internal/controller"
	"llm-security-bas/internal/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查
	systemCtrl := controller.NewSystemController()
	r.GET("/api/health", systemCtrl.Health)

	// API v1
	api := r.Group("/api/v1")
	{
		// 目标管理
		targetCtrl := controller.NewTargetController(db)
		target := api.Group("/target")
		{
			target.POST("/add", targetCtrl.Add)
			target.POST("/update", targetCtrl.Update)
			target.POST("/delete", targetCtrl.Delete)
			target.POST("/detail", targetCtrl.Detail)
			target.POST("/page", targetCtrl.Page)
			target.POST("/options", targetCtrl.Options)
			target.POST("/test", targetCtrl.Test)
		}

		// 测试用例
		testcaseCtrl := controller.NewTestCaseController(db)
		testcase := api.Group("/testcase")
		{
			testcase.POST("/add", testcaseCtrl.Add)
			testcase.POST("/update", testcaseCtrl.Update)
			testcase.POST("/delete", testcaseCtrl.Delete)
			testcase.POST("/detail", testcaseCtrl.Detail)
			testcase.POST("/page", testcaseCtrl.Page)
			testcase.POST("/stats", testcaseCtrl.Stats)
			testcase.POST("/batch-status", testcaseCtrl.BatchStatus)
		}

		// 任务管理
		taskCtrl := controller.NewTaskController(db)
		task := api.Group("/task")
		{
			task.POST("/add", taskCtrl.Add)
			task.POST("/delete", taskCtrl.Delete)
			task.POST("/detail", taskCtrl.Detail)
			task.POST("/page", taskCtrl.Page)
			task.POST("/start", taskCtrl.Start)
			task.POST("/cancel", taskCtrl.Cancel)
			task.POST("/progress", taskCtrl.Progress)
			task.POST("/results", taskCtrl.Results)
		}

		// 系统管理
		system := api.Group("/system")
		{
			system.POST("/info", systemCtrl.Info)
			system.POST("/config/get", systemCtrl.ConfigGet)
			system.POST("/config/update", systemCtrl.ConfigUpdate)
		}
	}

	// 静态文件服务
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		// 服务静态资源
		r.Static("/assets", filepath.Join(staticDir, "assets"))
		r.StaticFile("/favicon.svg", filepath.Join(staticDir, "favicon.svg"))

		// SPA 路由支持 - 所有非API请求返回index.html
		r.NoRoute(func(c *gin.Context) {
			// 如果是API请求,返回404
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
				return
			}
			c.File(filepath.Join(staticDir, "index.html"))
		})
	}

	return r
}
