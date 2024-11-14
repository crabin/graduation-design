/*
 * @author Crabin
 */

package test

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"os/exec"
	"strings"
	"testing"
	"unicode/utf8"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsCmdHandle() {
	r := gin.Default()
	r.GET("/ws", handleWebSocket)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func TestWsCmd(t *testing.T) {
	//cmd := exec.Command("ping", "www.baidu.com")
	cmdStr := "go run ..\\..\\..\\myImport/gobuster-3.4.0/main.go dir -u https://www.cnblogs.com/crabin/ -w common.txt"
	cmdParts := strings.Fields(cmdStr)
	//cmd := exec.Command("go", "run ../myImport/gobuster-3.4.0/main.go dir -u https://www.cnblogs.com/crabin/ -w G:\\WorkSpace\\Go_WorkSpace\\graduation-design\\gin-vue-admin-main\\server\\myImport\\dirb\\common.txt")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(stdout)

	//设置编码

	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) > 0 {
			r, _ := utf8.DecodeLastRune(b)
			fmt.Printf("%c", r)
		}
		fmt.Println(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		//_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		cmdStr := "go run ./myImport/gobuster-3.4.0/main.go dirb -u https://www.cnblogs.com/crabin/ -w G:\\WorkSpace\\Go_WorkSpace\\graduation-design\\gin-vue-admin-main\\server\\myImport\\dirb\\common.txt"
		cmdParts := strings.Fields(cmdStr)

		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			break
		}

		if err := cmd.Start(); err != nil {
			log.Println(err)
			break
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			message := scanner.Text()
			err = conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				break
			}
		}

		if err := cmd.Wait(); err != nil {
			log.Println(err)
			break
		}

	}
}
