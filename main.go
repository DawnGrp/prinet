package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
var PORT = ":10220"

//HOSTNAME ...
var HOSTNAME string

//LOCALIP 本地IP地址
var LOCALIP string

type Data struct {
	Cmd   string `json:"cmd"`
	Param string `json:"param"`
	Body  string `json:"body"`
	IP    string `json:"ip"`
}

type Client struct {
	Conn net.Conn
	Name string
}

func main() {
	go listenMsg()
	initLanIPs()
	touch()
	time.Sleep(3 * time.Second)

	ChatRoom()
}

func listenMsg() {
	addr, err := net.ResolveUDPAddr("udp", PORT)
	if err != nil {
		fmt.Println("resolve udp addr", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("listen udp", err)
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

		//解析到JSON中
		jsonData := Data{}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			fmt.Println("json format err", err.Error())
			continue
		}

		//fmt.Println("Serv Received:", string(data))

		go instructionSets(jsonData)

		//算出哈希，返回，检验正确性。
		md5string := byte2MD5string(data)
		_, err = conn.WriteToUDP([]byte(md5string), rAddr)
		if err != nil {
			fmt.Println("write md5", err)
			continue
		}

		//fmt.Println("Serv Send:", md5string)
	}

}

//通过检查指定IP是否存在链接，如果存在，发送消息
func sendMsg(ip string, data Data) (err error) {

	message, err := json.Marshal(data)

	if len(message) > BUFFERLEN {
		return fmt.Errorf("message to long")
	}

	//取出conn，如果没有，重新创建conn
	client, ok := LANIPS.Load(ip)
	c, ok := client.(Client)
	if !ok {
		return fmt.Errorf("not client object")
	}

	if !ok || c.Conn == nil {
		c.Conn, err = net.Dial("udp", ip+PORT)
		if err != nil {
			fmt.Println("connect to ", ip, err.Error())
			return err
		}

		LANIPS.Store(ip, c)
	}

	err = c.Conn.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		fmt.Println("set deadline ", ip, err.Error())
		return
	}
	//写入数据
	_, err = c.Conn.Write([]byte(message))
	if err != nil {
		return
	}

	//接收数据到hexByte
	hexByte := make([]byte, 32)
	_, err = c.Conn.Read(hexByte)
	if err != nil {
		return
	}

	//检查消息md5值与收到的是否相同
	if byte2MD5string(message) != string(hexByte) {
		fmt.Printf("%s != %s", byte2MD5string(message), string(hexByte))
		return fmt.Errorf("%s != %s", byte2MD5string(message), string(hexByte))
	}

	//fmt.Println("Send Ok!")
	return nil
}

//检查局域网内在线设备，保存到LANIPS列表中
func touch() {

	data := Data{
		IP:   LOCALIP,
		Cmd:  "uname",
		Body: "",
	}

	err := sendMsg("192.168.3.255", data)
	if err != nil {
		fmt.Println("touche", err.Error())
	}

}

func initLanIPs() (ips []*net.IPNet) {
	HOSTNAME, _ = os.Hostname()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本地网络信息:", err)
	}
	for _, a := range addrs {
		if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				fmt.Println("IP:", ip.IP)
				fmt.Println("子网掩码:", ip.Mask)
				LOCALIP = ip.IP.String()
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
		LANIPS.LoadOrStore(inetNtoA(i), Client{Conn: nil, Name: ""})
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

func instructionSets(data Data) (err error) {

	switch data.Cmd {
	case "uname":
		data2Send := Data{
			Cmd:   "mname",
			Body:  HOSTNAME,
			Param: "",
			IP:    LOCALIP,
		}

		sendMsg(data.IP, data2Send)
	case "mname":
		client, ok := LANIPS.Load(data.IP)
		if !ok {
			return fmt.Errorf("not found %s", data.IP)
		}

		c, ok := client.(Client)
		if !ok {
			return fmt.Errorf("not client object %s", data.IP)
		}

		c.Name = data.Body

		LANIPS.Store(data.IP, c)

	case "talk":
		if data.IP != LOCALIP {

			client, ok := LANIPS.Load(data.IP)
			if ok {
				fmt.Println(client.(Client).Name, ":", data.Body)
			}

		}
	}

	return
}

//ChatRoom 启动一个聊天室
func ChatRoom() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Chat Room")
	fmt.Println("---------------------")

	for {
		fmt.Print("Me: ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		LANIPS.Range(func(ip interface{}, client interface{}) bool {

			c, ok := client.(Client)

			if !ok || c.Name == "" {
				return true
			}

			err := sendMsg(ip.(string), Data{Cmd: "talk", Body: text, IP: LOCALIP})

			if err != nil {
				fmt.Println("send error:", err.Error())
			}

			fmt.Print("[" + c.Name + "] ")

			return true
		})
		fmt.Println("")
	}
}
