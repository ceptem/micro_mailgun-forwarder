package main

import (
	"fmt"
	"log"
	"net/http"
)

var profiles = make(map[string]Profile)

func main() {
	fmt.Println("Info: Finding settings")
	FindSettings()
	fmt.Println("Info: Switching UID (possibly dropping privileges)")
	ConfigureRuntime()
	fmt.Println("Info: Loading profiles")
	LoadProfiles()
	fmt.Println("Info: Starting service")
	http.HandleFunc("/", ForwardEmail)
	err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", Settings.Address, Settings.Port),
		nil,
	)
	if err != nil {
		log.Fatal("Alert: Failed to start serving: ", err)
	}
}
