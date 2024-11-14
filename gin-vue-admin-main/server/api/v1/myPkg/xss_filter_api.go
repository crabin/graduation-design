/*
 * @author Crabin
 */

package myPkg

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	myRes "github.com/flipped-aurora/gin-vue-admin/server/model/my/response"
	"github.com/gin-gonic/gin"
)

type XSSFilterApi struct {
}

// XSS过滤处理函数
func (m *XSSFilterApi) XSSFilterHandle(c *gin.Context) {
	reqBody := make(map[string]interface{})
	c.BindJSON(&reqBody)

	filterText := reqBody["text"].(string)

	safeHtml, xssSentences := xssFilterServer.XSSFilter(filterText)

	response.OkWithData(myRes.FilterRes{
		SafeHtml:     safeHtml,
		XssSentences: xssSentences,
	}, c)
}
