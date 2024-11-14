/*
 * @author Crabin
 */

package myPkg

import (
	"bytes"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SqlInject struct {
}

var ()

func getClient() *http.Client {
	client := &http.Client{
		Timeout: 20 * time.Second, // 设置读取超时时间
	}
	return client
}

func setHeader(req *http.Request) {
	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//reqPost.Header.Set("Content-Type", "application/json")
	} else {

	}

	// 设置随机数种子
	rand.Seed(time.Now().Unix())
	req.Header.Set("User-Agent", common.UserAgents[rand.Intn(len(common.UserAgents))])

}

// 设置请求参数
func setParams() {

}

func (s *SqlInject) CheckUrlBlindInjectionGet(urlStr string, params []string) bool {
	for _, injection := range common.Injections {
		client := getClient()
		// 发送请求 get
		//urlStrP := urlStr + "/?id=1'+and+(select+*+from+(select+if(1=1,sleep(5),1))x);--+"

		urlStrP, paramsUrl := utils.HandleUrl(urlStr)
		if paramsUrl != nil && len(paramsUrl) > 0 {
			for _, param := range paramsUrl {
				urlStrP = utils.MontageURL(urlStrP, param+"="+injection)
			}
		} else {
			for _, param := range common.ParamsCommon {
				urlStrP = utils.MontageURL(urlStrP, param+"="+injection)
			}
			for _, param := range params {
				urlStrP = utils.MontageURL(urlStrP, param+"="+injection)
			}
		}

		//urlStrP := utils.MontageURL(urlStr, "id="+injections[i])

		getReq, err := http.NewRequest(http.MethodGet, urlStrP, nil)
		if err != nil {
			return false
		}
		// 构造包含延迟的注入语句
		delay := time.Duration(common.InjectionTime) * time.Second

		setHeader(getReq)

		// 发送GET请求
		start := time.Now()

		// 获取响应
		getResp, err := client.Do(getReq)
		if err != nil {
			return false
		}
		defer getResp.Body.Close()
		getDody, err := ioutil.ReadAll(getResp.Body)
		if err != nil {
			return false
		}

		if getDody == nil {
			log.Println(urlStr + "返回为空")
		}

		elapsed := time.Since(start)
		// 根据响应状态码来判断是否存在 SQL 注入漏洞
		if getResp.StatusCode == http.StatusOK {
			// 检查响应时间是否大于注入语句中的延迟时间
			if elapsed > delay {
				log.Println(urlStr + "\t存在sql注入漏洞--get--时间盲注")
				return true
			}
		}

	}
	return false
}

func (s *SqlInject) CheckUrlBlindInjectionPost(urlStr string, params []string) bool {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr = urlStr + "/"
	}
	for _, injection := range common.Injections {

		client := getClient()
		// post请求判断
		postData := url.Values{}
		if params == nil || len(params) <= 0 {
			for _, param := range common.ParamsCommon {
				postData.Set(param, injection)
			}
		} else {
			for _, param := range params {
				postData.Set(param, injection)
			}
		}
		// 发送请求 get
		//urlStrP := urlStr + "/?id=1'+and+(select+*+from+(select+if(1=1,sleep(5),1))x);--+"

		// 将表单数据编码为application/x-www-form-urlencoded格式
		encodedFormData := postData.Encode()
		// 创建bytes.Buffer对象，并将编码后的表单数据写入其中
		requestBody := bytes.NewBufferString(encodedFormData)

		reqPost, err := http.NewRequest(http.MethodPost, urlStr, requestBody)

		setHeader(reqPost)
		//reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//// 设置随机数种子
		//rand.Seed(time.Now().Unix())
		//reqPost.Header.Set("User-Agent", common.UserAgents[rand.Intn(len(common.UserAgents))])
		//reqPost.Header.Set("Content-Type", "application/json")

		if err != nil {
			return false
		}
		// 构造包含延迟的注入语句
		delay := time.Duration(common.InjectionTime) * time.Second

		// 发送GET请求
		start := time.Now()

		// 获取响应
		postResp, err := client.Do(reqPost)
		if err != nil {
			return false
		}
		defer postResp.Body.Close()
		getDody, err := ioutil.ReadAll(postResp.Body)
		if err != nil {
			return false
		}

		if getDody == nil {
			log.Println(urlStr + "返回为空")
		}

		//
		elapseds := time.Since(start)
		// 响应状态码是否正常
		if postResp.StatusCode == http.StatusOK {
			// 检查响应时间是否大于注入语句中设置的延迟时间
			if elapseds > delay {
				log.Println(urlStr + "\t存在sql注入漏洞--post--时间盲注")
				return true
			}
		}

	}
	return false
}

// 测试get请求
func (s *SqlInject) CheckUrlGet(urlStr string, params []string, payload string) bool {
	client := getClient()
	// 发送请求 get
	getReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return false
	}

	setHeader(getReq)

	urlStr, paramsUrl := utils.HandleUrl(urlStr)
	getQ := getReq.URL.Query()
	if paramsUrl != nil && len(paramsUrl) > 0 {
		for key := range getQ {
			getQ.Del(key)
		}
		for _, param := range paramsUrl {
			// 请求参数
			getQ.Set(param, payload)
		}
	} else {
		//使用字典
		for _, param := range common.ParamsCommon {
			getQ.Set(param, payload)
		}
		for _, param := range params {
			getQ.Set(param, payload)
		}

	}

	getReq.URL.RawQuery = getQ.Encode()

	// 获取响应
	getResp, err := client.Do(getReq)
	if err != nil {
		return false
	}
	defer getResp.Body.Close()

	// 判断是否存在SQL注入漏洞
	getBody, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		return false
	}
	if strings.Contains(string(getBody), "SQL syntax") || strings.Contains(string(getBody), "mysql_fetch_array()") {
		log.Println(urlStr + "\t存在sql注入漏洞--get")
		return true
	}
	return false
}

// 测试post请求
func (s *SqlInject) CheckUrlPost(urlStr string, params []string, payload string) bool {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr = urlStr + "/"
	}
	client := getClient()
	// post请求判断
	postData := url.Values{}
	if params == nil || len(params) <= 0 {
		for _, param := range common.ParamsCommon {
			postData.Set(param, payload)
		}
	} else {
		for _, param := range params {
			postData.Set(param, payload)
		}
	}

	// 将表单数据编码为application/x-www-form-urlencoded格式
	encodedFormData := postData.Encode()
	// 创建bytes.Buffer对象，并将编码后的表单数据写入其中
	requestBody := bytes.NewBufferString(encodedFormData)

	reqPost, err := http.NewRequest(http.MethodPost, urlStr, requestBody)

	setHeader(reqPost)

	// 获取响应
	postResp, err := client.Do(reqPost)
	if err != nil {
		return false
	}
	defer postResp.Body.Close()

	// 判断是否存在SQL注入漏洞
	postDody, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return false
	}

	if strings.Contains(string(postDody), "SQL syntax") || strings.Contains(string(postDody), "mysql_fetch_array()") {
		log.Println(urlStr + "\t存在sql注入漏洞--post")
		return true
	}
	return false
}
