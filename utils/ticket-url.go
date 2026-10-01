package utils

import (
	"net/url"
	"strings"
)

// approved ticket provider domains (subdomains are allowed)
var approvedTicketHosts = []string{
	"ticketmaster.com",
	"ticketmaster.ca",
	"ticketmaster.co.uk",
	"livenation.com",
	"ticketweb.com",
}

// IsApprovedTicketURL allows only https links to approved hostnames
func IsApprovedTicketURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if u.Scheme != "https" || u.User != nil {
		return false
	}

	host := strings.ToLower(u.Hostname())
	if host == "" {
		return false
	}

	for _, approved := range approvedTicketHosts {
		if host == approved || strings.HasSuffix(host, "."+approved) {
			return true
		}
	}
	return false
}
