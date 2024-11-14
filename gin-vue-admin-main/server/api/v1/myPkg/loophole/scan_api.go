/*
 * @author Crabin
 */

package myPgk

import (
	"bufio"
	"bytes"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/my"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/file"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/util"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/poc/rule"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type ScanApi struct {
}

type scanSerializer struct {
	// 单个url
	Target  string   `json:"target" binding:"required"`
	Type    string   `json:"type" binding:"required,oneof=multi all"` // multi or all
	VulList []string `json:"vul_list"`
	Remarks string   `json:"remarks"`
}

// @Summary scan url
// @Tags Scan
// @Description 扫描单个url
// @accept json
// @Produce  json
// @Param scan body scanSerializer true "扫描参数"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/scan/url [post]
func (m *ScanApi) Url(c *gin.Context) {
	scan := scanSerializer{}
	err := c.ShouldBindJSON(&scan)
	if err != nil {
		response.FailWithMessage("测试url不可为空，扫描类型为multi或all", c)
		return
	}
	// url去首尾空格
	scan.Target = strings.TrimSpace(scan.Target)
	oreq, err := util.GenOriginalReq(scan.Target)
	if err != nil || oreq == nil {
		response.FailWithMessage("原始请求生成失败", c)
		return
	}

	// 插件列表
	plugins, err := rule.LoadDbPlugin(scan.Type, scan.VulList)
	if err != nil || plugins == nil {
		response.FailWithMessage("poc插件加载失败"+err.Error(), c)
		return
	}
	token := c.Request.Header.Get("x-token")
	j := utils.NewJWT()
	// parseToken 解析token包含的信息
	claims, _ := j.ParseToken(token)

	// 创建任务
	task := my.Task{
		Operator: claims.Username,
		Remarks:  scan.Remarks,
		Target:   scan.Target,
	}
	my.AddTask(&task)
	taskItem := &rule.TaskItem{
		OriginalReq: oreq,
		Plugins:     plugins,
		Task:        &task,
	}

	response.OkWithMessage("任务下发成功", c)
	go rule.TaskProducer(taskItem)
	go rule.TaskConsumer()
	return
}

// @Summary scan raw
// @Tags Scan
// @Description 传文件：请求报文
// @Accept multipart/form-data
// @Param type formData string true "扫描类型：multi / all"
// @Param vul_list formData swaggerArray false "vul_id 列表"
// @Param remarks formData string false "备注"
// @Param target formData file true "file"
// @accept json
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/scan/raw [post]
func (m *ScanApi) Raw(c *gin.Context) {
	scanType := c.PostForm("type")
	vulList := c.PostFormArray("vul_list")
	remarks := c.PostForm("remarks")

	if scanType != "multi" && scanType != "all" {
		response.FailWithMessage("扫描类型为multi或all", c)
		return
	}

	target, err := c.FormFile("target")
	if err != nil {
		response.FailWithMessage("文件上传失败", c)
		return
	}
	// 存文件
	filePath := file.UploadTargetsPath(path.Ext(target.Filename))
	err = c.SaveUploadedFile(target, filePath)

	if err != nil || !file.Exists(filePath) {
		response.FailWithMessage("文件保存失败", c)
		return
	}

	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		response.FailWithMessage("请求报文文件解析失败", c)
		return
	}

	oreq, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	if err != nil || oreq == nil {
		response.FailWithMessage("生成原始请求失败", c)
		return
	}
	if !oreq.URL.IsAbs() {
		scheme := "http"
		oreq.URL.Scheme = scheme
		oreq.URL.Host = oreq.Host
	}

	plugins, err := rule.LoadDbPlugin(scanType, vulList)
	if err != nil || plugins == nil {
		response.FailWithMessage("插件加载失败"+err.Error(), c)
		return
	}

	oReqUrl := oreq.URL.String()

	token := c.Request.Header.Get("Authorization")
	claims, _ := util.ParseToken(token)

	task := my.Task{
		Operator: claims.Username,
		Remarks:  remarks,
		Target:   oReqUrl,
	}
	my.AddTask(&task)
	taskItem := &rule.TaskItem{
		OriginalReq: oreq,
		Plugins:     plugins,
		Task:        &task,
	}

	response.OkWithMessage("任务下发成功", c)
	go rule.TaskProducer(taskItem)
	go rule.TaskConsumer()
	return
}

// @Summary scan list
// @Tags Scan
// @Description 传文件：url列表
// @Accept multipart/form-data
// @Param type formData string true "扫描类型：multi / all"
// @Param vul_list formData swaggerArray false "vul_id 列表"
// @Param remarks formData string false "备注"
// @Param target formData file true "file"
// @Produce  json
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/scan/list [post]
func (m *ScanApi) List(c *gin.Context) {
	scanType := c.PostForm("type")
	vulList := c.PostFormArray("vul_list")
	remarks := c.PostForm("remarks")

	if scanType != "multi" && scanType != "all" {
		response.FailWithMessage("扫描类型为multi或all", c)
		return
	}

	target, err := c.FormFile("target")
	if err != nil {
		response.FailWithMessage("文件上传失败", c)
		return
	}
	// 存文件
	filePath := file.UploadTargetsPath(path.Ext(target.Filename))
	err = c.SaveUploadedFile(target, filePath)

	if err != nil || !file.Exists(filePath) {
		response.FailWithMessage("文件保存失败", c)
		return
	}

	// 加载poc
	plugins, err := rule.LoadDbPlugin(scanType, vulList)
	if err != nil {
		response.FailWithMessage("插件加载失败"+err.Error(), c)
		return
	}
	if len(plugins) == 0 {
		response.FailWithMessage("插件加载失败"+err.Error(), c)
		return
	}
	targets := file.ReadingLines(filePath)

	token := c.Request.Header.Get("Authorization")
	claims, _ := util.ParseToken(token)

	var oReqList []*http.Request
	var taskList []*my.Task

	for _, url := range targets {
		oreq, err := util.GenOriginalReq(url)
		if err != nil {
			continue
		}
		task := my.Task{
			Operator: claims.Username,
			Remarks:  remarks,
			Target:   url,
		}
		my.AddTask(&task)

		oReqList = append(oReqList, oreq)
		taskList = append(taskList, &task)
	}

	if len(oReqList) == 0 || len(taskList) == 0 {
		response.FailWithMessage("url列表加载失败", c)
		return
	}
	response.OkWithMessage("任务下发成功", c)
	for index, oreq := range oReqList {
		taskItem := &rule.TaskItem{
			OriginalReq: oreq,
			Plugins:     plugins,
			Task:        taskList[index],
		}
		go rule.TaskProducer(taskItem)
		go rule.TaskConsumer()
	}
	return
}
