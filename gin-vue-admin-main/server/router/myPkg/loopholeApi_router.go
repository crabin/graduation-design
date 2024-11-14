/*
 * @author Crabin
 */

package myPkg

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type LoopholeApi struct {
}

func (m *LoopholeApi) InitLoopholeApiRouter(Router *gin.RouterGroup) {
	MyRouterGroup := Router.Group("my")
	initApi := v1.ApiGroupApp.MypkgApiGroup.InitApi
	MyRouterGroup.GET("/init", initApi.InitVulData)

	{
		pluginRoutes := MyRouterGroup.Group("poc")
		pluginsApi := v1.ApiGroupApp.MypkgApiGroup.PluginsApi
		{
			// all
			pluginRoutes.GET("/", pluginsApi.Get)
			// 增
			pluginRoutes.POST("/", pluginsApi.Add)
			// 改
			pluginRoutes.PUT("/:id/", pluginsApi.Update)
			// 详情
			pluginRoutes.GET("/:id/", pluginsApi.Detail)
			// 删
			pluginRoutes.DELETE("/:id/", pluginsApi.Delete)
			// 测试单个poc
			pluginRoutes.POST("/run/", pluginsApi.Run)
			// 上传yaml文件
			pluginRoutes.POST("/upload/", pluginsApi.UploadYaml)
			// 下载yaml文件
			pluginRoutes.POST("/download/", pluginsApi.DownloadYaml)
		}

		vulRoutes := MyRouterGroup.Group("/vul")
		vulnerabilityApi := v1.ApiGroupApp.MypkgApiGroup.VulnerabilityApi
		{
			// basic
			vulRoutes.GET("/basic/", vulnerabilityApi.Basic)
			// all
			vulRoutes.GET("/", vulnerabilityApi.Get)
			// 增
			vulRoutes.POST("/", vulnerabilityApi.Add)
			// 改
			vulRoutes.PUT("/:id/", vulnerabilityApi.Update)
			// 详情
			vulRoutes.GET("/:id/", vulnerabilityApi.Detail)
			// 删
			vulRoutes.DELETE("/:id/", vulnerabilityApi.Delete)
		}

		appRoutes := MyRouterGroup.Group("/product")
		webAppApi := v1.ApiGroupApp.MypkgApiGroup.WebAppApi
		{
			// all
			appRoutes.GET("/", webAppApi.Get)
			// 增
			appRoutes.POST("/", webAppApi.Add)
			// 改
			appRoutes.PUT("/:id/", webAppApi.Update)
			// 详情
			appRoutes.GET("/:id/", webAppApi.Detail)
			// 删
			appRoutes.DELETE("/:id/", webAppApi.Delete)
		}

		scanRoutes := MyRouterGroup.Group("/scan")
		scanApi := v1.ApiGroupApp.MypkgApiGroup.ScanApi
		{
			scanRoutes.POST("/url/", scanApi.Url)
			scanRoutes.POST("/raw/", scanApi.Raw)
			scanRoutes.POST("/list/", scanApi.List)
		}

		taskRoutes := MyRouterGroup.Group("/task")
		taskApi := v1.ApiGroupApp.MypkgApiGroup.TaskApi
		{
			// all
			taskRoutes.GET("/", taskApi.Get)
			// 删
			taskRoutes.DELETE("/:id/", taskApi.Delete)
		}

		resultRoutes := MyRouterGroup.Group("/result")
		resultApi := v1.ApiGroupApp.MypkgApiGroup.ResultApi
		{
			// all
			resultRoutes.GET("/", resultApi.Get)
			// 删
			resultRoutes.DELETE("/:id/", resultApi.Delete)
		}
	}
	{
		serverRoutes := MyRouterGroup.Group("/server")
		serverApi := v1.ApiGroupApp.MypkgApiGroup.ScanServerApi
		{
			serverRoutes.POST("/", serverApi.ScanServer)
		}
	}
}
