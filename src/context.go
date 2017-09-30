package main

import "net/url"

type Context struct {
	Domain  string
	URL     *url.URL
	From    string
	Subject string
	Body    string
	Phone   string
}
