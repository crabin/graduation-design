package myPkg

import "github.com/flipped-aurora/gin-vue-admin/server/service/myPkg/cmsDiscriminate"

type ServiceGroup struct {
	// Code generated by github.com/flipped-aurora/gin-vue-admin/server Begin; DO NOT EDIT.
	PortScanApi
	ScanIp
	CrawlAllApi
	SqlInject
	IpInfoServer
	XSSFilterServer
	cmsDiscriminate.CMSService
	// Code generated by github.com/flipped-aurora/gin-vue-admin/server End; DO NOT EDIT.
}
