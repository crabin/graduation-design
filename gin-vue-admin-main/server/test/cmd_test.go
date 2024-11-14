/*
 * @author Crabin
 */

package test

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCmd(t *testing.T) {
	//cmd := exec.Command("ping", "www.baidu.com")

	out, _ := utils.Mkdir("out")
	out = filepath.Join(out, "com_dirb.txt")
	cmdStr := common.GobusterPath + " dir -u https://www.cnblogs.com/crabin/ -w " + common.CommTextPath
	cmdParts := strings.Fields(cmdStr)
	//cmd := exec.Command("go", "1")
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
