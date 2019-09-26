package main

import (
	notifier "github.com/deckarep/gosx-notifier"
)

//OutputToNotificationCenter ...
func OutputToNotificationCenter(msg string) error {

	notification := notifier.Notification{
		Title:   "消息",
		Message: msg,
		Sound:   notifier.Glass,
	}

	return notification.Push()
}
