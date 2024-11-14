/*
 * @author Crabin
 */

package myPkg

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/wsConn"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ipinfo/go/v2/ipinfo"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	myRes "github.com/flipped-aurora/gin-vue-admin/server/model/my/response"
)

type IpInfoApi struct {
}

func (m *IpInfoApi) GetIPInfo(c *gin.Context) {
	reqBody := make(map[string]interface{})
	c.BindJSON(&reqBody)

	ipAddress := reqBody["ip"].(string)

	//去掉首尾空格
	ipAddress = strings.TrimSpace(ipAddress)

	if ipAddress == "" {
		response.FailWithMessage("请输入您要查询的信息", c)
		return
	}

	res := myRes.IpInfos{}

	if !utils.IsValidIP(ipAddress) {
		// 是网址
		urlInit := ipAddress
		url, err := url.Parse(ipAddress)
		if url.Host == "" || err != nil {
			response.FailWithMessage("输入正确的网址或者IP", c)
			return
		}
		ips := utils.GetIpLookupWithUrl(url.Host)
		ipAddress = ips[0].String()
		serviceType, _ := utils.CheckProtocolType(urlInit)
		res.ServerType = serviceType
		res.XPoweredBy, _ = utils.CheckXPoweredBy(urlInit)
		res.WebInfo = utils.GetWebInfo(urlInit)
	}

	// params: httpClient, cache, token. `http.DefaultClient` and no cache will be used in case of `nil`.
	client := ipinfo.NewClient(nil, nil, common.IpInfoToten)
	info, err := client.GetIPInfo(net.ParseIP(ipAddress))

	if err != nil {
		log.Fatal(err)
		return
	}
	res.Info = info
	response.OkWithData(res, c)
}

func (m *IpInfoApi) GetIpAddress(c *gin.Context) {
	ipAddress := utils.GetClientIP(c)
	if ipAddress == "" {
		response.FailWithMessage("获取失败", c)
		return
	}
	res := map[string]string{}
	res["Ip"] = ipAddress
	response.OkWithData(res, c)
}

// 追踪路由
func (m *IpInfoApi) TracertHostWs(c *gin.Context) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 1000,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}

	//log.Println("建立连接")
	// 升级http为websokect
	wsSocket, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println(err)
	}
	wConn = wsConn.New(wsSocket)

	data, err := wConn.ReadMessage()
	if err != nil {
		wConn.Close()
		return
	}

	if data.Data != nil {
		//接收前端传过来信息，输入信息
		urlStr := string(data.Data)
		url, err := url.Parse(urlStr)
		if err != nil || url.Host == "" {
			d1, _ := utils.StructToBytes(myRes.IpInfos{Status: 0, Msg: "输入的格式不正确"})
			wConn.WriteMessage(1, d1)
			defer wConn.Close()
			log.Println("输入的格式不正确")
			return
		}
		//开启追踪
		go ipInfoServer.Tracert(url.Host, wConn)
	}

	if err := wConn.WriteMessage(data.MessageType, data.Data); err != nil {
		wConn.Close()
		return
	}

}
