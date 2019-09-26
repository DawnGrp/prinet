package main

import (
	toast "gopkg.in/toast.v1"
)

func OutputToNotificationCenter(msg string) error {

	notification := toast.Notification{
		AppID:   "Prinet", // Shows up in the action center (lack of accent is due to encoding issues)
		Title:   "消息",
		Message: msg,
		Actions: []toast.Action{
			{"关闭", "确定"},
		},
	}

	return notification.Push()
}
