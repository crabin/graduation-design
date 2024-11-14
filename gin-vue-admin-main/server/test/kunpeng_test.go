/*
 * @author Crabin
 */

package test

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/kunpeng"
	"testing"
)

func TestKunpeng(t *testing.T) {
	// 开启日志打印
	kunpeng.ShowLog()
	pluginInfo := kunpeng.GetPlugins()
	// 获取插件信息
	fmt.Println(pluginInfo)

	// 扫描目标
	task := kunpeng.Task{
		Type:   "service",
		Netloc: "192.168.46.128:3306",
		Target: "mysql",
		Meta: kunpeng.Meta{
			System:   "",
			PathList: []string{},
			FileList: []string{},
			PassList: []string{},
		},
	}
	jsonBytes, _ := json.Marshal(task)
	result := kunpeng.Greeter.Check(string(jsonBytes))
	fmt.Println(result)

	task2 := kunpeng.Task{
		Type:   "web",
		Netloc: "http://www.google.cn",
		Target: "web",
		Meta: kunpeng.Meta{
			System:   "",
			PathList: []string{},
			FileList: []string{},
			PassList: []string{},
		},
	}

	// time.Sleep(time.Second * 21)
	jsonBytes, _ = json.Marshal(task2)
	result = kunpeng.Greeter.Check(string(jsonBytes))
	fmt.Println(result)
}
