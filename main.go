package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
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

//Data 数据
type Data struct {
	Cmd   string `json:"cmd"`
	Param string `json:"param"`
	Body  string `json:"body"`
}

//Client 终端
type Client struct {
	Conn net.Conn
	Name string
}

func main() {
	go listenMsg()
	defer quit()
	time.Sleep(1 * time.Second)
	ChatRoomUI()

}

func listenMsg() {
	addr, err := net.ResolveUDPAddr("udp", PORT)
	if err != nil {
		printMsg("resolve udp addr" + err.Error())
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		printMsg("listen udp" + err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, BUFFERLEN)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			printMsg("Serv read:" + err.Error())
			continue
		}

		if _, err := checkClient(rAddr.IP.String()); err != nil {
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
			printMsg("json format err:" + err.Error())
			continue
		}

		go instructionSets(rAddr.IP.String(), jsonData)

		//算出哈希，返回，检验正确性。
		md5string := byte2MD5string(data)
		_, err = conn.WriteToUDP([]byte(md5string), rAddr)
		if err != nil {
			printMsg("write md5:" + err.Error())
			continue
		}

		//printMsg("Serv Send:", md5string, rAddr.IP.String())
	}

}

//通过检查指定IP是否存在链接，如果存在，发送消息
func sendMsg(ip string, data Data) (err error) {

	message, err := json.Marshal(data)

	if len(message) > BUFFERLEN {
		return fmt.Errorf("message to long")
	}

	c, err := checkClient(ip)
	if err != nil {
		return fmt.Errorf("no client")
	}

	err = c.Conn.SetDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		printMsg("set deadline " + ip + err.Error())
		return
	}
	//写入数据
	_, err = c.Conn.Write([]byte(message))
	if err != nil {
		c.Conn.Close()
		LANIPS.Delete(ip)
		return
	}

	//接收数据到hexByte
	hexByte := make([]byte, 32)
	_, err = c.Conn.Read(hexByte)
	if err != nil {
		c.Conn.Close()
		LANIPS.Delete(ip)
		return
	}

	//检查消息md5值与收到的是否相同
	if byte2MD5string(message) != string(hexByte) {
		fmt.Printf("%s != %s", byte2MD5string(message), string(hexByte))
		return fmt.Errorf("%s != %s", byte2MD5string(message), string(hexByte))
	}

	return nil
}

//检查局域网内在线设备，保存到LANIPS列表中
func touch() {

	data := Data{
		Cmd:  "uname",
		Body: "",
	}

	err := sendMsg("255.255.255.255", data)
	if err != nil {
		printMsg("touch uname:" + err.Error())
	}

	data.Cmd = "mname"
	data.Body = HOSTNAME
	err = sendMsg("255.255.255.255", data)
	if err != nil {
		printMsg("touch mname:" + err.Error())
	}

}

//initLocalInfo ...
func initLocalInfo() (ips []*net.IPNet) {
	HOSTNAME, _ = os.Hostname()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("无法获取本地网络信息:", err)
	}
	for _, a := range addrs {
		if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				printMsg("IP:" + ip.IP.String())
				printMsg("子网掩码:" + ip.Mask.String())
				LOCALIP = ip.IP.String()

			}
		}
	}

	return
}

func byte2MD5string(plaindata []byte) (md5string string) {
	m := md5.New()
	m.Write(plaindata)
	return hex.EncodeToString(m.Sum(nil))
}

func checkClient(ip string) (c Client, err error) {
	//取出conn，如果没有，重新创建conn
	client, ok := LANIPS.Load(ip)
	if !ok {
		conn, err := net.Dial("udp", ip+PORT)
		if err != nil {
			printMsg("connect to :" + ip + err.Error())
			return c, err
		}

		client = Client{Conn: conn, Name: ""}
		LANIPS.Store(ip, client)
	}
	c, ok = client.(Client)
	if !ok {
		return c, fmt.Errorf("bad Type Client")
	}
	return c, err
}

func instructionSets(ip string, data Data) (err error) {

	switch data.Cmd {
	case "uname":
		data2Send := Data{
			Cmd:   "mname",
			Body:  HOSTNAME,
			Param: "",
			//IP:    LOCALIP,
		}

		err := sendMsg(ip, data2Send)
		if err != nil {
			printMsg("re my name :" + err.Error())
		}
		printMsg("被询问" + ip)
	case "mname":
		client, ok := LANIPS.Load(ip)
		if !ok {
			return fmt.Errorf("not found %s", ip)
		}

		c, ok := client.(Client)
		if !ok {
			return fmt.Errorf("not client object %s", ip)
		}

		printMsg("收到答复:" + c.Name + data.Body)
		c.Name = data.Body

		LANIPS.Store(ip, c)

		refreshClients()

	case "talk":

		client, ok := LANIPS.Load(ip)
		if ok {
			//printMsg(client.(Client).Name, ":", data.Body)
			// textBox.SetText(fmt.Sprintf("%s %s:%s", textBox.GetText(false), client.(Client).Name, data.Body))
			// app.Draw()

			printMsg(client.(Client).Name + ":" + data.Body)
		}

	case "quit":
		client, ok1 := LANIPS.Load(ip)

		if c, ok2 := client.(Client); ok1 && ok2 && c.Conn != nil {
			c.Conn.Close()
		}

		LANIPS.Delete(ip)
	}

	return
}

func quit() {

	data := Data{
		Cmd:  "quit",
		Body: "",
	}

	err := sendMsg("255.255.255.255", data)
	if err != nil {
		fmt.Println("i am quit", err.Error())
	}

}
