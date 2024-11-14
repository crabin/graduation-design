/*
 * @author Crabin
 */

package myPgk

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole"
	"github.com/gin-gonic/gin"
)

type InitApi struct {
}

func (m *InitApi) InitVulData(c *gin.Context) {
	loophole.InitAll()
	loophole.HotConf()
	response.OkWithMessage("初始化数据完成！", c)
}
