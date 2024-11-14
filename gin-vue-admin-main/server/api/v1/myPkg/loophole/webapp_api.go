/*
 * @author Crabin
 */

package myPgk

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/my"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type WebAppApi struct {
}

// @Summary product detail
// @Tags Product
// @Description 详情
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/{id}/ [get]
func (m *WebAppApi) Detail(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	var data interface{}
	if my.ExistWebappById(id) {
		data = my.GetWebapp(id)
		response.OkWithData(data, c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}
}

// @Summary product list
// @Tags Product
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.WebappSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/ [get]
func (m *WebAppApi) Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := my.WebappSearchField{Search: ""}
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}

	apps := my.GetWebapps(page, pageSize, &field)
	data["data"] = apps
	total := my.GetWebappsTotal(&field)
	data["total"] = total
	response.OkWithData(data, c)
	return
}

// @Summary product add
// @Tags Product
// @Description 新增
// @Produce  json
// @Security token
// @Param plugin body db.Webapp true "webapp"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/ [post]
func (m *WebAppApi) Add(c *gin.Context) {
	app := my.Webapp{}
	err := c.BindJSON(&app)
	if err != nil {
		response.FailWithMessage("参数校验不通过", c)
		return
	}
	if my.ExistWebappByName(app.Name) {
		response.FailWithMessage("漏洞名称已存在", c)
		return
	} else {
		my.AddWebapp(app)
		response.OkWithData(app, c)
		return
	}
}

// @Summary product update
// @Tags Product
// @Description 更新
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Param webapp body db.Webapp true "webapp"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/{id}/ [put]
func (m *WebAppApi) Update(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	app := my.Webapp{}
	err := c.ShouldBindJSON(&app)
	if err != nil {
		response.FailWithMessage("组件名称不可为空", c)
		return
	}
	if my.ExistWebappById(id) {
		my.EditWebapp(id, app)
		response.OkWithData(app, c)
	} else {
		response.FailWithMessage("record not found", c)
		return
	}
}

// @Summary product delete
// @Tags Product
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/product/{id}/ [delete]
func (m *WebAppApi) Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	if my.ExistWebappById(id) {
		my.DeleteWebapp(id)
		response.OkWithMessage("删除成功", c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}

}
