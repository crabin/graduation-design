/*
 * @author Crabin
 */

package cmsDiscriminate

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/service/myPkg/cmsDiscriminate/engine"
	"github.com/flipped-aurora/gin-vue-admin/server/service/myPkg/cmsDiscriminate/until"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/wsConn"
	"log"
	"strings"
	"sync"
	"time"
)

type CMSService struct {
}

func (m *CMSService) CMSDiscriminate(wConn *wsConn.WsConnection, targets string) {

	domains := strings.Split(targets, ";")

	if domains == nil || len(domains) == 0 {
		return
	}

	// 加载指纹
	sortPairs, webdata := until.ParseCmsDataFromFile(common.CmsJsonPath)
	var waitGroup sync.WaitGroup

	// 开始并发相关
	t1 := time.Now()
	ResultChian := make(chan string)
	fmt.Println("Load url:", domains)
	for _, domain := range domains {
		go func(d string) {
			newWorker := engine.NewWorker(7, d, &waitGroup, ResultChian)
			if !newWorker.Checkout() {
				return
			}
			newWorker.Start()
			for _, v := range sortPairs {
				tmp_job := engine.JobStruct{d, v.Path, webdata[v.Path]}
				//fmt.Println(tmp_job)
				newWorker.Add(tmp_job)
			}
		}(domain)
	}
	time.Sleep(time.Second * 2)
	go func() {
		for {
			// 把消息传输给客户端
			r := <-ResultChian
			log.Println("result=======" + r)
			wConn.WriteMessage(common.SuccessType, []byte(r))
		}
	}()
	log.Println("初始化完成")

	waitGroup.Wait()
	elapsed := time.Since(t1)
	fmt.Println("耗时：", elapsed.String())
	//wConn.WriteMessage(common.SuccessType, []byte("耗时："+elapsed.String()))
	time.Sleep(2 * time.Second)
	// 完成
	wConn.WriteMessage(common.SuccessType, []byte(common.Accomplish))
}
