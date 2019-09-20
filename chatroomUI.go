// Demo code for unicode support (demonstrates wide Chinese characters).
package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	input := tview.NewInputField().
		SetLabel("我")

	form := tview.NewForm()

	form.AddFormItem(input)
	form.SetBorder(false).SetTitle("输入框").SetTitleAlign(tview.AlignLeft)

	clients := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("clients")
	text := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("text")
	grid := tview.NewGrid().
		SetRows(0, 3).
		SetColumns(0, 30).
		SetBorders(true).
		AddItem(form, 1, 0, 1, 2, 0, 0, true).
		AddItem(text, 0, 0, 1, 1, 0, 0, false).
		AddItem(clients, 0, 1, 1, 1, 0, 0, false)

	pages.AddPage("base", grid, true, true)
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
