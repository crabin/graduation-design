/*
 * @author Crabin
 */

package test

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"strconv"
	"testing"
)

var (
	host     = "192.168.13.128"
	port     = 22
	user     = "root"
	password = "1234567899"
)

func TestSSHConn(t *testing.T) {
	/*
		cmd := "gobuster dirb -u https://www.cnblogs.com/crabin/ -w /usr/share/wordlists/dirb/common.txt"

		// 创建 SSH 客户端配置
		sshConfig := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// 连接 SSH 服务器
		conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), sshConfig)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// 创建 SSH 会话
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()

		// 执行命令
		output, err := session.CombinedOutput(cmd)
		if err != nil {
			panic(err)
		}

		// 打印命令输出
		fmt.Println(string(output))*/
	// 连接到远程服务器
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 执行命令
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}

	const cmd = "gobuster dirb -u https://www.cnblogs.com/crabin/ -w /usr/share/wordlists/dirb/common.txt"
	if err := session.Run(cmd); err != nil {
		panic(err)
	}

	// 实时输出命令的标准输出和错误输出
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	// 等待命令执行完成
	if err := session.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("Done.")
}
