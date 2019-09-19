package test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

func TestHex(t *testing.T) {
	data := []byte("")
	//data = bytes.TrimSpace(data)
	t.Log(data, string(data))
	arpLan()
}
func String2MD5(plainText string) (md5string string) {

	fmt.Println(plainText)
	m := md5.New()

	md5.Sum([]byte(plainText))

	return hex.EncodeToString(m.Sum(nil))
}

func arpLan() {

	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmd := exec.Command("arp", "-a")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(opBytes))

}
