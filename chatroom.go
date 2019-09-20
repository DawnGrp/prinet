package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rivo/tview"
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

//ChatRoomUI 带UI的聊天室
func ChatRoomUI() {

	box := tview.NewBox().SetBorder(true).SetTitle("ChatBox!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}

}
