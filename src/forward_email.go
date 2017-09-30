package main

import (
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"net/http"
)

func ForwardEmail(w http.ResponseWriter, r *http.Request) {
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

	context := Context{
		Domain:  profile.Domain,
		URL:     r.URL,
		From:    r.Form["from"][0],
		Subject: r.Form["subject"][0],
		Body:    r.Form["body"][0],
		Phone:   r.Form["phone"][0],
	}

	from, err := ApplyTemplate(profile.From, context)
	if err != nil {
		fmt.Println("Error: Failed to write from:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing from failed\"}")
		return
	}
	subject, err := ApplyTemplate(profile.Subject, context)
	if err != nil {
		fmt.Println("Error: Failed to write subject:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing subject failed\"}")
		return
	}
	body, err := ApplyTemplate(profile.Body, context)
	if err != nil {
		fmt.Println("Error: Failed to write body:", err)
		fmt.Fprintf(w, "{\"status\":\"Writing body failed\"}")
		return
	}
	to, err := ApplyTemplate(profile.To, context)
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
