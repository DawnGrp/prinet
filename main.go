package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//LANIPS 局域网的IP
var LANIPS = sync.Map{}

//BUFFERLEN 数据接收块大小
var BUFFERLEN = 1024

//PORT 本地监听端口
var PORT = ":8888"

//HOSTNAME ...
var HOSTNAME string

func main() {
	HOSTNAME, _ = os.Hostname()
	go listenMsg()
	initLanIPs()
	touch()
	time.Sleep(3 * time.Second)

	LANIPS.Range(func(ip interface{}, conn interface{}) bool {
		if conn != nil {
			fmt.Println(sendMsg(ip.(string), []byte("How are you")))
		}

		return true
	})

	time.Sleep(5 * time.Second)
}

func listenMsg() {
	addr, err := net.ResolveUDPAddr("udp", PORT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, BUFFERLEN)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Serv read", err)
			continue
		}

		//去掉多余字节
		index := bytes.IndexByte(data, 0)
		if index > -1 {
			data = data[0:index]
		}

		strData := string(data)
		fmt.Println("Serv Received:", strData)

		md5string := byte2MD5string(data)

		//算出哈希，返回，检验正确性。
		_, err = conn.WriteToUDP([]byte(md5string), rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Serv Send:", md5string)
	}

}

//通过检查指定IP是否存在链接，如果存在，发送消息
func sendMsg(ip string, message []byte) (err error) {

	//取出conn，如果没有，重新创建conn
	conn, ok := LANIPS.Load(ip)
	if !ok || conn == nil {
		conn, err = net.Dial("udp", ip+PORT)
		if err != nil {
			fmt.Println("connect to ", ip, err.Error())
			return
		}
		LANIPS.Store(ip, conn)
	}
	err = conn.(net.Conn).SetDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		fmt.Println("set deadline ", ip, err.Error())
		return
	}
	//写入数据
	_, err = conn.(net.Conn).Write([]byte(message))
	if err != nil {
		return
	}

	//接收数据到hexByte
	hexByte := make([]byte, 32)
	_, err = conn.(net.Conn).Read(hexByte)
	if err != nil {
		return
	}

	//检查消息md5值与收到的是否相同
	if byte2MD5string(message) != string(hexByte) {
		fmt.Printf("%s != %s", byte2MD5string(message), string(hexByte))
		return fmt.Errorf("%s != %s", byte2MD5string(message), string(hexByte))
	}

	fmt.Println("Send Ok!")
	return nil
}

//检查局域网内在线设备，保存到LANIPS列表中
func touch() {

	LANIPS.Range(func(ip interface{}, conn interface{}) bool {

		go func(ip string) {
			err := sendMsg(ip, []byte("[Hello]"))
			if err != nil {
				fmt.Println("touche", err.Error())
			}
		}(ip.(string))

		return true
	})

}

func initLanIPs() (ips []*net.IPNet) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本地网络信息:", err)
	}
	for _, a := range addrs {
		if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				fmt.Println("IP:", ip.IP)
				fmt.Println("子网掩码:", ip.Mask)
				lanIPs(ip)
			}
		}
	}

	return
}

func lanIPs(ipNet *net.IPNet) {
	ip := ipNet.IP.To4()

	var min, max uint32

	for i := 0; i < 4; i++ {
		b := uint32(ip[i] & ipNet.Mask[i])
		min += b << ((3 - uint(i)) * 8)
	}
	one, _ := ipNet.Mask.Size()
	max = min | uint32(math.Pow(2, float64(32-one))-1)
	log.Printf("内网IP范围:%d - %d", min, max)
	// max 是广播地址，忽略
	// i & 0x000000ff  == 0 是尾段为0的IP，根据RFC的规定，忽略
	for i := min; i < max; i++ {
		if i&0x000000ff == 0 {
			continue
		}
		LANIPS.LoadOrStore(inetNtoA(i), nil)
	}

	return
}

func inetNtoA(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func byte2MD5string(plaindata []byte) (md5string string) {
	m := md5.New()
	m.Write(plaindata)
	return hex.EncodeToString(m.Sum(nil))
}

func instructionSets(cmd string) (err error) {

	switch {
	case strings.HasPrefix(cmd, "[uname]"):

	case strings.HasPrefix(cmd, "[mname]"):

	case strings.HasPrefix(cmd, "[talk]"):

	}

	return
}
