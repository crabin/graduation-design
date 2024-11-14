package kunpeng

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/opensec-cn/kunpeng/config"
	"github.com/opensec-cn/kunpeng/plugin"
	_ "github.com/opensec-cn/kunpeng/plugin/go"
	_ "github.com/opensec-cn/kunpeng/plugin/json"
	"github.com/opensec-cn/kunpeng/util"
	"github.com/opensec-cn/kunpeng/web"
)

var VERSION string

type greeting string

func (g greeting) Check(taskJSON string) []map[string]interface{} {
	var task plugin.Task
	json.Unmarshal([]byte(taskJSON), &task)
	return plugin.Scan(task)
}

func (g greeting) GetPlugins() []map[string]interface{} {
	return plugin.GetPlugins()
}

func (g greeting) SetConfig(configJSON string) {
	config.Set(configJSON)
}

func (g greeting) ShowLog() {
	config.SetDebug(true)
}

func (g greeting) GetVersion() string {
	return VERSION
}

func (g greeting) StartBuffer() {
	util.Logger.StartBuffer()
}

func (g greeting) GetLog(sep string) string {
	return util.Logger.BufferContent(sep)
}

//export StartWebServer
func StartWebServer(bindAddr *C.char) {
	go web.StartServer(C.GoString(bindAddr))
}

//export Check
func Check(task *C.char) *C.char {
	util.Logger.Info(C.GoString(task))
	var m plugin.Task
	err := json.Unmarshal([]byte(C.GoString(task)), &m)
	if err != nil {
		util.Logger.Error(err.Error())
		return C.CString("[]")
	}
	util.Logger.Info(m)
	result := plugin.Scan(m)
	if len(result) == 0 {
		return C.CString("[]")
	}
	b, err := json.Marshal(result)
	if err != nil {
		util.Logger.Error(err.Error())
		return C.CString("[]")
	}
	return C.CString(string(b))
}

//export GetPlugins
func GetPlugins() *C.char {
	var result string
	plugins := plugin.GetPlugins()
	b, err := json.Marshal(plugins)
	if err != nil {
		util.Logger.Error(err.Error())
		return C.CString("[]")
	}
	result = string(b)
	return C.CString(result)
}

//export SetConfig
func SetConfig(configJSON *C.char) {
	config.Set(C.GoString(configJSON))
}

//export ShowLog
func ShowLog() {
	config.SetDebug(true)
}

//export GetVersion
func GetVersion() *C.char {
	return C.CString(VERSION)
}

//export StartBuffer
func StartBuffer() {
	util.Logger.StartBuffer()
}

//export GetLog
func GetLog(sep *C.char) *C.char {
	return C.CString(util.Logger.BufferContent(C.GoString(sep)))
}

var Greeter greeting

type Meta struct {
	System   string   `json:"system"`
	PathList []string `json:"pathlist"`
	FileList []string `json:"filelist"`
	PassList []string `json:"passlist"`
}

type Task struct {
	Type   string `json:"type"`
	Netloc string `json:"netloc"`
	Target string `json:"target"`
	Meta   Meta   `json:"meta"`
}

func main() {
	// 开启日志打印
	ShowLog()
	pluginInfo := GetPlugins()
	// 获取插件信息
	fmt.Println(pluginInfo)

	// 扫描目标
	task := Task{
		Type:   "service",
		Netloc: "192.168.46.128:3306",
		Target: "mysql",
		Meta: Meta{
			System:   "",
			PathList: []string{},
			FileList: []string{},
			PassList: []string{},
		},
	}
	jsonBytes, _ := json.Marshal(task)
	result := Greeter.Check(string(jsonBytes))
	fmt.Println(result)

	task2 := Task{
		Type:   "web",
		Netloc: "http://www.google.cn",
		Target: "web",
		Meta: Meta{
			System:   "",
			PathList: []string{},
			FileList: []string{},
			PassList: []string{},
		},
	}

	// time.Sleep(time.Second * 21)
	jsonBytes, _ = json.Marshal(task2)
	result = Greeter.Check(string(jsonBytes))
	fmt.Println(result)

}
