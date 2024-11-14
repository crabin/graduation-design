/*
 * @author Crabin
 */

package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// 验证url格式是否正确
func MatchUrl(urlStr string) bool {
	if l := strings.Split(urlStr, "?"); len(l) > 2 {
		return false
	}

	urlS, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	if urlS.Host == "localhost:8080" {
		return true
	}

	//reg := "^(https?://)?([\\da-z.-]+)\\.([a-z.]{2,6})(:[0-9]+)?([/\\w.-]*)*/?\\?([\\w=&]*)$\n"
	reg := `^(http|https)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z0-9]{2,3}(:[a-zA-Z0-9]*)?\/?.*$`

	//re := regexp.MustCompile("^(http(s)?://)?(localhost:[0-9]+|([a-zA-Z0-9]+[.])?[a-zA-Z0-9]+([.][a-zA-Z]{2,})*)(:[0-9]+)?((/|\\?).*)?$")
	re, err := regexp.Compile(reg)
	if err != nil {
		return false
	}

	if re.MatchString(urlStr) {
		return true
	}
	return false
}

// 验证域名是否正确
func MatchRealm(realm string) bool {
	re := regexp.MustCompile("^[a-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?$")
	matches := re.FindStringSubmatch(realm)
	if len(matches) > 0 {
		return true
	}
	return false
}

// 在url后面拼接 ,, params 中的字符串为 eg： id=asdf
func MontageURL(urlStr string, params string) string {

	if len(params) == 0 {
		return urlStr
	}
	if strings.Contains(urlStr, "?") {
		return urlStr + "&" + params
	}
	if strings.HasSuffix(urlStr, "/") {
		//是否结尾为 /
		urlStr = urlStr + "?"
	} else {
		urlStr = urlStr + "/?"
	}
	params = strings.ReplaceAll(params, " ", "+")
	if !strings.Contains(params, "=") {
		// 不是id=asdf这种格式
		return urlStr
	}
	urlStr = urlStr + params

	return urlStr
}

// www.baidu.com/li/r/a?id=1&name=2&pa=3
// 处理url，提取url中的参数  return www.baidu.com/li/r/a []string{id,name,pa}
func HandleUrl(urlStr string) (string, []string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
		return urlStr, nil
	}
	// 获取所有参数
	params := u.Query()
	res := []string{}
	// 遍历所有参数
	for key, _ := range params {
		res = append(res, key)
	}
	return u.Scheme + "://" + u.Host + u.Path, res
}
