// Demo code for unicode support (demonstrates wide Chinese characters).
package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	notifier "github.com/deckarep/gosx-notifier"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"gopkg.in/toast.v1"
)

var textBox *tview.TextView
var clientsBox *tview.TextView
var app *tview.Application

//ChatRoomUI 界面
func ChatRoomUI() {
	app = tview.NewApplication()
	pages := tview.NewPages()

	input := tview.NewInputField()
	input.SetFieldBackgroundColor(tcell.ColorDarkRed)
	input.SetLabel(" Say: ")
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		txt := input.GetText()

		if event.Key() == tcell.KeyEnter && len(input.GetText()) > 0 {

			LANIPS.Range(func(ip interface{}, client interface{}) bool {

				c, ok := client.(Client)
				if !ok || c.Name == "" {
					return true
				}

				err := sendMsg(ip.(string), Data{Cmd: "talk", Body: txt})
				if err != nil {

					textBox.SetText(fmt.Sprintf("%s %s:%s", textBox.GetText(false), client.(Client).Name, err.Error()))

				}
				//textBox.SetText(fmt.Sprintf("%s\n%s:%s", textBox.GetText(false), client.(Client).Name, "["+c.Name+"] "))

				input.SetText("")
				return true
			})
		}
		return event

	})

	clientsBox = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(" 无在线主机")
	clientsBox.SetBackgroundColor(tcell.ColorDarkGreen)
	textBox = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(" 聊天内容：")
	textBox.SetBackgroundColor(tcell.ColorDarkSlateGray)
	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0, 20).
		SetBorders(true).
		AddItem(input, 1, 0, 1, 2, 0, 0, true).
		AddItem(textBox, 0, 0, 1, 1, 0, 0, false).
		AddItem(clientsBox, 0, 1, 1, 1, 0, 0, false)

	pages.AddPage("base", grid, true, true)

	go func() {
		initLocalInfo()
		touch()
		refreshClients()
	}()

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

func refreshClients() {
	//刷新客户端
	if clientsBox != nil {
		clients := " "
		LANIPS.Range(func(ip, client interface{}) bool {
			clients += client.(Client).Name + "\n "
			return true
		})
		clientsBox.SetText(clients)
		app.Draw()
	}
}

func printMsg(msg string) {

	//terminal-notifier -title '💰' -message 'Check your Apple stock!'

	switch {
	case strings.Contains(runtime.GOOS, "windows"):
		//OutputToNotificationCenterWindows(msg)
	case strings.Contains(runtime.GOOS, "darwin"):
		OutputToNotificationCenterDarwin(msg)
	case strings.Contains(runtime.GOOS, "linux"):
		OutputToNotificationCenterLinux(msg)
	}

	textBox.SetText(fmt.Sprintf("%s %s", textBox.GetText(false), msg))
	app.Draw()
}

func OutputToNotificationCenterLinux(msg string) error {

	exec.Command("notify-send", msg).Run()
	return nil
}

func OutputToNotificationCenterWindows(msg string) error {

	notification := toast.Notification{
		AppID:   "Prinet", // Shows up in the action center (lack of accent is due to encoding issues)
		Title:   "消息",
		Message: msg,
	}

	return notification.Push()
}

func OutputToNotificationCenterDarwin(msg string) error {

	notification := notifier.Notification{
		Title:   "消息",
		Message: msg,
		Sound:   notifier.Glass,
	}

	return notification.Push()
}
