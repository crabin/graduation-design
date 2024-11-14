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

type TaskApi struct {
}

// @Summary task list
// @Tags Task
// @Description 列表
// @Produce  json
// @Security token
// @Param page query int true "Page"
// @Param pagesize query int true "Pagesize"
// @Param object query db.TaskSearchField false "field"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/task/ [get]
func (m *TaskApi) Get(c *gin.Context) {
	data := make(map[string]interface{})
	field := my.TaskSearchField{Search: ""}

	// 分页
	page, _ := com.StrTo(c.Query("page")).Int()
	pageSize, _ := com.StrTo(c.Query("pagesize")).Int()

	// 查询条件
	if arg := c.Query("search"); arg != "" {
		field.Search = arg
	}

	tasks := my.GetTask(page, pageSize, &field)
	data["data"] = tasks

	total := my.GetTaskTotal(&field)
	data["total"] = total

	response.OkWithData(data, c)
	return
}

// @Summary task delete
// @Tags Task
// @Description 删除
// @Produce  json
// @Security token
// @Param id path int true "ID"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/task/{id}/ [delete]
func (m *TaskApi) Delete(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	if my.ExistTaskByID(id) {
		my.DeleteTask(id)
		response.OkWithMessage("删除成功", c)
		return
	} else {
		response.FailWithMessage("record not found", c)
		return
	}

}
