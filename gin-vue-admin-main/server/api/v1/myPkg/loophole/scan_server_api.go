/*
 * @author Crabin
 */

package myPgk

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type ScanServerApi struct {
}

type scanServerSerializer struct {
	// 单个url
	Target  string   `json:"target" binding:"required"`
	Type    string   `json:"type" binding:"required,oneof=multi all"` // multi or all
	VulList []string `json:"vul_list"`
	Remarks string   `json:"remarks"`
}

func (m *ScanServerApi) ScanServer(c *gin.Context) {
	scanServer := scanServerSerializer{}

	err := c.ShouldBindJSON(&scanServer)
	if err != nil {
		response.FailWithMessage("测试url不可为空，扫描类型为multi或all", c)
		return
	}
	// url去首尾空格
	scanServer.Target = strings.TrimSpace(scanServer.Target)
	// 扫描目标
	//task := kunpeng.Task{
	//	Type:   "service",
	//	Netloc: scanServer.Target,
	//	Target: "mysql",
	//	Meta: kunpeng.Meta{
	//		System:   "",
	//		PathList: []string{},
	//		FileList: []string{},
	//		PassList: []string{},
	//	},
	//}
	//fmt.Println(task)
	//jsonBytes, _ := json.Marshal(task)
	//result := kunpeng.Greeter.Check(string(jsonBytes))
	//fmt.Println(result)
	//
	//response.OkWithData(result, c)
}
