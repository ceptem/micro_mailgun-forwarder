package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

import "gopkg.in/mailgun/mailgun-go.v1"

const (
	profileDirectory = "/etc/ceptem/us/mailgun-forward-email/"
)

type ProfileConfig struct {
	Domain     string `json:"domain"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
}

type Profile struct {
	Domain     string
	PrivateKey string
	PublicKey  string
	From       *template.Template
	To         *template.Template
	Subject    *template.Template
	Body       *template.Template
}

var profiles = make(map[string]Profile)

func loadTemplate(name string, body string) *template.Template {
	template, err := template.New(name).Parse(body)
	if err != nil {
		fmt.Println("Alert: Failed to parse template:", err)
		os.Exit(1)
	}
	return template
}

func loadProfile(fn string, config *ProfileConfig) *Profile {
	return &Profile{
		Domain:     config.Domain,
		PrivateKey: config.PrivateKey,
		PublicKey:  config.PublicKey,
		From:       loadTemplate(fn+"/from", config.From),
		To:         loadTemplate(fn+"/to", config.From),
		Subject:    loadTemplate(fn+"/subject", config.Subject),
		Body:       loadTemplate(fn+"/body", config.Body),
	}
}

func loadProfiles() {
	direntries, err := ioutil.ReadDir(profileDirectory)
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
		fh, err := os.Open(profileDirectory + fn)
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

		profile := loadProfile(fn, &profileConfig)
		if profile == nil {
			continue
		}

		fmt.Println("Info: Loaded profile:", cn)
		profiles[cn] = *profile
	}
}

type Variables struct {
	Domain  string
	URL     *url.URL
	From    string
	Subject string
	Body    string
	Phone   string
}

func executeTemplate(
	tpl *template.Template,
	variables Variables,
) (string, error) {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, variables)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func forwardEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Info: Received forwarding request")
	r.ParseForm()

	valid := true
	valid = valid && len(r.Form["profile"]) == 1
	valid = valid && len(r.Form["profile"][0]) > 1
	valid = valid && len(r.Form["from"]) == 1
	valid = valid && len(r.Form["from"][0]) > 1
	valid = valid && len(r.Form["body"]) == 1
	valid = valid && len(r.Form["body"][0]) > 1
	valid = valid && len(r.Form["phone"]) == 1
	valid = valid && len(r.Form["subject"]) == 1
	if !valid {
		fmt.Println("Error: Invalid request arguments")
		fmt.Fprintf(w, "{\"status\":\"Invalid request\"}")
		return
	}

	profile, exists := profiles[r.Form["profile"][0]]
	if !exists {
		fmt.Println("Error: Invalid profile requested")
		fmt.Fprintf(w, "{\"status\":\"Invalid request\"}")
		return
	}

	var variables = Variables{
		Domain:  profile.Domain,
		URL:     r.URL,
		From:    r.Form["from"][0],
		Subject: r.Form["subject"][0],
		Body:    r.Form["body"][0],
		Phone:   r.Form["phone"][0],
	}

	from, err := executeTemplate(profile.From, variables)
	if err != nil {
		fmt.Println("Error: Failed to write from:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing from failed\"}")
		return
	}
	subject, err := executeTemplate(profile.Subject, variables)
	if err != nil {
		fmt.Println("Error: Failed to write subject:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing subject failed\"}")
		return
	}
	body, err := executeTemplate(profile.Body, variables)
	if err != nil {
		fmt.Println("Error: Failed to write body:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing body failed\"}")
		return
	}
	to, err := executeTemplate(profile.To, variables)
	if err != nil {
		fmt.Println("Error: Failed to write to:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing to failed\"}")
		return
	}

	mg := mailgun.NewMailgun(
		profile.Domain,
		profile.PrivateKey,
		profile.PublicKey,
	)
	message := mailgun.NewMessage(
		from,
		subject,
		body,
		to,
	)

	resp, id, err := mg.Send(message)
	if err != nil {
		fmt.Println("Error: Failed to forward:", err)
		fmt.Fprintf(w, "{\"status\":\"Forwarding failed\"}")
		return
	}
	fmt.Printf("Info: Queued by relay: id=%s response=%s\n", id, resp)
	fmt.Fprintf(w, "{\"status\":\"Forwarding succeeded\"}")
}

func main() {
	fmt.Println("Info: Loading profiles")
	loadProfiles()
	fmt.Println("Info: Starting service")
	http.HandleFunc("/", forwardEmail)
	err := http.ListenAndServe("0.0.0.0:6543", nil)
	if err != nil {
		log.Fatal("Alert: Failed to start serving: ", err)
	}
}
