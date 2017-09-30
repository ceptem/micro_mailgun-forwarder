package main

import (
	"syscall"
)

func ConfigureRuntime() {
	if Settings.GroupID != -1 && Settings.GroupID != syscall.Getgid() {
		syscall.Setuid(Settings.GroupID)
	}
	if Settings.UserID != -1 && Settings.UserID != syscall.Getuid() {
		syscall.Setuid(Settings.UserID)
	}
}
