package rule

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/my"
	"github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/conf"
	log2 "github.com/flipped-aurora/gin-vue-admin/server/myImport/loophole/pkg/logging"
	"github.com/panjf2000/ants/v2"
	"gopkg.in/yaml.v2"
	_ "log"
	"net/http"
	"sync"
)

const (
	LoadMulti = "multi"
)

func ParseJsonPoc(jsonByte []byte) (*Poc, error) {
	poc := &Poc{}
	err := yaml.Unmarshal(jsonByte, poc)
	if poc.Name == "" {
		errMsg := "poc解析失败，poc名称不可为空"
		log2.Error("rule/plugin.go:ParseJsonPoc Err", errMsg)
		return nil, errors.New(errMsg)
	}
	return poc, err
}

func ParseYamlPoc(yamlByte []byte) (*Poc, error) {
	poc := &Poc{}
	err := yaml.Unmarshal(yamlByte, poc)
	if poc.Name == "" {
		errMsg := "poc parse error"
		log2.Error("rule/plugin.go:ParseJsonPoc Err", errMsg)
		return nil, errors.New(errMsg)
	}
	return poc, err
}

// 限制并发
type TaskItem struct {
	OriginalReq *http.Request // 原始请求
	Plugins     []Plugin      // 检测插件
	Task        *my.Task      // 所属任务
}

var TaskChannel chan *TaskItem

func InitTaskChannel() {
	// channel 限制 target 并发
	concurrent := 10
	if conf.GlobalConfig.PluginsConfig.Concurrent != 0 {
		concurrent = conf.GlobalConfig.PluginsConfig.Concurrent
	}
	TaskChannel = make(chan *TaskItem, concurrent)
}

func (item *TaskItem) Verify() error {
	errMsg := ""
	if item.Task == nil {
		errMsg = "task create fail"
		log2.Error("[rule/parallel.go:Verify error]", errMsg)
		return errors.New(errMsg)
	}
	if item.OriginalReq == nil {
		errMsg = "not original request"
		log2.Error("[rule/parallel.go:Verify error]", errMsg)
		return errors.New(errMsg)
	}
	if len(item.Plugins) == 0 {
		errMsg = "not plugin"
		log2.Error("[rule/parallel.go:Verify error]", errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// 从数据库 中加载 POC
func LoadDbPlugin(lodeType string, array []string) ([]Plugin, error) {
	// 数据库数据
	var dbPluginList []my.Plugin
	// plugin对象
	var plugins []Plugin
	switch lodeType {
	case LoadMulti:
		// 多个
		tx := global.GVA_DB.Where("vul_id IN ? AND enable = ?", array, 1).Find(&dbPluginList)
		if tx.Error != nil {
			log2.Error("[rule/parallel.go:LoadDbPlugin load multi err]", tx.Error)
			return nil, tx.Error
		}
	default:
		// 默认执行全部启用规则
		tx := global.GVA_DB.Where("enable = ?", 1).Find(&dbPluginList)
		if tx.Error != nil {
			log2.Error("[rule/parallel.go:LoadDbPlugin load all err]", tx.Error)
			return nil, tx.Error
		}
	}
	log2.Info("[rule/parallel.go:LoadDbPlugin load plugin number]", len(dbPluginList))

	for _, v := range dbPluginList {
		poc, err := ParseJsonPoc(v.JsonPoc)
		if err != nil {
			continue
		}
		plugin := Plugin{
			VulId:   v.VulId,
			Affects: v.Affects,
			JsonPoc: poc,
			Enable:  v.Enable,
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

func TaskProducer(item *TaskItem) {
	TaskChannel <- item
}

func TaskConsumer() {
	for item := range TaskChannel {
		// 校验格式
		err := item.Verify()
		if err != nil {
			log2.Error("[rule/poc.go:WriteTaskError scan error] ", err)
			my.ErrorTask(item.Task.Id)
			continue
		}
		RunPlugins(item)
	}
}

// 并发测试
func RunPlugins(item *TaskItem) {
	// 限制插件并发数
	var wg sync.WaitGroup
	parallel := conf.GlobalConfig.PluginsConfig.Parallel
	p, _ := ants.NewPoolWithFunc(parallel, func(item interface{}) {
		RunPoc(item, false)
		wg.Done()
	})
	defer p.Release()
	oreq := item.OriginalReq
	plugins := item.Plugins
	task := item.Task

	log2.Info("[rule/parallel.go:TaskConsumer start scan]", oreq.URL.String())
	for i := range plugins {
		item := &ScanItem{oreq, &plugins[i], task}
		wg.Add(1)
		p.Invoke(item)
	}
	wg.Wait()
	my.DownTask(task.Id)
	//my.AddResult(my.Result{
	//	Vul:        ,
	//	Detail:     task.,
	//	PluginId:   "",
	//	PluginName: "",
	//	TaskId:     0,
	//})
}
