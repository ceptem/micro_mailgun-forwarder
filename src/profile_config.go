package main

type ProfileConfig struct {
	Domain     string `json:"domain"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
}
