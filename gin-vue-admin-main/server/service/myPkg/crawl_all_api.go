/*
 * @author Crabin
 */

package myPkg

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type CrawlAllApi struct {
}

var vulnerableURL = common.VulnerableURL{}

// 根据这个url返回的结果递归爬取所有url
func (w *CrawlAllApi) GetAllLinks(target string) []common.VulnerableURL {
	log.Println("开始爬取URL")
	url, err := url.Parse(target)
	if err != nil {
		log.Println(err)
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	//首先测试这个是否可以访问
	resp, err := client.Get(target)
	if err != nil || resp == nil {
		//访问异常
		return nil
	}
	defer resp.Body.Close()

	host := url.Host
	//拿到这个URL的匹配正则表达式
	regexpStr := regexp.MustCompile(target)
	re := regexp.MustCompile(`(/[a-zA-Z0-9_]+)+/\d+`)
	// Create a Collector
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile(regexpStr.String()),
			regexp.MustCompile(re.String()),
		),
	)
	// Visit only links with http and https schemes
	c.AllowedDomains = []string{host}
	c.SetRequestTimeout(time.Second * 10)
	var res = []common.VulnerableURL{}
	var resMap = make(map[string]bool)
	selector := "a[href],area[href],form[action],form[method]"
	// 找到并访问
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		//fmt.Println(e)
		link := e.Attr("href")
		// 处理form
		formVulnerableURL := findForm(e)
		if formVulnerableURL.URL != "" {
			res = append(res, formVulnerableURL)
			resMap[formVulnerableURL.URL] = true
			return
		}

		if !strings.Contains(link, "http") && strings.Compare(link, "javascript:;") != 0 {
			//没有前缀,拼接一下Less-55 -> http://localhost:82/Less-55
			link = url.Scheme + "://" + url.Host + "/" + link
		}
		//fmt.Println(link)
		if strings.Contains(link, "http") {
			// Extract request URL
			// 这个请求是否已经判断过
			if isImageURL(link) && resMap[link] {
				return
			}
			methods := []string{http.MethodGet, http.MethodPost}
			for _, method := range methods {
				resp, err := http.NewRequest(method, link, nil)
				if resp == nil || err != nil {
					//log.Fatal(err)
					return
				}
				vulnerableURL = common.VulnerableURL{
					URL:    link,
					Method: method,
				}
			}
			res = append(res, vulnerableURL)
			resMap[link] = true
			//fmt.Println(reqURL.URL)
		}
		// 继续访问
		e.Request.Visit(link)
	})
	// 开始
	if err := c.Visit(target); err != nil {
		log.Fatal(err)
	}
	return res
}

func isImageURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
}

func findForm(e *colly.HTMLElement) common.VulnerableURL {
	// 不是form
	if e.Name != "form" {
		return common.VulnerableURL{}
	}
	action := e.Attr("action")
	method := e.Attr("method")
	// url为空
	_, err := url.Parse(action)
	if len(action) == 0 || err != nil {
		return common.VulnerableURL{}
	}
	//获取表单属性
	var params []string
	// 获取form标签内的所有input、select和textarea标签并遍历
	e.DOM.Find("input, select, textarea").Each(func(_ int, s *goquery.Selection) {
		// 获取标签类型和名称
		tagName := s.Get(0).Data
		name, _ := s.Attr("name")

		// 如果是select标签，还需要获取选项的值
		if tagName == "select" {
			var options []string
			s.Find("option").Each(func(_ int, option *goquery.Selection) {
				value, _ := option.Attr("value")
				options = append(options, value)
			})
			paramValue := strings.Join(options, ",")
			params = append(params, tagName+":"+paramValue)
			return
		}
		params = append(params, tagName+":"+name)
	})

	vulnerableURL = common.VulnerableURL{
		URL:    action,
		Method: method,
		Params: params,
	}
	return vulnerableURL
}
