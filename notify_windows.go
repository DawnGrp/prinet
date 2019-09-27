package main

import (
	toast "gopkg.in/toast.v1"
)

func OutputToNotificationCenter(msg string) error {

	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog", // Shows up in the action center (lack of accent is due to encoding issues)
		Title:   "prinet message",
		Message: msg,
		Actions: []toast.Action{
			{"protocol", "Close", ""},
		},
	}

	return notification.Push()
}
