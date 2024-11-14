package myPkg

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	myRes "github.com/flipped-aurora/gin-vue-admin/server/model/my/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/wsConn"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

type PortScanApi struct {
}

var (
	wConn *wsConn.WsConnection
)

func (m *PortScanApi) StateScan(c *gin.Context) {
	reqBody := make(map[string]interface{})
	c.BindJSON(&reqBody)

	scanIp := reqBody["ip"].(string)
	scanPort := reqBody["port"].(string)
	log.Println(scanIp)
	log.Println(scanPort)
	if scanIp == "" || scanPort == "" {
		response.Result(
			201,
			nil,
			"ip或端口号不能为空", c)
		return
	}

	var params = myRes.Params{
		Ip:      scanIp,
		Port:    scanPort,
		Process: 10,
		Timeout: 100,
		Debug:   1,
	}
	debug := true
	//初始化
	scanIP := portIp.NewScanIp(params.Timeout, params.Process, debug)
	ips, err := scanIP.GetAllIp(scanIp)
	if err != nil {
		response.Result(
			201,
			nil,
			fmt.Sprintf("ip解析错误！ %v", err.Error()), c)
		return
	}

	//扫所有的ip
	filePath, _ := utils.Mkdir("log")
	fileName := filePath + params.Ip + "_port.txt"

	var openPorts []int
	stateTime := time.Now()
	for i := 0; i < len(ips); i++ {
		ports := scanIP.GetIpOpenPort(ips[i], params.Port, wConn)
		if len(ports) > 0 {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				if err := f.Close(); err != nil {
					fmt.Println(err)
				}
				continue
			}
			if _, err := f.WriteString(fmt.Sprintf("%v【%v】开放:%v \n", time.Now().Format("2006-01-02 15:04:05"), ips[i], ports)); err != nil {
				if err := f.Close(); err != nil {
					fmt.Println(err)
				}
				continue
			}
			openPorts = append(openPorts, ports...)
		}
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(stateTime).Seconds()
	log.Println(openPorts)
	response.Result(
		200,
		myRes.Ports{openPorts, int(math.Floor(elapsedTime))},
		"扫描完成",
		c,
	)
}

// ws服务
func (m *PortScanApi) PortScanWs(c *gin.Context) {

	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 1000,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}

	log.Println("建立连接")
	// 升级http为websokect

	wsSocket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}
	wConn = wsConn.New(wsSocket)
	//defer wsSocket.Close()
	//wConn.WriteMessage(websocket.TextMessage, []byte("建立连接"))
	for {
		data, err := wConn.ReadMessage()
		if err != nil {
			wConn.Close()
			return
		}
		if err := wConn.WriteMessage(data.MessageType, data.Data); err != nil {
			wConn.Close()
			return
		}
	}
}

/*
demo
func (m *MyApi) Ws(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	// 升级HTTP请求为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 处理WebSocket消息...
	conn.WriteMessage(websocket.TextMessage, []byte("建立连接"))
}*/
