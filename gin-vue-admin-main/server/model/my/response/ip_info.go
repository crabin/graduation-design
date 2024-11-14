/*
 * @author Crabin
 */

package response

import (
	"github.com/ipinfo/go/v2/ipinfo"
	"time"
)

type IpInfos struct {
	ServerType string              `json:"serverType"` // 服务类型
	Info       *ipinfo.Core        `json:"info"`       // IP地理位置信息
	Time       time.Duration       `json:"time"`       // 时间
	Status     int                 `json:"status"`     // 连接代码，是否错误，0：error ; 1: success\
	Msg        string              `json:"msg"`        //信息
	XPoweredBy string              `json:"xPoweredBy"` // 使用的编程语言标签
	WebInfo    map[string]struct{} `json:"webInfo"`    // 内置信息
}
