/*
 * @author Crabin
 */

package myPkg

import (
	"bufio"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DirbScanApi struct {
}

func (d *DirbScanApi) DirbScanHandler(c *gin.Context) {
	// 升级为 websokect
	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 1000,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	stateTime := time.Now()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	host := c.Query("host")

	//查看这个是否是url格式
	if !utils.MatchUrl(host) {
		conn.WriteMessage(websocket.TextMessage, []byte("\n\n输入格式不正确！"))
		return
	}

	out, _ := utils.Mkdir("out")
	out = filepath.Join(out, timer.GetTimeFormat()+"_dirb.txt")

	cmdStr := common.GobusterPath + " dir -u " + host + " -w " + common.CommTextPath + " -o " + out + " -q -z"
	cmdParts := strings.Fields(cmdStr)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	// 结果计数
	resNum := 0
	// 实时返回
	for scanner.Scan() {
		message := scanner.Text()
		log.Println(message)
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			break
		}
		resNum++
	}

	// 等待cmd运行结束
	if err := cmd.Wait(); err != nil {
		log.Println("执行完成1")
		log.Println(err)
	}
	//执行结束，返回结果

	endTime := time.Now()
	elapsedTime := endTime.Sub(stateTime).Seconds()
	msg := "\n\n总记录：" + strconv.Itoa(resNum) + "条  总耗时：" + strconv.Itoa(int(math.Floor(elapsedTime))) + "秒"
	conn.WriteMessage(websocket.TextMessage, []byte(msg))

}
