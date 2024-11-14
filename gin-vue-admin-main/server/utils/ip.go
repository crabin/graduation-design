/*
 * @author Crabin
 */

package utils

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/gin-gonic/gin"
	"github.com/ipinfo/go/v2/ipinfo"
	"golang.org/x/net/html"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

/*
 *获取请求的IP地址, hostname: www.google.com  return 这个域名的ip列表
 */
func getIPAddressFromHost(hostname string) ([]net.IP, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	return ips, nil
}

// RemoteAddr
func handlerRemoteAddr(r *http.Request) string {
	clientIP := r.RemoteAddr
	fmt.Println("Client IP:", clientIP)
	return clientIP
}

// X-Real-IP  X-Forwarded-For RemoteAddr
func GetClientIP(c *gin.Context) string {
	clientIP := c.Request.Header.Get("X-Real-IP")
	resIp := clientIP
	// 检查X-Real-IP头部是否存在
	if clientIP == "" {
		clientIP = c.Request.Header.Get("X-Forwarded-For")
		if clientIP != "" {
			//X-Forwarded-For 头部存在
			resIp = strings.Split(clientIP, ",")[0]
		}
		// X-Real-IP 和 X-Forwarded-For 头部不存在，则返回远程地址
		remoteAddr := c.Request.RemoteAddr
		if remoteAddr == "" {
			return ""
		}
		resIp = strings.Split(remoteAddr, ":")[0]
	}
	if resIp == "127.0.0.1" {
		//是本地在测试
		ipinfo.NewClient(nil, nil, common.IpInfoToten)
		resIp, _ = ipinfo.GetIPAddr()
	}

	//X-Real-IP头部存在
	return resIp
}

// 校验ip格式
func IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// 通过传入url获取其IP地址
func GetIpLookupWithUrl(url string) []net.IP {
	ips, err := net.LookupIP(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ips
}

// 检查服务协议类型
func CheckProtocolType(urlStr string) (string, error) {
	urls, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	urlR := fmt.Sprintf("%s://%s:%s", urls.Scheme, urls.Host, urls.Port())
	req, err := http.NewRequest("OPTIONS", urlR, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return strings.ToUpper(resp.Header.Get("Server")), nil
}

// 检查网站后端编程语言
func CheckXPoweredBy(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	// 获取HTTP响应头信息
	header := resp.Header
	for k, v := range header {
		if strings.ToLower(k) == "x-powered-by" {
			return v[0], nil
			/*poweredBy := strings.ToLower(v[0])
			for lang, name := range common.Languages {
				if strings.Contains(poweredBy, lang) {
					fmt.Printf("Powered by %v\n", name)
					return name, nil
				}
			}*/
		}
	}
	buf := make([]byte, 1024)
	resp.Body.Read(buf)
	body := string(buf)

	// 通过HTML源码中的meta标签获取后端编程语言信息
	langRegex := regexp.MustCompile(`(?i:<meta\s.*?name=["']?generator["']?.*?>)`)
	langMatch := langRegex.FindString(body)
	if langMatch != "" {
		lang := regexp.MustCompile(`(?i:content=["']?(.*?)["']?)`).FindStringSubmatch(langMatch)
		if len(lang) == 2 {
			fmt.Printf("Backend Programming Language: %s\n", lang[1])
			return lang[1], nil
		}
	}

	return "", nil
}

// 检查是否存在CDN
func CheckCDN(urlStr string) (bool, error) {
	//检查url
	urlCheck, err := url.Parse(urlStr)
	if err != nil || urlCheck.Host == "" {
		log.Fatal(err)
		return false, err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	resp, _ := client.Get(urlStr)

	if resp == nil {
		//log.Fatal("连接错误" + urlStr)
		return false, nil
	}
	defer resp.Body.Close()

	//HTTP响应头中是否存在Via或CDN-Cache字段
	via := resp.Header.Get("Via")
	CDNCache := resp.Header.Get("CDNCache")
	if via != "" || CDNCache != "" {
		return true, nil
	}

	// 从header中获取server信息
	server := resp.Header.Get("Server")

	hasCDN := false
	// 匹配server是否包含CDN
	if server != "" {
		hasCDN, _ = regexp.MatchString("cdn", server)
	}

	if hasCDN {
		return hasCDN, nil
	}

	//检查HTML文档中的所有<link>元素
	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
		return false, err
	}
	hasCDN = checkCDN(doc)

	return hasCDN, nil
}

// 遍历HTML文档中的所有<link>元素，查找是否存在rel属性中包含cdn字符串的元素。如果存在，则判定该网站存在CDN
func checkCDN(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "link" {
		for _, attr := range n.Attr {
			if attr.Key == "rel" && strings.Contains(attr.Val, "cdn") {
				return true
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if checkCDN(c) {
			return true
		}
	}
	return false
}
