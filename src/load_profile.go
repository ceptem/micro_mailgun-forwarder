package main

func LoadProfile(fn string, config *ProfileConfig) *Profile {
	return &Profile{
		Domain:     config.Domain,
		PrivateKey: config.PrivateKey,
		PublicKey:  config.PublicKey,
		From:       LoadTemplate(fn+"/from", config.From),
		To:         LoadTemplate(fn+"/to", config.From),
		Subject:    LoadTemplate(fn+"/subject", config.Subject),
		Body:       LoadTemplate(fn+"/body", config.Body),
	}
}
