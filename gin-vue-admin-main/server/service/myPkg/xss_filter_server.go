/*
 * @author Crabin
 */

package myPkg

import (
	"encoding/hex"
	"github.com/feiin/go-xss"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"net/url"
	"regexp"
	"strconv"
)

type XSSFilterServer struct {
}

func (m *XSSFilterServer) XSSFilter(XSSAttack string) (string, []string) {
	source := XSSAttack
	source = DecodeHex(source)
	// 多次解码字符串
	source, _ = DecodeMultilayer(source)

	//fmt.Println(source)

	options := xss.XssOption{}

	// 只允许a标签，该标签只允许href, title, target这三个属性
	options.WhiteList = map[string][]string{
		"p":        {"style"},
		"textarea": {"style"},
		"div":      {"style"},
		"a":        {"title"},
	}

	//去掉不在白名单上的标签
	options.StripIgnoreTag = true

	//仅去掉指定的不在白名单上的标签及标签体 ；“”(空数组) 去掉所有不在白名单上的标签
	options.StripIgnoreTagBody = []string{"script"}

	x := xss.NewXSS(options)
	safeHtml := x.Process(source)

	//log.Println("前：" + source)
	//log.Println("后：" + safeHtml)

	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(source, safeHtml, false)

	// 储存疑似xss语句
	xssSentences := []string{}
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			log.Println(diff.Text)
			xssSentences = append(xssSentences, diff.Text)
		}
	}

	hasXSS := safeHtml == source
	log.Println(hasXSS)
	log.Println("一共有" + strconv.Itoa(len(xssSentences)))
	return safeHtml, xssSentences
}

// 多次解码字符串
func DecodeMultilayer(encoded string) (string, error) {
	for {
		decoded, err := url.QueryUnescape(encoded)
		if err != nil {
			return "", err
		}
		if decoded == encoded {
			break
		}
		encoded = decoded
	}
	return encoded, nil
}

// 这个函数使用正则表达式找到输入字符串中的所有16进制编码，并使用hex.DecodeString将其解码为原始字节。然后替换原始字符串中的编码为解码后的字符串。最后返回解码后的字符串。
func DecodeHex(input string) string {
	// 匹配 &#xHH; 格式的16进制编码
	re := regexp.MustCompile(`&#x([0-9A-Fa-f]+);`)
	// 通过正则表达式找到所有的16进制编码
	hexStrings := re.FindAllStringSubmatch(input, -1)
	for _, hexStr := range hexStrings {
		// 使用hex.DecodeString解码16进制字符串
		decoded, err := hex.DecodeString(hexStr[1])
		if err == nil {
			// 替换原来的编码为解码后的字符串
			input = re.ReplaceAllString(input, string(decoded))
		}
	}
	return input
}
