// Package config 配置信息定义
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Timeout         int      `json:"timeout"`
	Aider           string   `json:"aider"`
	HTTPProxy       string   `json:"http_proxy"`
	PassList        []string `json:"pass_list"`
	ExtraPluginPath string   `json:"extra_plugin_path"`
}

// Debug 为True时打印过程日志
var Debug bool

// Config 全局配置信息
var Config config

func init() {
	//加载字典
	list, err := LoadDictionary("")
	if err != nil {
		fmt.Println("加载字典文件失败：", err)
		return
	}
	Config.PassList = list
	passList := []string{
		"{user}", "{user}123",
	}
	for _, s := range passList {
		Config.PassList = append(Config.PassList, s)
	}

	Config.Timeout = 15
	Debug = false
}

// Set 设置配置信息
func Set(configJSON string) {
	json.Unmarshal([]byte(configJSON), &Config)
	if Config.Timeout == 0 {
		Config.Timeout = 15
	}
	if !strings.HasSuffix(Config.ExtraPluginPath, "/") {
		Config.ExtraPluginPath = Config.ExtraPluginPath + "/"
	}
}

// SetDebug 是否开启debug，即打印日志
func SetDebug(debug bool) {
	Debug = debug
}

// 加载字典
func LoadDictionary(filePath string) ([]string, error) {
	if filePath == "" {
		//弱口令字典
		filePath = "files/weakPass.txt"
	}

	filePath, _ = filepath.Abs(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 逐行读取并保存到数组中
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
