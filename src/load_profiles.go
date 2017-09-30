package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func LoadProfiles() {
	direntries, err := ioutil.ReadDir(Settings.ProfilesDir)
	if err != nil {
		fmt.Println("Warn: Cannot access profile directory:", err)
	}

	for _, dirent := range direntries {
		fn := dirent.Name()

		if strings.HasPrefix(fn, ".") {
			continue
		}
		if !strings.HasSuffix(fn, ".json") {
			continue
		}
		cn := strings.TrimSuffix(fn, ".json")
		fh, err := os.Open(Settings.ProfilesDir + fn)
		if err != nil {
			fmt.Println("Warn: Cannot open profile:", err)
			continue
		}
		defer fh.Close()
		parser := json.NewDecoder(fh)
		var profileConfig ProfileConfig
		err = parser.Decode(&profileConfig)
		if err != nil {
			fmt.Println("Warn: Cannot parse profile:", fn+":", err)
			continue
		}

		profile := LoadProfile(fn, &profileConfig)
		if profile == nil {
			continue
		}

		fmt.Println("Info: Loaded profile:", cn)
		profiles[cn] = *profile
	}
}
