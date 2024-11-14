/*
 * @author Crabin
 */

package myPgk

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/my"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/util"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/poc/rule"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"gopkg.in/yaml.v2"
	"gorm.io/datatypes"
	"io/ioutil"
	"log"
	"strconv"
)

type PluginsApi struct {
}

type Serializer struct {
	// 返回给前端的字段
	DespName    string         `json:"desp_name"`
	Id          int            `json:"id"`
	VulId       string         `json:"vul_id"`
	Affects     string         `json:"affects"`
	JsonPoc     datatypes.JSON `json:"json_poc"`
	Enable      bool           `json:"enable"`
	Description int            `json:"description"`
}

type RunSerializer struct {
	// 运行单个
	Target  string         `json:"target" binding:"required"`
	VulId   string         `json:"vul_id"`
	Affects string         `json:"affects" binding:"required"`
	JsonPoc datatypes.JSON `json:"json_poc"`
}

type DownloadSerializer struct {
	// 下载 yaml
	JsonPoc datatypes.JSON `json:"json_poc"`
}

func (m *PluginsApi) GetPlugin(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	plugin := my.GetPlugin(id)
	response.OkWithData(plugin, c)
	return
}

// @Summary plugin detail
// @Tags Plugin
// @Description 详情
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/{id}/ [get]
func (m *PluginsApi) Detail(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var data interface{}
	if id < 0 {
		response.FailWithMessage("ID必须大于0", c)
		return
	}
	if my.ExistPluginByID(id) {
		data = my.GetPlugin(id)
		response.OkWithData(data, c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}
}

// @Summary plugin list
// @Tags Plugin
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.PluginSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/ [get]
func (m *PluginsApi) Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := my.PluginSearchField{Search: "", EnableField: -1, AffectsField: ""}

	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}
	if arg := c.Query("enableField"); arg != "" {
		enable, _ := strconv.ParseBool(arg)
		if enable {
			field.EnableField = 1
		} else {
			field.EnableField = 0
		}
	}
	if arg := c.Query("affectsField"); arg != "" {
		field.AffectsField = arg
	}
	plugins := my.GetPlugins(page, pageSize, &field)

	var pluginRespData []Serializer

	for _, plugin := range plugins {
		var despName string
		if plugin.Vulnerability != nil {
			despName = plugin.Vulnerability.NameZh
		} else {
			despName = ""
		}

		pluginRespData = append(pluginRespData, Serializer{
			DespName:    despName,
			Id:          plugin.Id,
			VulId:       plugin.VulId,
			Affects:     plugin.Affects,
			JsonPoc:     plugin.JsonPoc,
			Enable:      plugin.Enable,
			Description: plugin.Desc,
		})
	}
	data["data"] = pluginRespData
	total := my.GetPluginsTotal(&field)
	data["total"] = total
	response.OkWithData(data, c)
	return
}

// @Summary plugin add
// @Tags Plugin
// @Description 新增
// @Produce  json
// @Security token
// @Param plugin body rule.Plugin true "plugin"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/ [post]
func (m *PluginsApi) Add(c *gin.Context) {
	plugin := my.Plugin{}
	err := c.ShouldBindJSON(&plugin)
	if err != nil {
		response.FailWithMessage("参数不合法", c)
		return
	}
	// 漏洞编号自动生成
	vulId, err := my.GenPluginVulId()
	if err != nil {
		response.FailWithMessage("漏洞编号生成失败", c)
		return
	}
	plugin.VulId = vulId
	my.AddPlugin(plugin)
	response.OkWithData(plugin, c)
	return
}

// @Summary plugin update
// @Tags Plugin
// @Description 更新
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Param plugin body rule.Plugin true "plugin"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/{id}/ [put]
func (m *PluginsApi) Update(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	plugin := my.Plugin{}
	err := c.ShouldBindJSON(&plugin)
	if err != nil {
		response.FailWithMessage("参数不合法", c)
		return
	}
	if my.ExistPluginByID(id) {
		my.EditPlugin(id, plugin)
		response.OkWithData(plugin, c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}
}

// @Summary plugin delete
// @Tags Plugin
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/{id}/ [delete]
func (m *PluginsApi) Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	if my.ExistPluginByID(id) {
		my.DeletePlugin(id)
		response.OkWithMessage("删除成功", c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}
}

// @Summary plugin run
// @Tags Plugin
// @Description 运行
// @Produce  json
// @Security token
// @Param run body RunSwigger false "run"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/run/ [post]
func (m *PluginsApi) Run(c *gin.Context) {
	run := RunSerializer{}
	err := c.ShouldBindJSON(&run)
	if err != nil {
		response.FailWithMessage("漏洞编号以poc-开头 测试url和漏洞类型不可为空", c)
		return
	}

	oreq, err := util.GenOriginalReq(run.Target)
	if err != nil {
		response.FailWithMessage("目标连通性不通过/原始请求生成失败", c)
		return
	}
	poc, err := rule.ParseJsonPoc(run.JsonPoc)
	if err != nil {
		log.Fatal("[plugins.go run] fail to load plugins", err)
		response.FailWithMessage("规则加载失败", c)
		return
	}

	task := my.Task{
		Remarks: "single poc",
		Target:  run.Target,
	}
	my.AddTask(&task)

	currentPlugin := rule.Plugin{
		Affects: run.Affects,
		JsonPoc: poc,
		VulId:   run.VulId,
	}

	item := &rule.ScanItem{OriginalReq: oreq, Plugin: &currentPlugin, Task: &task}

	result, err := rule.RunPoc(item, true)
	if err != nil {
		my.ErrorTask(task.Id)
		response.FailWithMessage("规则运行失败："+err.Error(), c)
		return
	}
	my.DownTask(task.Id)
	response.OkWithData(result, c)
	return
}

// @Summary download yaml
// @Tags Plugin
// @Description 下载yaml
// @Security token
// @Param run body DownloadSwigger true "json_poc"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/download/ [post]
func (m *PluginsApi) DownloadYaml(c *gin.Context) {
	download := DownloadSerializer{}
	err := c.ShouldBindJSON(&download)
	if err != nil {
		response.FailWithMessage("规则格式不正确", c)
		return
	}
	poc, err := rule.ParseJsonPoc(download.JsonPoc)
	if err != nil {
		log.Fatal("[plugins.go download] fail to load plugins", err)
		response.FailWithMessage("规则解析失败", c)
		return
	}
	content, err := yaml.Marshal(poc)
	if err != nil {
		log.Fatal("[plugins.go download] fail to marshal yaml", err)
		response.FailWithMessage("yaml生成失败", c)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", poc.Name))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Writer.Write(content)
}

// @Summary upload yaml
// @Tags Plugin
// @Description 上传yaml
// @Accept multipart/form-data
// @Param yaml formData file true "file"
// @accept json
// @Security token
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/poc/upload/ [post]
func (m *PluginsApi) UploadYaml(c *gin.Context) {
	file, _, err := c.Request.FormFile("yaml")
	if err != nil {
		response.FailWithMessage("文件上传失败", c)
		return
	}
	// 获取yaml内容
	content, err := ioutil.ReadAll(file)
	if err != nil {
		response.FailWithMessage("文件读取失败，请检查后重试", c)
		return
	}
	poc, err := rule.ParseYamlPoc(content)
	if err != nil {
		response.FailWithMessage("yaml解析失败，请检查后重试", c)
		return
	}
	// todo slice to map
	toMap := TempPoc{
		Params: poc.Params,
		Name:   poc.Name,
		Set:    SliceToMap(poc.Set),
		Rules:  poc.Rules,
		Groups: poc.Groups,
		Detail: rule.Detail{},
	}
	data := make(map[string]interface{})
	data["json_poc"] = toMap
	response.OkWithData(data, c)
}

type TempPoc struct {
	Params []string               `json:"params"`
	Name   string                 `json:"name"`
	Set    map[string]string      `json:"set"`
	Rules  []rule.Rule            `json:"rules"`
	Groups map[string][]rule.Rule `json:"groups"`
	Detail rule.Detail            `json:"detail"`
}

func SliceToMap(slice yaml.MapSlice) map[string]string {
	m := make(map[string]string)
	for _, v := range slice {
		m[v.Key.(string)] = v.Value.(string)
	}
	return m
}
