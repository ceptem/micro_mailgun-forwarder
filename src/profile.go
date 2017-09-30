package main

import "text/template"

type Profile struct {
	Domain     string
	PrivateKey string
	PublicKey  string
	From       *template.Template
	To         *template.Template
	Subject    *template.Template
	Body       *template.Template
}
