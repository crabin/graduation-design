/*
 * @author Crabin
 */

package common

import (
	"os"
	"path/filepath"
)

var (
	// 根目录： G:\WorkSpace\Go_WorkSpace\graduation-design\gin-vue-admin-main\server
	RootDir, err = os.Getwd()

	// G:\WorkSpace\Go_WorkSpace\graduation-design\gin-vue-admin-main\server\myImport
	MyImportPath = filepath.Join(RootDir, "myImport")

	DirbPath = filepath.Join(RootDir, "myImport", "dirb")

	// 字典文件目录： G:\WorkSpace\Go_WorkSpace\graduation-design\gin-vue-admin-main\server\myImport\...
	CommTextPath = filepath.Join(RootDir, "myImport", "dirb", "common.txt")

	// 工具Gobuster目录： G:\WorkSpace\Go_WorkSpace\graduation-design\gin-vue-admin-main\server\myImport\...
	GobusterPath = filepath.Join(RootDir, "myImport", "gobuster-master", "main.exe")

	// 输出文件夹路径
	OutPath = filepath.Join(RootDir, "myImport", "out")

	// cms指纹识别cms.json路径
	CmsJsonPath = filepath.Join(RootDir, "myImport", "myFile", "cms.json")

	// cms指纹识别waf.txt路径，识别waf
	CmsWafPath = filepath.Join(RootDir, "myImport", "myFile", "waf.txt")
)
