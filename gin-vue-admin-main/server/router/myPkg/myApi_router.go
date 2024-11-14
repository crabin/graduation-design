package myPkg

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type MyApi struct {
}

func (m *MyApi) InitStateScanRouter(Router *gin.RouterGroup) {
	MyRouterGroup := Router.Group("my")
	portScanApi := v1.ApiGroupApp.MypkgApiGroup.PortScanApi
	dirbScanApi := v1.ApiGroupApp.MypkgApiGroup.DirbScanApi
	dnsScanApi := v1.ApiGroupApp.MypkgApiGroup.DnsScanApi
	sqlInjectApi := v1.ApiGroupApp.MypkgApiGroup.SqlInjectApi
	ipInfoApi := v1.ApiGroupApp.MypkgApiGroup.IpInfoApi
	xssFilterApi := v1.ApiGroupApp.MypkgApiGroup.XSSFilterApi
	cmsDiscriminateApi := v1.ApiGroupApp.MypkgApiGroup.CmsDiscriminateApi
	pluignsApi := v1.ApiGroupApp.MypkgApiGroup.PluginsApi
	{
		MyRouterGroup.POST("stateScan", portScanApi.StateScan)
		MyRouterGroup.GET("portWs", portScanApi.PortScanWs)
		MyRouterGroup.GET("dirbWs", dirbScanApi.DirbScanHandler)
		MyRouterGroup.GET("DNSWs", dnsScanApi.DnsScanHandler)
		MyRouterGroup.POST("checkSqlInject", sqlInjectApi.CheckSqlInjection)
		MyRouterGroup.POST("getIPInfo", ipInfoApi.GetIPInfo)
		MyRouterGroup.GET("getIpAddress", ipInfoApi.GetIpAddress)
		MyRouterGroup.GET("tracertHostWs", ipInfoApi.TracertHostWs)
		MyRouterGroup.POST("XSSFilter", xssFilterApi.XSSFilterHandle)
		MyRouterGroup.GET("cmsDiscriminate", cmsDiscriminateApi.CmsDiscriminateHandler)
		MyRouterGroup.GET("testDB", pluignsApi.GetPlugin)
	}
}
