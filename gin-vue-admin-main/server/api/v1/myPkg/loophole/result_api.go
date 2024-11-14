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

type ResultApi struct {
}

// @Summary result list
// @Tags Result
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.ResultSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/result/ [get]
func (m *ResultApi) Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := my.ResultSearchField{Search: "", TaskField: -1, VulField: -1}
	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}
	if arg := c.Query("taskField"); arg != "" {
		enable := com.StrTo(arg).MustInt()
		field.TaskField = enable
	}
	if arg := c.Query("vulField"); arg != "" {
		vul := com.StrTo(arg).MustInt()
		field.VulField = vul
	}
	results := my.GetResult(page, pageSize, &field)
	data["data"] = results

	total := my.GetResultTotal(&field)
	data["total"] = total

	response.OkWithData(data, c)
	return
}

// @Summary result delete
// @Tags Result
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/result/{id}/ [delete]
func (m *ResultApi) Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	if my.ExistResultByID(id) {
		my.DeleteResult(id)
		response.OkWithMessage("删除成功", c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}
}
