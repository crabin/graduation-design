package myPkg

import (
	myPgk "github.com/flipped-aurora/gin-vue-admin/server/api/v1/myPkg/loophole"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
)

type ApiGroup struct {
	// Code generated by github.com/flipped-aurora/gin-vue-admin/server Begin; DO NOT EDIT.
	PortScanApi
	DirbScanApi
	DnsScanApi
	SqlInjectApi
	IpInfoApi
	XSSFilterApi
	CmsDiscriminateApi
	myPgk.PluginsApi
	myPgk.VulnerabilityApi
	myPgk.WebAppApi
	myPgk.TaskApi
	myPgk.ResultApi
	myPgk.ScanApi
	myPgk.InitApi
	myPgk.ScanServerApi
	// Code generated by github.com/flipped-aurora/gin-vue-admin/server End; DO NOT EDIT.
}

var (
	portScanService       = service.ServiceGroupApp.MypkgServiceGroup.PortScanApi
	portIp                = service.ServiceGroupApp.MypkgServiceGroup.ScanIp
	crawlAllApi           = service.ServiceGroupApp.MypkgServiceGroup.CrawlAllApi
	sqlInject             = service.ServiceGroupApp.MypkgServiceGroup.SqlInject
	ipInfoServer          = service.ServiceGroupApp.MypkgServiceGroup.IpInfoServer
	xssFilterServer       = service.ServiceGroupApp.MypkgServiceGroup.XSSFilterServer
	cmsDiscriminateserver = service.ServiceGroupApp.MypkgServiceGroup.CMSService
)
