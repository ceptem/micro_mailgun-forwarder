package main

import (
	"fmt"
	"os"
	"strconv"
)

func FindStringSetting(name string, defaultValue string) string {
	cand := os.Getenv(name)
	if cand != "" {
		return cand
	}
	return defaultValue
}

func FindIntSetting(name string, defaultValue int) int {
	cand := os.Getenv(name)
	if cand != "" {
		val, err := strconv.Atoi(cand)
		if err != nil {
			fmt.Println("Alert: Failed to read", name+":", err)
			os.Exit(1)
		}
		return val
	}
	return defaultValue
}

func FindUint16Setting(name string, defaultValue uint16) uint16 {
	cand := os.Getenv(name)
	if cand != "" {
		val, err := strconv.Atoi(cand)
		if err != nil || val < 0 || val > 0xFFFF {
			fmt.Println("Alert: Failed to read", name+":", err)
			os.Exit(1)
		}
		return uint16(val)
	}
	return defaultValue
}

func FindSettings() {
	Settings.UserID = FindIntSetting("UID", DefaultUserID)
	Settings.GroupID = FindIntSetting("GID", DefaultGroupID)
	Settings.Port = FindUint16Setting("PORT", DefaultListenPort)
	Settings.Address = FindStringSetting("ADDRESS", DefaultListenAddress)
	Settings.ProfilesDir = FindStringSetting(
		"PROFILEDIR", DefaultProfileDirectory)
}
