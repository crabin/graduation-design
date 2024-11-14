/*
 * @author Crabin
 */

package myPkg

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/wsConn"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

type CmsDiscriminateApi struct {
}

func (m *CmsDiscriminateApi) CmsDiscriminateHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 1000,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	// 升级http为websokect
	wsSocket, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println(err)
	}
	wConn = wsConn.New(wsSocket)

	targets := ""

	data, err := wConn.ReadMessage()
	if err != nil {
		wConn.Close()
		return
	}

	if data != nil {
		//接收前端传过来信息，输入信息
		targets = string(data.Data)
		url, err := url.Parse(targets)
		if err != nil || url.Host == "" {
			wConn.WriteMessage(common.SuccessType, []byte("输入的格式不正确！"))
			defer wConn.Close()
			log.Println("输入的格式不正确")
			return
		}
		// 开始识别cms指纹
		go cmsDiscriminateserver.CMSDiscriminate(wConn, targets)
	}

}
