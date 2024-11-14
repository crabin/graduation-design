/*
 * @author Crabin
 */

package utils

import (
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"io/ioutil"
	"log"
	"net/http"
)

// 基于wappalyzer工具获取web的信息
func GetWebInfo(urlStr string) map[string]struct{} {
	resp, err := http.DefaultClient.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)

	wappalyzerClient, err := wappalyzer.New()

	fingerprints := wappalyzerClient.Fingerprint(resp.Header, data)

	return fingerprints
	// Output: map[Acquia Cloud Platform:{} Amazon EC2:{} Apache:{} Cloudflare:{} Drupal:{} PHP:{} Percona:{} React:{} Varnish:{}]
}

/**
原理

wappalyzer 通过两种方式来区分和识别得到的数据：

指纹库
wappalyzer 使用指纹库来识别 Web 应用程序使用的技术和框架，包括后端 API。
指纹库是一组正则表达式、脚本和模式，通过匹配页面内容、HTTP 响应头、HTML 标记等方式来识别 Web 应用程序使用的特定技术和框架。
指纹库可以通过更新配置文件或者使用提供的公共指纹库来获得最新的信息。

线索
wappalyzer 通过使用预定义的规则和检测器来寻找线索来识别 Web 应用程序使用的技术和框架。
例如，如果页面中出现了一个名为 wp-content 的文件或目录，那么 wappalyzer 可能会推断该应用程序是基于WordPress 的。
类似地，如果页面中包含某些 JavaScript 对象或函数，wappalyzer 可能会根据这些信息推断出页面使用了 jQuery 或其他 JavaScript 库。

需要注意的是，由于 wappalyzer 使用的是一组预先定义好的规则和指纹库，它可能会出现误判的情况，
例如将一个基于 Flask 的 Web 应用程序错误地识别为基于 Django 的应用程序。
为了避免这种情况，可以使用人工审核和校正的方法来推动更加准确的识别。
*/
