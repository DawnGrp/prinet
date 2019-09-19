package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/text"
)

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
		if len(text) == 0 {
			continue
		}
		fmt.Print("receiv:")
		LANIPS.Range(func(ip interface{}, client interface{}) bool {

			c, ok := client.(Client)

			if !ok || c.Name == "" {
				//fmt.Println("send To ", ip, ok, c.Name)
				return true
			}

			err := sendMsg(ip.(string), Data{Cmd: "talk", Body: text})

			if err != nil {
				fmt.Println("send error:", err.Error())
			}

			fmt.Print("[" + c.Name + "] ")

			return true
		})
		fmt.Println("")
	}
}

func ChatRoomUI() {

	txt := text.New("hello world")
	app, _ := gowid.NewApp(gowid.AppArgs{View: txt})
	app.SimpleMainLoop()

	app.SimpleMainLoop()
}
