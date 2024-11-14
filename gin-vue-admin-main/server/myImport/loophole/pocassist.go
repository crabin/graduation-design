package loophole

import (
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/conf"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/logging"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/util"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/poc/rule"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
)

func InitAll() {
	// config 必须最先加载
	conf.Setup()
	logging.Setup()
	util.Setup()
	rule.Setup()
}

// 使用viper 对配置热加载
func HotConf() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("cmd.HotConf, fail to get current path: %v", err)
	}
	// 配置文件路径 当前文件夹 + config.yaml
	configFile := path.Join(dir, conf.ConfigFileName)
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	// watch 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		log.Println("config file changed:", e.Name)
		InitAll()
	})
}

//
//func RunApp() {
//	app := cli.NewApp()
//	app.Name = conf.ServiceName
//	app.Usage = conf.Website
//	app.Version = conf.Version
//
//	app.Flags = []cli.Flag{
//		&cli.StringFlag{
//			// 后端端口
//			Name:    "port",
//			Aliases: []string{"p"},
//			Value:   conf.DefaultPort,
//			Usage:   "web server `PORT`",
//		},
//	}
//	app.Action = RunServer
//
//	err := app.Run(os.Args)
//	if err != nil {
//		log.Fatalf("cli.RunApp err: %v", err)
//		return
//	}
//}
//
