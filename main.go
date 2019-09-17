package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sync"
	"time"
)

//LANIPS 局域网的IP
var LANIPS = sync.Map{}

//BUFFERLEN 数据接收块大小
var BUFFERLEN = 1024

//TOUCHSTR 试探是否在线的字符串
var TOUCHSTR = "hello"

//PORT 本地监听端口
var PORT = ":8888"

func main() {
	initLanIPs()
	go listenMsg()

	time.Sleep(2 * time.Second)
	touch()
	time.Sleep(3 * time.Second)

	LANIPS.Range(func(ip interface{}, conn interface{}) bool {
		if conn != nil {
			sendMsg(ip.(string), "How are you")
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
			fmt.Println(err)
			continue
		}

		strData := string(data)
		fmt.Println("Received:", strData)

		//TODO: 算出哈希，返回，检验正确性。
		_, err = conn.WriteToUDP(data, rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Send:", strData)
	}

}

func sendMsg(ip string, message string) bool {

	if c, ok := LANIPS.Load(ip); ok && c != nil {
		i, err := c.(net.Conn).Write([]byte(message))
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		fmt.Println("send...", ip, message, i)
		return true
	}

	return false
}

func touch() {

	LANIPS.Range(func(ip interface{}, conn interface{}) bool {

		go func(ip string) {
			conn, err := net.Dial("udp", ip+PORT)
			if err != nil {
				fmt.Println("connect to ", ip, err.Error())
				return
			}

			conn.Write([]byte(TOUCHSTR))
			msg := make([]byte, len(TOUCHSTR))
			_, err = conn.Read(msg)

			if err == nil && string(msg) == TOUCHSTR {
				LANIPS.Store(ip, conn)
			} else {
				fmt.Println(err, string(msg), string(msg) == TOUCHSTR)
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
