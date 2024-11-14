/*
 * @author Crabin
 */

package myPkg

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	myRes "github.com/flipped-aurora/gin-vue-admin/server/model/my/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/wsConn"
	"github.com/ipinfo/go/v2/ipinfo"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

type IpInfoServer struct {
}

// 追踪这个域名，查找真实IP
func (m *IpInfoServer) Tracert(host string, wsConn *wsConn.WsConnection) {
	log.Println(host)
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
		wsConn.WriteMessage(1, []byte(err.Error()))
	}
	var dst net.IPAddr
	for _, ip := range ips {
		if ip.To4() != nil {
			dst.IP = ip
			fmt.Printf("using %v for tracing an IP packet route to %s\n", dst.IP, host)
			break
		}
	}
	if dst.IP == nil {
		log.Fatal("no A record found")
	}

	c, err := net.ListenPacket("ip4:1", "0.0.0.0") // ICMP for IPv4
	if err != nil {
		log.Fatal(err)
		wsConn.WriteMessage(1, []byte(err.Error()))
	}
	defer c.Close()
	p := ipv4.NewPacketConn(c)

	if err1 := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		log.Fatal(err1)
		wsConn.WriteMessage(1, []byte(err1.Error()))
	}
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			//这里的ID是进程号，用于区分不同的程序，因为这个字段在报文中是16位的，所以和0xffff做了与运算
			ID:   os.Getpid() & 0xffff,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	rb := make([]byte, 1500)
	flag := false
	for i := 1; i <= 31; i++ { // up to 64 hops
		wm.Body.(*icmp.Echo).Seq = i
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatal(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
		}
		//这里是ICMP报文中的序列号，用于区分发送的第几个ICMP数据报
		if err := p.SetTTL(i); err != nil {
			log.Fatal(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
		}

		begin := time.Now()
		if _, err := p.WriteTo(wb, nil, &dst); err != nil {
			log.Fatal(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
		}
		if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			log.Fatal(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
		}
		n, cm, peer, err := p.ReadFrom(rb)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				fmt.Printf("%v\t*\n", i)
				continue
			}
			log.Fatal(err)
		}
		rm, err := icmp.ParseMessage(1, rb[:n])
		if err != nil {
			log.Fatal(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
		}
		rtt := time.Since(begin)

		res := myRes.IpInfos{}

		switch rm.Type {
		case ipv4.ICMPTypeTimeExceeded:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\t%+v\n", i, peer, names, rtt, cm)
		case ipv4.ICMPTypeEchoReply:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\t%+v\n", i, peer, names, rtt, cm)
			flag = true
		default:
			log.Printf("unknown ICMP message: %+v\n", rm)
		}
		client := ipinfo.NewClient(nil, nil, common.IpInfoToten)
		info, err := client.GetIPInfo(net.ParseIP(peer.String()))
		if err != nil {
			log.Fatal(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
			return
		}
		res.Info = info
		res.Time = rtt
		res.Status = 1
		d1, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			wsConn.WriteMessage(1, []byte(err.Error()))
			return
		}
		wsConn.WriteMessage(1, d1)
		if flag {
			break
		}
	}
	log.Println("追踪完成！")
	d, _ := json.Marshal(myRes.IpInfos{
		Status: 2,
		Msg:    "追踪完成",
	})
	wsConn.WriteMessage(1, d)
}
