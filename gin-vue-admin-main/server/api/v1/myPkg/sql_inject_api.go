/*
 * @author Crabin
 */

package myPkg

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	myRes "github.com/flipped-aurora/gin-vue-admin/server/model/my/response"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"sync"
	"time"
)

type SqlInjectApi struct {
}

func (m *SqlInjectApi) CheckSqlInjection(c *gin.Context) {
	reqBody := make(map[string]interface{})
	c.BindJSON(&reqBody)
	target := reqBody["host"].(string)
	paramStr := reqBody["params"].(string)

	stateTime := time.Now()
	//拿到所有的接口
	urlsAll := crawlAllApi.GetAllLinks(target)
	if urlsAll == nil {
		response.FailWithMessage("这个url访问异常", c)
		return
	}

	log.Println("爬取API结束，耗时：", time.Since(stateTime))

	// 去掉重复
	urlsAll = removeDuplicates(urlsAll)

	//接收所有可能的url
	urlsCh := make(chan common.VulnerableURL, len(urlsAll))

	//参数处理
	params, err := parseParams(paramStr)
	if err != nil {
		log.Println(err)
	}
	log.Println("开始测试，测试目标一共：", len(urlsAll))
	//开始测试
	checkURLs(urlsAll, urlsCh, params)

	//urlsIn := []common.VulnerableURL{}
	//总计
	total := []common.VulnerableURL{}
	generalMeasurement := []common.VulnerableURL{}
	timeBlind := []common.VulnerableURL{}
	spadeTime := time.Since(stateTime)
	for vulnUrl := range urlsCh {
		total = append(total, vulnUrl)
		if vulnUrl.Type == 1 || vulnUrl.Type == 12 {
			generalMeasurement = append(generalMeasurement, vulnUrl)
		} else if vulnUrl.Type == 2 || vulnUrl.Type == 22 {
			timeBlind = append(timeBlind, vulnUrl)
		}
	}

	fmt.Printf("存在sql漏洞的url一共有：%v \n", len(total))
	fmt.Printf("get,post请求payload测出一共有：%v \n", len(generalMeasurement))
	fmt.Printf("时间盲注测出一共有：%v \n", len(timeBlind))
	fmt.Printf("耗时：%v s", spadeTime)

	response.OkWithDetailed(
		myRes.Sql{
			Urls:               total,
			GeneralMeasurement: generalMeasurement,
			TimeBlind:          timeBlind,
			SpadeTime:          int(spadeTime.Seconds()),
		},
		"测试完成",
		c,
	)
}

func checkURLs(vulnerableURLs []common.VulnerableURL, ch chan common.VulnerableURL, params []string) {
	var wg sync.WaitGroup
	wg.Add(len(vulnerableURLs))
	for _, vulnerableURL := range vulnerableURLs {
		urlInit := vulnerableURL.URL
		method := vulnerableURL.Method
		go func(urlInit string) {
			defer wg.Done()

			//首先使用报错检测
			for _, payload := range common.Payloads {
				if sqlInject.CheckUrlGet(urlInit, params, payload) {
					vulnerableURL.Type = 1
					vulnerableURL.URL = urlInit
					vulnerableURL.Params = params
					vulnerableURL.Method = method
					ch <- vulnerableURL
					return
				}
				if sqlInject.CheckUrlPost(urlInit, params, payload) {
					vulnerableURL.Type = 12
					vulnerableURL.URL = urlInit
					vulnerableURL.Params = params
					vulnerableURL.Method = method
					ch <- vulnerableURL
					return
				}
			}
			//使用时间盲注检测
			if sqlInject.CheckUrlBlindInjectionGet(urlInit, params) {
				vulnerableURL.Type = 2
				vulnerableURL.URL = urlInit
				vulnerableURL.Params = params
				vulnerableURL.Method = method
				ch <- vulnerableURL
			}
			if sqlInject.CheckUrlBlindInjectionPost(urlInit, params) {
				vulnerableURL.Type = 22
				vulnerableURL.URL = urlInit
				vulnerableURL.Params = params
				vulnerableURL.Method = method
				ch <- vulnerableURL
			}
		}(urlInit)
	}
	wg.Wait()
	close(ch)
}

// 去掉数组重复项
func removeDuplicates(elements []common.VulnerableURL) []common.VulnerableURL {
	// Define a map
	encountered := map[string]bool{}
	//定义返回结果
	res := []common.VulnerableURL{}
	for _, element := range elements {
		if encountered[element.URL] != true {
			encountered[element.URL] = true
			res = append(res, element)
		}
	}
	return res
}

func parseParams(params string) ([]string, error) {
	if len(params) == 0 {
		return nil, errors.New("参数为空")
	}
	res := []string{}

	//去掉空格
	params = strings.ReplaceAll(params, " ", "")

	res = strings.Split(params, ";")

	if len(res) == 0 {
		return nil, errors.New("发生未知错误")
	}

	return res, nil
}
