package main

import "os/exec"

func OutputToNotificationCenter(msg string) error {

	exec.Command("notify-send", msg).Run()
	return nil
}
